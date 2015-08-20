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
)

const (
	PlayBackgroundMusic twodee.GameEventType = iota
	PlayBossMusic
	PauseMusic
	ResumeMusic
	PlayBossDeathEffect
	PlayColorChangeEffect
	PlayPlayerDeathEffect
	PlayRollEffect
	ShakeCamera
	ChangeColor
	BossColor
	BossDied
	PlayerDied
	SENTINEL
)

const (
	NumGameEventTypes = int(SENTINEL)
)

type ColorEvent struct {
	twodee.BasicGameEvent
	Color mgl32.Vec3
	Add   bool
}

func NewColorEvent(color mgl32.Vec3, add bool) *ColorEvent {
	return &ColorEvent{
		*twodee.NewBasicGameEvent(ChangeColor),
		color,
		add,
	}
}

type ShakeEvent struct {
	twodee.BasicGameEvent
	Millis    int32
	Amplitude float32
	Frequency float32
	Decay     float32
	Priority  int32
}

func NewShakeEvent(priority, ms int32, amplitude, freq, decay float32) *ShakeEvent {
	return &ShakeEvent{
		BasicGameEvent: *twodee.NewBasicGameEvent(ShakeCamera),
		Millis:         ms,
		Amplitude:      amplitude,
		Frequency:      freq,
		Decay:          decay,
		Priority:       priority,
	}
}

type BossColorEvent struct {
	twodee.BasicGameEvent
	Color mgl32.Vec3
}

func NewBossColorEvent(color mgl32.Vec3) *BossColorEvent {
	return &BossColorEvent{
		*twodee.NewBasicGameEvent(BossColor),
		color,
	}
}

type BossDiedEvent struct {
	twodee.BasicGameEvent
	Name string
}

func NewBossDiedEvent(name string) *BossDiedEvent {
	return &BossDiedEvent{
		BasicGameEvent: *twodee.NewBasicGameEvent(BossDied),
		Name:           name,
	}
}

type PlayerDiedEvent struct {
	twodee.BasicGameEvent
}

func NewPlayerDiedEvent() *PlayerDiedEvent {
	return &PlayerDiedEvent{
		*twodee.NewBasicGameEvent(PlayerDied),
	}
}
