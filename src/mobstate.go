// Copyright 2015 Pikkpoiss
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"../lib/twodee"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"time"
)

type Mob interface {
	Bored(time.Duration) bool
	SetFrames(f []int)
	Detect(dist float32) bool
	Pos() twodee.Point
	Bounds() twodee.Rectangle
	SearchPattern() []mgl32.Vec2
	Speed() float32
	MoveTo(twodee.Point)
	ShouldSwing(p mgl32.Vec2) bool
}

type Mobile struct {
	DetectionRadius float32
	BoredThreshold  time.Duration
	speed           float32
	searchPattern   []mgl32.Vec2
}

func (m *Mobile) Bored(d time.Duration) bool {
	return d >= m.BoredThreshold
}

func (m *Mobile) Detect(d float32) bool {
	return d <= m.DetectionRadius
}

func (m *Mobile) SearchPattern() []mgl32.Vec2 {
	return m.searchPattern
}

func (m *Mobile) Speed() float32 {
	return m.speed
}

// TODO: this should probably just kill the player?
func (m *Mobile) HandleCollision(p *Player) {}

// MoveMob moves the mob along the given vector, which should be normalized.
func MoveMob(m Mob, v mgl32.Vec2, l *Level) {
	// Fuck all this collision noise.
	pos := m.Pos()
	end := twodee.Pt(pos.X()+v[0], pos.Y()+v[1])
	m.MoveTo(end)
}

// MobState is implemented by various states responsible for controlling mobile
// entities.
// ExamineWorld examines the current state of the game and determines which
// state
type MobState interface {
	// ExamineWorld examines the current state of the game and determines
	// which state the mobile should transition to. It's legal to return
	// either the current MobState or nil, indicating that the mobile
	// should transition to a previous state.
	ExamineWorld(Mob, *Level) (newState MobState)
	// Update should be called each frame and may update values in the
	// current state or call functions on the mob.
	Update(Mob, time.Duration)
	// Enter should be called when entering this MobState.
	Enter(Mob)
	// Enter should be called when exiting this MobState.
	Exit(Mob)
}

type BaseState struct {
	Name string
}

func (s *BaseState) ExamineWorld(m Mob, l *Level) MobState {
	return s
}
func (s *BaseState) Update(m Mob, d time.Duration) {}
func (s *BaseState) Enter(m Mob) {
	//	fmt.Printf("Entering the %v state.\n", s.Name)
}
func (s *BaseState) Exit(m Mob) {
	//	fmt.Printf("Exiting the %v state.\n", s.Name)
}

// VegState encapsulates the state of being a vegetable.
type VegState struct {
	*BaseState
}

func NewVegState() *VegState {
	return &VegState{&BaseState{"Veggie"}}
}

// ExamineWorld always returns a new SearchState.
func (s *VegState) ExamineWorld(m Mob, l *Level) MobState {
	return &SearchState{
		Pattern:        m.SearchPattern(),
		targetPointIdx: 0,
		path:           []twodee.GridPoint{},
		pathIdx:        0,
		pathAge:        maxPathAge,
		BaseState:      &BaseState{"Search"}}
}

// SearchState is the state during which a mobile is aimlessly wandering,
// hoping to chance across the player.
// TODO: Implement Enter and Exit to have searching animations.
type SearchState struct {
	Pattern          []mgl32.Vec2
	targetPointIdx   int
	path             []twodee.GridPoint
	pathIdx, pathAge int
	*BaseState
}

// ExamineWorld returns HuntState if the player is seen, otherwise the mob
// continues wandering according to its search pattern.
func (s *SearchState) ExamineWorld(m Mob, l *Level) MobState {
	s.pathAge++
	g := l.BossCollisions
	l.SetBossPath(s.path)
	if playerSeen(m, l) {
		return NewHuntState()
	}
	if len(s.Pattern) == 0 {
		// Do nothing right now with no search pattern.
		return s
	}
	cTarget := s.Pattern[s.targetPointIdx]
	if s.pathAge > maxPathAge {
		if path := getPath(g, m.Pos(), twodee.Point{cTarget}); len(path) > 0 {
			s.path = path
			s.pathAge = 0
			s.pathIdx = 0
		}
		// Otherwise retain current path.
	}
	if len(s.path) == 0 {
		return s // Crap, still no path generated...
	}
	mv := mgl32.Vec2{m.Pos().X(), m.Pos().Y()}
	for s.pathIdx < len(s.path)-1 {
		tv := mgl32.Vec2{
			g.InversePosition(s.path[s.pathIdx].X, 0.5),
			g.InversePosition(s.path[s.pathIdx].Y, 0.5),
		}
		if tv.Sub(mv).Len() >= 2 {
			break
		}
		s.pathIdx++
	}
	tv := mgl32.Vec2{
		g.InversePosition(s.path[s.pathIdx].X, 0.5),
		g.InversePosition(s.path[s.pathIdx].Y, 0.5),
	}
	MoveMob(m, tv.Sub(mv).Normalize().Mul(m.Speed()), l)
	if s.pathIdx == len(s.path)-1 && tv.Sub(mv).Len() < 2 {
		s.targetPointIdx = (s.targetPointIdx + 1) % len(s.Pattern)
	}
	return s
}

const maxPathAge = 2

// HuntState is the state during which a mobile is actively hunting the player.
type HuntState struct {
	durSinceLastContact time.Duration
	path                []twodee.GridPoint
	pathIdx, pathAge    int
	*BaseState
}

func NewHuntState() *HuntState {
	return &HuntState{
		durSinceLastContact: 0,
		path:                []twodee.GridPoint{},
		pathIdx:             0,
		pathAge:             maxPathAge,
		BaseState:           &BaseState{"Hunt"},
	}
}

// ExamineWorld returns the current state if the player is currently seen or
// the mob is not yet tired of chasing. Otherwise, it returns nil.
func (s *HuntState) ExamineWorld(m Mob, l *Level) (ns MobState) {
	s.pathAge++
	g := l.BossCollisions
	if playerSeen(m, l) {
		// Try to generate a new path if the player is seen and our
		// last one is stale.
		if s.pathAge > maxPathAge {
			if path := getPath(g, m.Pos(), l.Player.Pos()); len(path) > 0 {
				s.path = path
				s.pathAge = 0
				s.pathIdx = 0
			}
			// Otherwise keep our current path, since we couldn't make
			// a new one.
		}
		s.durSinceLastContact = time.Duration(0)
		ns = s
	}
	l.SetBossPath(s.path)
	if !m.Bored(s.durSinceLastContact) {
		ns = s
	}
	if ns != nil {
		pv := mgl32.Vec2{l.Player.Pos().X(), l.Player.Pos().Y()}
		if m.ShouldSwing(pv) {
			// Return Swing state.
		}

		// We've passed the end of the path; there's nothing left to do.
		// Hopefully a new path will be generated within a few frames.
		if s.pathIdx == len(s.path) {
			return ns
		}
		mv := m.Pos().Vec2
		for s.pathIdx < len(s.path)-1 { // Never roll off the end.
			tv := mgl32.Vec2{
				g.InversePosition(s.path[s.pathIdx].X, 0.5),
				g.InversePosition(s.path[s.pathIdx].Y, 0.5),
			}
			if tv.Sub(mv).Len() >= 2 {
				break
			}
			s.pathIdx++
		}
		// Chase player!
		tv := mgl32.Vec2{
			g.InversePosition(s.path[s.pathIdx].X, 0.5),
			g.InversePosition(s.path[s.pathIdx].Y, 0.5),
		}
		MoveMob(m, tv.Sub(mv).Normalize().Mul(m.Speed()), l)
	}
	return ns
}

// Update resets the player's hiding timer if the player is seen, otherwise it
// increments.
func (s *HuntState) Update(m Mob, d time.Duration) {
	s.durSinceLastContact += d
}

// playerSeen returns true if the player is currently visible to the mob and
// within its detection radius.
func playerSeen(m Mob, l *Level) bool {
	c := l.BossCollisions
	mpv := m.Pos().Vec2
	ppv := l.Player.Pos().Vec2
	if c.CanSee(mpv, ppv, 0.5, 0.5) && m.Detect(mpv.Sub(ppv).Len()) {
		return true
	}
	return false
}

// getPath maps the provided start and end "world" coordinations into discrete
// grid-space, then runs A* search. The resultant slice is also in discrete
// grid-space, since portions of this slice may be thrown away. Calling
// code should therefore map back to "world" coordinates for use when moving.
func getPath(g *twodee.Grid, s, e twodee.Point) []twodee.GridPoint {
	// Need to map from points in the game to locations on the grid board.
	sx, sy := g.GridPosition(s.X(), 0.5), g.GridPosition(s.Y(), 0.5)
	ex, ey := g.GridPosition(e.X(), 0.5), g.GridPosition(e.Y(), 0.5)
	path, err := g.GetPath(sx, sy, ex, ey)
	if err != nil {
		if g.Get(sx, sy) { // Crap, we're inside of something. cheat!
			for _, sp := range neighbors(sx, sy) {
				sx, sy = sp.X, sp.Y
				if !g.Get(sx, sy) {

					//					fmt.Printf("Trying position (%v, %v) to (%v, %v)\n", sx, sy, ex, ey)
					if path, err := g.GetPath(sx, sy, ex, ey); err != nil {
						//						fmt.Printf("Shaken free with %v\n", path)
						return path
					}
				}
			}
			//			fmt.Printf("Still stuck, perhaps we should rand move.\n")
		}
		return []twodee.GridPoint{}
	}
	return path
}

func neighbors(x, y int32) []twodee.GridPoint {
	searchSize := 4
	r := make([]twodee.GridPoint, 0, int(math.Pow(float64(searchSize), 2.0))-1)
	for dx := -searchSize; dx <= searchSize; dx++ {
		for dy := -searchSize; dy <= searchSize; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			r = append(r, twodee.GridPoint{int32(dx), int32(dy)})
		}
	}
	for i, o := range r {
		r[i] = twodee.GridPoint{x + o.X, y + o.Y}
	}
	return r
}

func MaxInt(x, y int) int {
	if x >= y {
		return x
	}
	return y
}
