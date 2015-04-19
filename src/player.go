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
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"time"
)

type PlayerState int32

const (
	_                    = iota
	Standing PlayerState = 1 << iota
	Walking
	Rolling
	Dying
	Left
	Right
	Up
	Down
)

var PlayerAnimations = map[PlayerState][]int{
	Standing | Up:    []int{6},
	Standing | Down:  []int{0},
	Standing | Left:  []int{3},
	Standing | Right: []int{3},
	Walking | Up:     []int{7, 6, 8, 6},
	Walking | Down:   []int{1, 0, 2, 0},
	Walking | Left:   []int{4, 3, 5, 3},
	Walking | Right:  []int{4, 3, 5, 3},
	Rolling | Up:     []int{19, 20, 21, 22, 23},
	Rolling | Down:   []int{14, 15, 16, 17, 18},
	Rolling | Left:   []int{9, 10, 11, 12, 13},
	Rolling | Right:  []int{9, 10, 11, 12, 13},
	Dying:            []int{24, 24, 24, 25, 25, 25, 26, 26, 26, 27, 27, 27, 27, 27, 27},
}

type Player struct {
	*twodee.AnimatingEntity
	events    *twodee.GameEventHandler
	dx        float32
	dy        float32
	rolldx    float32
	rolldy    float32
	speed     float32
	rollspeed float32
	rolling   bool
	State     PlayerState
	Dead      bool
}

func NewPlayer(events *twodee.GameEventHandler, sheet *twodee.Spritesheet) *Player {
	var (
		frame = sheet.GetFrame("player_00")
	)
	return &Player{
		AnimatingEntity: twodee.NewAnimatingEntity(
			0, 0,
			frame.Width, frame.Height,
			0.0,
			twodee.Step10Hz,
			PlayerAnimations[Standing|Down],
		),
		events:    events,
		dx:        0.0,
		dy:        0.0,
		speed:     0.05,
		rollspeed: 0.10,
		rolling:   false,
		Dead:      false,
		State:     Standing | Up,
	}
}

func (p *Player) SpriteConfig(sheet *twodee.Spritesheet) twodee.SpriteConfig {
	var (
		frame          = sheet.GetFrame(fmt.Sprintf("player_%02d", p.Frame()))
		pt             = p.Pos()
		scaleX float32 = 1.0
	)
	if p.State&Left == Left {
		scaleX = -1.0
	}
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			pt.X, pt.Y, 0,
			0, 0, 0,
			scaleX, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}

func (p *Player) Bottom() float32 {
	return p.AnimatingEntity.Bounds().Min.Y
}

func (p *Player) HandleCollision(player *Player) {
}

func (p *Player) UpdateLevel(elapsed time.Duration, level *Level) {
	var (
		isMoving = p.dx != 0 || p.dy != 0
	)
	if !p.Dead {
		if !p.rolling && isMoving {
			var (
				magX = math.Abs(float64(p.dx))
				magY = math.Abs(float64(p.dy))
			)
			p.swapState(Rolling|Standing, Walking)
			if magX > magY {
				if p.dx > 0 {
					p.swapState(Left|Up|Down, Right)
				} else {
					p.swapState(Up|Right|Down, Left)
				}
			} else {
				if p.dy > 0 {
					p.swapState(Left|Right|Down, Up)
				} else {
					p.swapState(Left|Up|Right, Down)
				}
			}
			p.move(mgl32.Vec2{p.dx, p.dy}.Normalize().Mul(p.speed), level)
		} else if p.rolling && isMoving {
			p.swapState(Walking|Standing, Rolling)
			p.move(mgl32.Vec2{p.rolldx, p.rolldy}.Normalize().Mul(p.rollspeed), level)
		} else {
			p.swapState(Rolling|Walking, Standing)
		}
	}
	p.AnimatingEntity.Update(elapsed)
}

func (p *Player) move(vec mgl32.Vec2, level *Level) {
	var (
		bounds = p.Bounds()
		pos    = p.Pos()
	)
	vec = level.Collisions.FixMove(mgl32.Vec4{
		bounds.Min.X,
		bounds.Min.Y,
		bounds.Max.X,
		bounds.Max.Y,
	}, vec, 0.5, 0.5)
	p.MoveTo(twodee.Pt(pos.X+vec[0], pos.Y+vec[1]))
}

func (p *Player) MoveX(mag float32) {
	p.dx = mag
}

func (p *Player) MoveY(mag float32) {
	p.dy = mag
}

func (p *Player) Die() {
	if !p.Dead {
		p.Dead = true
		p.swapState(Left|Right|Down|Up|Walking|Rolling|Standing, Dying)
	}
}

func (p *Player) Roll() {
	if p.rolling {
		return
	}
	p.rolling = true
	p.rolldx = p.dx
	p.rolldy = p.dy
	p.SetCallback(func() {
		p.swapState(Walking|Rolling, Standing)
		p.rolling = false
	})
	p.events.Enqueue(NewShakeEvent(0, 500, 0.08, 4.0, 1.0))
}

func (p *Player) remState(state PlayerState) {
	p.setState(p.State & ^state)
}

func (p *Player) addState(state PlayerState) {
	p.setState(p.State | state)
}

func (p *Player) swapState(rem, add PlayerState) {
	p.setState(p.State & ^rem | add)
}

func (p *Player) setState(state PlayerState) {
	if state != p.State {
		p.State = state
		if frames, ok := PlayerAnimations[p.State]; ok {
			p.SetFrames(frames)
		}
	}
}
