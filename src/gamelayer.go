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
	"time"
)

type GameLayer struct {
	cameraBounds twodee.Rectangle
	camera       *twodee.Camera
	sprite       *twodee.SpriteRenderer
	app          *Application
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
	if l.sprite != nil {
		l.sprite.Delete()
	}
	if l.sprite, err = twodee.NewSpriteRenderer(l.camera); err != nil {
		return
	}
	return
}

func (l *GameLayer) Delete() {
	if l.sprite != nil {
		l.sprite.Delete()
		l.sprite = nil
	}
}

func (l *GameLayer) Render() {
}

func (l *GameLayer) Update(elapsed time.Duration) {
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
