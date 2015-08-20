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
	"image/color"
	"math"
)

type Hud struct {
	blackLine1       *twodee.LineGeometry
	blackLine2       *twodee.LineGeometry
	blackLine3       *twodee.LineGeometry
	levelRedLine     *twodee.LineGeometry
	levelGreenLine   *twodee.LineGeometry
	levelBlueLine    *twodee.LineGeometry
	bossRedLine      *twodee.LineGeometry
	bossGreenLine    *twodee.LineGeometry
	bossBlueLine     *twodee.LineGeometry
	blackStyle       *twodee.LineStyle
	whiteStyle       *twodee.LineStyle
	redStyle         *twodee.LineStyle
	greenStyle       *twodee.LineStyle
	blueStyle        *twodee.LineStyle
	levelRed         float32
	levelGreen       float32
	levelBlue        float32
	levelRedOffset   float32
	levelGreenOffset float32
	levelBlueOffset  float32
	bossRed          float32
	bossGreen        float32
	bossBlue         float32
	bossRedOffset    float32
	bossGreenOffset  float32
	bossBlueOffset   float32
}

func newHud() (hud *Hud, err error) {
	hud = &Hud{
		blackLine1:     twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4.6}, mgl32.Vec2{7.7, 4.6}}, false),
		blackLine2:     twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4.3}, mgl32.Vec2{7.7, 4.3}}, false),
		blackLine3:     twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4}, mgl32.Vec2{7.7, 4}}, false),
		levelRedLine:   twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4.6}, mgl32.Vec2{5.8, 4.6}}, false),
		levelGreenLine: twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4.3}, mgl32.Vec2{5.8, 4.3}}, false),
		levelBlueLine:  twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4}, mgl32.Vec2{5.8, 4}}, false),
		bossRedLine:    twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4.6}, mgl32.Vec2{5.87, 4.6}}, false),
		bossGreenLine:  twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4.3}, mgl32.Vec2{5.87, 4.3}}, false),
		bossBlueLine:   twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4}, mgl32.Vec2{5.87, 4}}, false),
		blackStyle: &twodee.LineStyle{
			Thickness: 0.15,
			Color:     color.RGBA{0, 0, 0, 128},
			Inner:     0.0,
		},
		whiteStyle: &twodee.LineStyle{
			Thickness: 0.15,
			Color:     color.RGBA{255, 255, 255, 128},
			Inner:     0.0,
		},
		redStyle: &twodee.LineStyle{
			Thickness: 0.15,
			Color:     color.RGBA{255, 0, 0, 128},
			Inner:     0.0,
		},
		greenStyle: &twodee.LineStyle{
			Thickness: 0.15,
			Color:     color.RGBA{0, 255, 0, 128},
			Inner:     0.0,
		},
		blueStyle: &twodee.LineStyle{
			Thickness: 0.15,
			Color:     color.RGBA{0, 0, 255, 128},
			Inner:     0.0,
		},
		levelRed:         0.0,
		levelGreen:       0.0,
		levelBlue:        0.0,
		levelRedOffset:   0.0,
		levelGreenOffset: 0.0,
		levelBlueOffset:  0.0,
		bossRed:          0.0,
		bossGreen:        0.0,
		bossBlue:         0.0,
		bossRedOffset:    0.0,
		bossGreenOffset:  0.0,
		bossBlueOffset:   0.0,
	}
	return
}

func (h *Hud) UpdateLines(l *Level, initialLoad bool) {

	// check if level's red value has changed
	if h.levelRed != l.Color[0] {
		// update stored red value
		h.levelRed = float32(math.Min(float64(l.Color[0]), 1.0))
		// update offset
		h.levelRedOffset = (7.7 - 5.8) * h.levelRed
		// update the level's red color line
		h.levelRedLine = twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4.6}, mgl32.Vec2{5.8 + h.levelRedOffset, 4.6}}, false)
	}

	// check if level's green value has changed
	if h.levelGreen != l.Color[1] {
		// update stored green value
		h.levelGreen = float32(math.Min(float64(l.Color[1]), 1.0))
		// update offset
		h.levelGreenOffset = (7.7 - 5.8) * h.levelGreen
		// update the level's green color line
		h.levelGreenLine = twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4.3}, mgl32.Vec2{5.8 + h.levelGreenOffset, 4.3}}, false)
	}

	// check if level's blue value has changed
	if h.levelBlue != l.Color[2] {
		// update stored blue value
		h.levelBlue = float32(math.Min(float64(l.Color[2]), 1.0))
		// update offset
		h.levelBlueOffset = (7.7 - 5.8) * h.levelBlue
		// update the level's blue color line
		h.levelBlueLine = twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4}, mgl32.Vec2{5.8 + h.levelBlueOffset, 4}}, false)
	}

	// update the color markers for the current boss
	if l.Boss != nil {

		if (h.bossRed != l.Boss.Color[0]) || initialLoad {
			// update stored boss red value
			h.bossRed = l.Boss.Color[0]
			// update offset
			h.bossRedOffset = (7.7 - 5.8) * h.bossRed
			// update the boss's red color marker
			h.bossRedLine = twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8 + h.bossRedOffset, 4.6}, mgl32.Vec2{5.87 + h.bossRedOffset, 4.6}}, false)
		}

		if (h.bossGreen != l.Boss.Color[1]) || initialLoad {
			// update stored boss green value
			h.bossGreen = l.Boss.Color[1]
			// update offset
			h.bossGreenOffset = (7.7 - 5.8) * h.bossGreen
			// update the boss's green color marker
			h.bossGreenLine = twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8 + h.bossGreenOffset, 4.3}, mgl32.Vec2{5.87 + h.bossGreenOffset, 4.3}}, false)
		}

		if (h.bossBlue != l.Boss.Color[2]) || initialLoad {
			// update stored boss blue value
			h.bossBlue = l.Boss.Color[2]
			// update offset
			h.bossBlueOffset = (7.7 - 5.8) * h.bossBlue
			// update the boss's blue color marker
			h.bossBlueLine = twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8 + h.bossBlueOffset, 4}, mgl32.Vec2{5.87 + h.bossBlueOffset, 4}}, false)
		}
	} else {
		h.bossRedLine = twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4.6}, mgl32.Vec2{5.87, 4.6}}, false)
		h.bossGreenLine = twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4.3}, mgl32.Vec2{5.87, 4.3}}, false)
		h.bossBlueLine = twodee.NewLineGeometry([]mgl32.Vec2{mgl32.Vec2{5.8, 4}, mgl32.Vec2{5.87, 4}}, false)
	}

	return
}
