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
	"sort"
	"time"
)

type Prop interface {
	SpriteConfig(sheet *twodee.Spritesheet) twodee.SpriteConfig
	Bottom() float32
	HandleCollision(p *Player)
	Bounds() twodee.Rectangle
	Update(elapsed time.Duration)
}

type StaticProp struct {
	twodee.Entity
	Name string
}

func NewStaticProp(x, y float32, sheet *twodee.Spritesheet, name string) *StaticProp {
	var (
		frame = sheet.GetFrame(name) // Ugh
	)
	return &StaticProp{
		Entity: twodee.NewBaseEntity(x, y, frame.Width, frame.Height, 0.0, 0),
		Name:   name,
	}
}

func (p *StaticProp) SpriteConfig(sheet *twodee.Spritesheet) twodee.SpriteConfig {
	var (
		pos = p.Entity.Pos()
	)
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			pos.X, pos.Y, 0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: sheet.GetFrame(p.Name).Frame,
	}
}

func (p *StaticProp) Bottom() float32 {
	return p.Entity.Bounds().Min.Y
}

func (p *StaticProp) HandleCollision(player *Player) {
}

func (p *StaticProp) Update(elapsed time.Duration) {
}

type PropList []Prop

func NewPropList() PropList {
	return make(PropList, 0)
}

func (l PropList) Len() int {
	return len(l)
}

func (l PropList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l PropList) Less(i, j int) bool {
	return l[i].Bottom() > l[j].Bottom() // Top to bottom
}

func (l PropList) SpriteConfigs(sheet *twodee.Spritesheet) (out []twodee.SpriteConfig) {
	out = make([]twodee.SpriteConfig, l.Len())
	sort.Sort(l)
	for i, prop := range l {
		out[i] = prop.SpriteConfig(sheet)
	}
	return
}

func (l PropList) CheckCollision(p *Player) {
	bounds := p.Bounds()
	for _, prop := range l {
		if prop.Bounds().Overlaps(bounds) {
			prop.HandleCollision(p)
		}
	}
}

func (l PropList) Update(elapsed time.Duration) {
	for _, prop := range l {
		prop.Update(elapsed)
	}
}
