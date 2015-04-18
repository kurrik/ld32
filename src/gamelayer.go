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
	twodee "../lib/twodee"
	"github.com/go-gl/mathgl/mgl32"
	"io/ioutil"
	"time"
)

const (
	PxPerUnit = 32
)

type GameLayer struct {
	shake         *twodee.ContinuousAnimation
	cameraBounds  twodee.Rectangle
	camera        *twodee.Camera
	sprite        *twodee.SpriteRenderer
	batch         *twodee.BatchRenderer
	effects       *EffectsRenderer
	app           *Application
	spritesheet   *twodee.Spritesheet
	spritetexture *twodee.Texture
	level         *Level
}

func NewGameLayer(winb twodee.Rectangle, app *Application) (layer *GameLayer, err error) {
	var (
		camera       *twodee.Camera
		cameraBounds = twodee.Rect(-10, -10, 10, 10)
	)
	if camera, err = twodee.NewCamera(cameraBounds, winb); err != nil {
		return
	}
	layer = &GameLayer{
		camera:       camera,
		cameraBounds: cameraBounds,
		app:          app,
	}
	err = layer.Reset()
	return
}

func (l *GameLayer) Reset() (err error) {
	l.Delete()
	if l.batch, err = twodee.NewBatchRenderer(l.camera); err != nil {
		return
	}
	if l.sprite, err = twodee.NewSpriteRenderer(l.camera); err != nil {
		return
	}
	if l.effects, err = NewEffectsRenderer(320, 320, 1.0); err != nil {
		return
	}
	if err = l.loadSpritesheet(); err != nil {
		return
	}
	if l.level, err = NewLevel("resources/background.tmx", l.spritesheet); err != nil {
		return
	}
	l.updateCamera(1.0)
	l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayMusic))
	return
}

func (l *GameLayer) Delete() {
	if l.batch != nil {
		l.batch.Delete()
		l.batch = nil
	}
	if l.sprite != nil {
		l.sprite.Delete()
		l.sprite = nil
	}
	if l.spritetexture != nil {
		l.spritetexture.Delete()
		l.spritetexture = nil
	}
	if l.effects != nil {
		l.effects.Delete()
		l.effects = nil
	}
}

func (l *GameLayer) Render() {
	if l.level != nil {
		l.effects.Bind()
		l.batch.Bind()
		if err := l.batch.Draw(l.level.Background, 0, 0, 0); err != nil {
			panic(err)
		}
		l.batch.Unbind()
		l.spritetexture.Bind()
		l.sprite.Draw(l.level.Props.SpriteConfigs(l.spritesheet))
		l.spritetexture.Unbind()
		l.effects.Unbind()
		l.effects.Draw()
	}
}

func (l *GameLayer) Update(elapsed time.Duration) {
	l.checkKeys()
	if l.shake != nil {
		l.shake.Update(elapsed)
	}
	l.updateCamera(0.05)
	if l.level != nil {
		l.level.Update(elapsed)
	}
}

func (l *GameLayer) updateCamera(scale float32) {
	var (
		pPt     = l.level.Player.Pos()
		cRect   = l.camera.WorldBounds
		cWidth  = cRect.Max.X - cRect.Min.X
		cHeight = cRect.Max.Y - cRect.Min.Y
		cMidX   = cRect.Min.X + (cWidth / 2.0)
		cMidY   = cRect.Min.Y + (cHeight / 2.0)
		pVec    = mgl32.Vec2{pPt.X, pPt.Y}
		cVec    = mgl32.Vec2{cMidX, cMidY}
		diff    = pVec.Sub(cVec)
		bounds  twodee.Rectangle
		adj     mgl32.Vec2
	)
	if diff.Len() > 1 {
		adj = diff.Mul(scale)
	} else {
		adj = mgl32.Vec2{0, 0}
	}
	if l.shake != nil {
		adj[1] += l.shake.Value()
	}
	bounds = twodee.Rect(
		cRect.Min.X+adj[0],
		cRect.Min.Y+adj[1],
		cRect.Max.X+adj[0],
		cRect.Max.Y+adj[1],
	)
	l.camera.SetWorldBounds(bounds)
}

func (l *GameLayer) ShakeCamera() {
	if l.shake == nil {
		decay := twodee.SineDecayFunc(
			time.Duration(500)*time.Millisecond,
			0.08, // Amplitude
			4.0,  // Frequency
			1.0,  // Decay
		)
		l.shake = twodee.NewContinuousAnimation(decay)
	}
	l.shake.Reset()
}

func (l *GameLayer) HandleEvent(evt twodee.Event) bool {
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		break
	case *twodee.MouseButtonEvent:
		break
	case *twodee.KeyEvent:
		//l.handleMovement(event)
		if event.Type == twodee.Release {
			break
		}
		switch event.Code {
		case twodee.KeyX:
			l.app.State.Exit = true
		case twodee.KeyZ:
			l.level.Player.Roll()
			l.ShakeCamera()
		case twodee.KeyM:
			if twodee.MusicIsPaused() {
				l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(ResumeMusic))
			} else {
				l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PauseMusic))
			}
		case twodee.KeyC:
			l.effects.Color = l.effects.Color.Add(mgl32.Vec3{0.1, 0.0, 0.0})
		}

	}
	return true
}

func (l *GameLayer) checkKeys() {
	var (
		events         = l.app.Context.Events
		down           = events.GetKey(twodee.KeyDown) == twodee.Press
		up             = events.GetKey(twodee.KeyUp) == twodee.Press
		left           = events.GetKey(twodee.KeyLeft) == twodee.Press
		right          = events.GetKey(twodee.KeyRight) == twodee.Press
		x      float32 = 0.0
		y      float32 = 0.0
	)
	switch {
	case down && !up:
		y = -1.0
	case up && !down:
		y = 1.0
	}
	switch {
	case left && !right:
		x = -1.0
	case right && !left:
		x = 1.0
	}
	l.level.Player.MoveX(x)
	l.level.Player.MoveY(y)
}

func (l *GameLayer) loadSpritesheet() (err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile("resources/spritesheet.json"); err != nil {
		return
	}
	if l.spritesheet, err = twodee.ParseTexturePackerJSONArrayString(
		string(data),
		PxPerUnit,
	); err != nil {
		return
	}
	if l.spritetexture, err = twodee.LoadTexture(
		"resources/"+l.spritesheet.TexturePath,
		twodee.Nearest,
	); err != nil {
		return
	}
	return
}
