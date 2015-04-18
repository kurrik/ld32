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
	"io/ioutil"
	"time"
)

const (
	PxPerUnit = 32
)

type GameLayer struct {
	cameraBounds  twodee.Rectangle
	camera        *twodee.Camera
	sprite        *twodee.SpriteRenderer
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
		level:        NewLevel(),
	}
	err = layer.Reset()
	return
}

func (l *GameLayer) Reset() (err error) {
	l.Delete()
	if l.sprite, err = twodee.NewSpriteRenderer(l.camera); err != nil {
		return
	}
	if err = l.loadSpritesheet(); err != nil {
		return
	}
	l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayMusic))
	return
}

func (l *GameLayer) Delete() {
	if l.sprite != nil {
		l.sprite.Delete()
		l.sprite = nil
	}
	if l.spritetexture != nil {
		l.spritetexture.Delete()
		l.spritetexture = nil
	}
}

func (l *GameLayer) Render() {
	if l.level != nil {
		l.spritetexture.Bind()
		l.sprite.Draw([]twodee.SpriteConfig{
			l.level.Player.SpriteConfig(l.spritesheet),
		})
		l.spritetexture.Unbind()
	}
}

func (l *GameLayer) Update(elapsed time.Duration) {
	if l.level != nil {
		l.level.Update(elapsed)
	}
}

func (l *GameLayer) HandleEvent(evt twodee.Event) bool {
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		break
	case *twodee.MouseButtonEvent:
		break
	case *twodee.KeyEvent:
		if event.Type == twodee.Release {
			break
		}
		switch event.Code {
		case twodee.KeyEscape:
			l.app.State.Exit = true
		}
	}
	return true
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
