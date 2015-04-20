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
	"fmt"
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"../lib/twodee"
)

type BossState int32

const (
	_                = iota
	Normal BossState = 1 << iota
	BossDying
)

var BossAnimations = map[BossState][]int{
	Normal:    []int{0, 1},
	BossDying: []int{0, 1, 2, 3, 4, 5, 5},
}

var BossMap = map[string]BossMaker{
	"boss1": MakeBoss1,
	"boss2": MakeBoss2,
}

type BossMaker func(x, y float32, events *twodee.GameEventHandler) *Boss

// MakeBoss1 returns a boss that searches left and right and gets bored easily.
func MakeBoss1(x, y float32, events *twodee.GameEventHandler) *Boss {
	sp := []mgl32.Vec2{
		mgl32.Vec2{x - 5, y},
		mgl32.Vec2{x + 5, y},
	}
	return NewBoss(&Mobile{
		DetectionRadius: 4,
		BoredThreshold:  5 * time.Second,
		speed:           0.04,
		searchPattern:   sp,
	}, []mgl32.Vec3{
		mgl32.Vec3{1.0, 0.0, 0.0},
		mgl32.Vec3{0.0, 1.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 1.0},
	}, events)
}

func MakeBoss2(x, y float32, events *twodee.GameEventHandler) *Boss {
	return NewBoss(&Mobile{
		DetectionRadius: 10,
		BoredThreshold:  20 * time.Second,
		speed:           0.04,
		searchPattern:   []mgl32.Vec2{},
	}, []mgl32.Vec3{
		mgl32.Vec3{1.0, 0.0, 0.0},
		mgl32.Vec3{0.0, 1.0, 0.0},
		mgl32.Vec3{0.0, 0.0, 1.0},
	}, events)
}

type Boss struct {
	*twodee.AnimatingEntity
	*Mobile
	// Likely don't need speed anymore on Boss, since it's on Mobile.
	dx, dy, speed float32
	StateStack    []MobState
	Color         mgl32.Vec3
	Colors        []mgl32.Vec3
	events        *twodee.GameEventHandler
	Dead          bool
}

func NewBoss(m *Mobile, colors []mgl32.Vec3, events *twodee.GameEventHandler) *Boss {
	b := &Boss{
		AnimatingEntity: twodee.NewAnimatingEntity(
			0, 0, 1, 1, 0,
			twodee.Step5Hz,
			BossAnimations[Normal],
		),
		Mobile:     m,
		dx:         0.0,
		dy:         0.0,
		speed:      0.04,
		StateStack: []MobState{NewVegState()},
		Colors:     colors,
		events:     events,
		Dead:       false,
	}
	b.NextColor()
	return b
}

func (b *Boss) NextColor() {
	if len(b.Colors) > 0 {
		b.Color = b.Colors[0]
		b.Colors = b.Colors[1:]
		b.events.Enqueue(NewBossColorEvent(b.Color))
	} else {
		b.events.Enqueue(NewBossDiedEvent())
	}
}

func (b *Boss) ExamineWorld(l *Level) {
	cState := b.StateStack[len(b.StateStack)-1]
	newState := cState.ExamineWorld(b, l)
	if newState == cState {
		return
	}
	if newState == nil { // Transition to last state.
		cState.Exit(b)
		b.StateStack = b.StateStack[:len(b.StateStack)-1]
		b.StateStack[len(b.StateStack)-1].Enter(b)
		return
	}
	cState.Exit(b)
	b.StateStack = append(b.StateStack, newState)
	newState.Enter(b)
}

func (b *Boss) Update(elapsed time.Duration) {
	b.AnimatingEntity.Update(elapsed)
	if !b.Dead {
		// Hrm, should update be fed to every state in the stack?
		for i := len(b.StateStack) - 1; i >= 0; i-- {
			b.StateStack[i].Update(b, elapsed)
		}
		//	b.StateStack[len(b.StateStack)-1].Update(b, elapsed)
	}
}

func (b *Boss) Bottom() float32 {
	return b.AnimatingEntity.Bounds().Min.Y
}

func (b *Boss) Die() {
	if !b.Dead {
		b.SetFrames(BossAnimations[BossDying])
		b.Dead = true
		b.events.Enqueue(twodee.NewBasicGameEvent(PlayBossDeathEffect))
	}
}

func (b *Boss) SpriteConfig(sheet *twodee.Spritesheet) twodee.SpriteConfig {
	frame := sheet.GetFrame(fmt.Sprintf("boss_%02d", b.Frame()))
	pt := b.Pos()
	scaleX := float32(1.0)
	// Implement facing left...
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			pt.X, pt.Y, 0,
			0, 0, 0,
			scaleX, 1.0, 1.0,
		},
		Frame: frame.Frame,
		Color: b.Color.Vec4(1.0),
	}
}

func (b *Boss) ShouldSwing(p mgl32.Vec2) bool {
	bv := mgl32.Vec2{b.Pos().X, b.Pos().Y}
	return p.Sub(bv).Len() < 1
}
