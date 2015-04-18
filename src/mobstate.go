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
	"time"

	"../lib/twodee"
	"github.com/go-gl/mathgl/mgl32"
)

type Mob interface {
	Bored(time.Duration) bool
	SetFrames(f []int)
	Detect(dist float32) bool
	Pos() twodee.Point
}

type Mobile struct {
	DetectionRadius float32
	BoredThreshold  time.Duration
}

func (m *Mobile) Bored(d time.Duration) bool {
	return d >= m.BoredThreshold
}

func (m *Mobile) Detect(d float32) bool {
	return d <= m.DetectionRadius
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

// SearchState is the state during which a mobile is aimlessly wandering,
// hoping to chance across the player.
type SearchState struct{}

// ExamineWorld returns HuntState if the player is seen, otherwise the mob
// continues wandering.
func (s *SearchState) ExamineWorld(m Mob, l *Level) MobState {
	if playerSeen(m, l) {
		return &HuntState{}
	}
	return s
}

func (s *SearchState) Update(m Mob, d time.Duration) {
}

func (s *SearchState) Enter(m Mob) {
	// TODO: set some hunting animation.
}

func (s *SearchState) Exit(m Mob) {
	// TODO: maybe something should happen when we start hunting?
}

// HuntState is the state during which a mobile is actively hunting the player.
type HuntState struct {
	durSinceLastContact time.Duration
}

// ExamineWorld returns the current state if the player is currently seen or
// the mob is not yet tired of chasing. Otherwise, it returns nil.
func (h *HuntState) ExamineWorld(m Mob, l *Level) MobState {
	if playerSeen(m, l) {
		h.durSinceLastContact = time.Duration(0)
		return h
	}
	if !m.Bored(h.durSinceLastContact) {
		return h
	}
	return nil
}

// Update resets the player's hiding timer if the player is seen, otherwise it
// increments.
func (h *HuntState) Update(m Mob, d time.Duration) {
	h.durSinceLastContact += d
}

func (h *HuntState) Enter(m Mob) {
}

func (h *HuntState) Exit(m Mob) {
}

// playerSeen returns true if the player is currently visible to the mob and
// within its detection radius.
func playerSeen(m Mob, l *Level) bool {
	c := l.Collisions
	mpv := mgl32.Vec2{m.Pos().X, m.Pos().Y}
	ppv := mgl32.Vec2{l.Player.Pos().X, l.Player.Pos().Y}
	if c.CanSee(mpv, ppv, 0.5, 0.5) && m.Detect(mpv.Sub(ppv).Len()) {
		return true
	}
	return false
}
