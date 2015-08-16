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
	"../lib/twodee/twodee"
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

type Plate struct {
	Prop
	Color   mgl32.Vec4
	Active  bool
	events  *twodee.GameEventHandler
	elapsed time.Duration
}

func NewPlate(x, y float32, color mgl32.Vec3, sheet *twodee.Spritesheet, events *twodee.GameEventHandler) *Plate {
	return &Plate{
		Prop: NewStaticProp(
			x, y,
			sheet,
			"plate.fw",
		),
		Color:  color.Vec4(1.0),
		events: events,
	}
}

func (p *Plate) SpriteConfig(sheet *twodee.Spritesheet) twodee.SpriteConfig {
	c := p.Prop.SpriteConfig(sheet)
	if !p.Active {
		c.Color = p.Color
	} else {
		c.Color = mgl32.Vec4{0.0, 0.0, 0.0, 1.0}
	}
	return c
}

func (p *Plate) Update(elapsed time.Duration) {
	if !p.Active {
		return
	}
	p.elapsed += elapsed
	if p.elapsed > time.Duration(5)*time.Second {
		p.events.Enqueue(NewColorEvent(p.Color.Vec3(), false))
		p.Active = false
	}
}

func (p *Plate) HandleCollision(player *Player) {
	if !p.Active {
		p.Active = true
		p.events.Enqueue(NewColorEvent(p.Color.Vec3(), true))
		p.elapsed = time.Duration(0)
	}
}
