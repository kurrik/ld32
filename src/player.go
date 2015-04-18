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
	Left
	Right
	Up
	Down
)

var PlayerAnimations = map[PlayerState][]int{
	Standing | Up:    []int{0},
	Standing | Down:  []int{1},
	Standing | Left:  []int{2},
	Standing | Right: []int{3},
	Walking | Up:     []int{4, 5},
	Walking | Down:   []int{6, 7},
	Walking | Left:   []int{8, 9},
	Walking | Right:  []int{8, 9},
	Rolling | Up:     []int{10, 11, 12, 13},
	Rolling | Down:   []int{10, 11, 12, 13},
	Rolling | Left:   []int{10, 11, 12, 13},
	Rolling | Right:  []int{10, 11, 12, 13},
}

type Player struct {
	*twodee.AnimatingEntity
	Dx    float32
	Dy    float32
	speed float32
	State PlayerState
}

func NewPlayer() *Player {
	return &Player{
		AnimatingEntity: twodee.NewAnimatingEntity(
			0, 0, 1, 1, 0,
			twodee.Step10Hz,
			PlayerAnimations[Standing|Up],
		),
		Dx:    0.0,
		Dy:    0.0,
		speed: 0.05,
		State: Standing | Up,
	}
}

func (p *Player) SpriteConfig(sheet *twodee.Spritesheet) twodee.SpriteConfig {
	var (
		frame          = sheet.GetFrame(fmt.Sprintf("numbered_squares_%02d", p.Frame()))
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

func (p *Player) Update(elapsed time.Duration) {
	if p.Dx != 0 || p.Dy != 0 {
		var (
			magX = math.Abs(float64(p.Dx))
			magY = math.Abs(float64(p.Dy))
		)
		p.swapState(Rolling|Standing, Walking)
		if magX > magY {
			if p.Dx > 0 {
				p.swapState(Left|Up|Down, Right)
			} else {
				p.swapState(Up|Right|Down, Left)
			}
		} else {
			if p.Dy > 0 {
				p.swapState(Left|Right|Down, Up)
			} else {
				p.swapState(Left|Up|Right, Down)
			}
		}
		p.move(mgl32.Vec2{p.Dx, p.Dy})
	} else {
		p.swapState(Rolling|Walking, Standing)
	}
	p.AnimatingEntity.Update(elapsed)
}

func (p *Player) move(vec mgl32.Vec2) {
	pos := p.Pos()
	vec = vec.Normalize().Mul(p.speed)
	p.MoveTo(twodee.Pt(pos.X+vec[0], pos.Y+vec[1]))
}

func (p *Player) Roll() {
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
		fmt.Printf("Setting new state: %v\n", state)
		p.State = state
		fmt.Printf("Frames: %v\n", PlayerAnimations)
		if frames, ok := PlayerAnimations[p.State]; ok {
			fmt.Printf("Setting frames: %v\n", frames)
			p.SetFrames(frames)
		}
	}
}
