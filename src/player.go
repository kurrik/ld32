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
	"time"
)

var PLAYER_FRAMES = []int{0, 1, 2, 3, 4}

type Player struct {
	*twodee.AnimatingEntity
	Dx    float32
	Dy    float32
	Speed float32
}

func NewPlayer() *Player {
	return &Player{
		AnimatingEntity: twodee.NewAnimatingEntity(0, 0, 1, 1, 0, twodee.Step10Hz, PLAYER_FRAMES),
		Dx:              0.0,
		Dy:              0.0,
		Speed:           0.05,
	}
}

func (p *Player) SpriteConfig(sheet *twodee.Spritesheet) twodee.SpriteConfig {
	var (
		frame = sheet.GetFrame(fmt.Sprintf("numbered_squares_%02d", p.Frame()))
		pt    = p.Pos()
	)
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			pt.X, pt.Y, 0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}

func (p *Player) Update(elapsed time.Duration) {
	if p.Dx != 0 || p.Dy != 0 {
		p.move(mgl32.Vec2{p.Dx, p.Dy})
	}
	p.AnimatingEntity.Update(elapsed)
}

func (p *Player) move(vec mgl32.Vec2) {
	pos := p.Pos()
	vec = vec.Normalize().Mul(p.Speed)
	p.MoveTo(twodee.Pt(pos.X+vec[0], pos.Y+vec[1]))
}
