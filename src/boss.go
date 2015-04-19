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
)

var BossMap = map[string]BossMaker{
	"boss1": MakeBoss1,
	"boss2": MakeBoss2,
}

type BossMaker func() *Boss

func MakeBoss1() *Boss {
	return NewBoss(&Mobile{0, 5 * time.Second})
}

func MakeBoss2() *Boss {
	return NewBoss(&Mobile{1, 20 * time.Second})
}

type Boss struct {
	*twodee.AnimatingEntity
	*Mobile
	dx, dy, speed float32
	State         MobState
}

func NewBoss(m *Mobile) *Boss {
	return &Boss{
		AnimatingEntity: twodee.NewAnimatingEntity(
			0, 0, 1, 1, 0,
			twodee.Step10Hz,
			PlayerAnimations[Standing|Up],
		),
		Mobile: m,
		dx:     0.0,
		dy:     0.0,
		speed:  0.04,
		State:  &SearchState{},
	}
}

func (b *Boss) ExamineWorld(l *Level) {
	newState := b.State.ExamineWorld(b, l)
	b.State = newState
}

func (b *Boss) Update(elapsed time.Duration) {
	b.State.Update(b, elapsed)
}
