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
	"image/color"
	"io/ioutil"
	"math"
	"time"
)

const (
	PxPerUnit = 32
)

type GameLayer struct {
	levels               map[string]string
	shake                *twodee.ContinuousAnimation
	cameraBounds         twodee.Rectangle
	camera               *twodee.Camera
	linesCamera          *twodee.Camera
	sprite               *twodee.SpriteRenderer
	lines                *twodee.LinesRenderer
	debugLines           *twodee.LinesRenderer
	batch                *twodee.BatchRenderer
	effects              *EffectsRenderer
	app                  *Application
	spritesheet          *twodee.Spritesheet
	spritetexture        *twodee.Texture
	level                *Level
	hud                  *Hud
	splash               string
	shakeObserverId      int
	shakePriority        int32
	bossDiedObserverId   int
	playerDiedObserverId int
}

func NewGameLayer(winb twodee.Rectangle, app *Application) (layer *GameLayer, err error) {
	var (
		camera       *twodee.Camera
		linesCamera  *twodee.Camera
		cameraBounds = twodee.Rect(-8, -5, 8, 5)
		hud          *Hud
	)
	if camera, err = twodee.NewCamera(cameraBounds, winb); err != nil {
		return
	}
	if linesCamera, err = twodee.NewCamera(cameraBounds, winb); err != nil {
		return
	}
	if hud, err = newHud(); err != nil {
		return
	}
	layer = &GameLayer{
		camera:       camera,
		linesCamera:  linesCamera,
		cameraBounds: cameraBounds,
		app:          app,
		levels: map[string]string{
			"main":  "resources/main.tmx",
			"boss1": "resources/boss1.tmx",
			"boss2": "resources/boss2.tmx",
		},
		shakePriority: -1,
		hud:           hud,
		splash:        "splash",
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
	if l.lines, err = twodee.NewLinesRenderer(l.linesCamera); err != nil {
		return
	}
	if l.debugLines, err = twodee.NewLinesRenderer(l.camera); err != nil {
		return
	}
	if l.effects, err = NewEffectsRenderer(512, 320, 1.0); err != nil {
		return
	}
	if err = l.loadSpritesheet(); err != nil {
		return
	}
	l.shakeObserverId = l.app.GameEventHandler.AddObserver(ShakeCamera, l.shakeCamera)
	l.bossDiedObserverId = l.app.GameEventHandler.AddObserver(BossDied, l.bossDied)
	l.playerDiedObserverId = l.app.GameEventHandler.AddObserver(PlayerDied, l.playerDied)
	l.loadLevel("main")
	l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayBackgroundMusic))
	return
}

func (l *GameLayer) loadLevel(name string) (err error) {
	var (
		path string
		ok   bool
	)
	if path, ok = l.levels[name]; !ok {
		return fmt.Errorf("Invalid level: %v", name)
	}
	if l.level != nil {
		l.level.Delete()
	}
	if l.level, err = NewLevel(name, path, l.spritesheet, l.app.GameEventHandler); err != nil {
		return
	}
	l.updateCamera(1.0)
	// check name of level being loaded
	// if "main" then trigger PlayBackgroundMusic event
	if name == "main" {
		l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayBackgroundMusic))
	} else {
		// else trigger PlayBossMusic event
		l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayBossMusic))
		l.hud.UpdateLines(l.level, true)
	}
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
	if l.lines != nil {
		l.lines.Delete()
		l.lines = nil
	}
	if l.debugLines != nil {
		l.debugLines.Delete()
		l.debugLines = nil
	}
	if l.effects != nil {
		l.effects.Delete()
		l.effects = nil
	}
	if l.shakeObserverId != 0 {
		l.app.GameEventHandler.RemoveObserver(ShakeCamera, l.shakeObserverId)
	}
	if l.bossDiedObserverId != 0 {
		l.app.GameEventHandler.RemoveObserver(BossDied, l.bossDiedObserverId)
	}
	if l.playerDiedObserverId != 0 {
		l.app.GameEventHandler.RemoveObserver(PlayerDied, l.playerDiedObserverId)
	}
	if l.level != nil {
		l.level.Delete()
		l.level = nil
	}
}

func (l *GameLayer) Render() {
	if l.splash != "" {
		l.spritetexture.Bind()
		splash := []twodee.SpriteConfig{l.getSplashSpriteConfig(l.splash, l.camera)}
		l.sprite.Draw(splash)
		l.spritetexture.Unbind()
	} else if l.level != nil {
		if l.level.Player.Dead {
			l.spritetexture.Bind()
			l.sprite.Draw([]twodee.SpriteConfig{l.level.Player.SpriteConfig(l.spritesheet)})
			l.spritetexture.Unbind()
		} else if l.level.Boss != nil && l.level.Boss.Dead {
			l.spritetexture.Bind()
			l.sprite.Draw([]twodee.SpriteConfig{l.level.Boss.SpriteConfig(l.spritesheet)})
			l.spritetexture.Unbind()
		} else {
			l.effects.Bind()
			l.batch.Bind()
			if err := l.batch.Draw(l.level.Background, 0, 0, 0); err != nil {
				panic(err)
			}
			l.batch.Unbind()
			l.spritetexture.Bind()
			if len(l.level.Plates) > 0 {
				l.sprite.Draw(l.level.Plates.SpriteConfigs(l.spritesheet))
			}
			if len(l.level.Props) > 0 {
				l.sprite.Draw(l.level.Props.SpriteConfigs(l.spritesheet))
			}
			l.spritetexture.Unbind()
			l.effects.Unbind()
			l.effects.Draw()

			modelview := mgl32.Ident4()

			l.hud.UpdateLines(l.level, false)

			l.lines.Bind()
			l.lines.Draw(l.hud.blackLine1, modelview, l.hud.blackStyle)
			l.lines.Draw(l.hud.blackLine2, modelview, l.hud.blackStyle)
			l.lines.Draw(l.hud.blackLine3, modelview, l.hud.blackStyle)
			l.lines.Draw(l.hud.levelRedLine, modelview, l.hud.redStyle)
			l.lines.Draw(l.hud.levelGreenLine, modelview, l.hud.greenStyle)
			l.lines.Draw(l.hud.levelBlueLine, modelview, l.hud.blueStyle)
			l.lines.Draw(l.hud.bossRedLine, modelview, l.hud.whiteStyle)
			l.lines.Draw(l.hud.bossGreenLine, modelview, l.hud.whiteStyle)
			l.lines.Draw(l.hud.bossBlueLine, modelview, l.hud.whiteStyle)
			l.lines.Unbind()

			if l.app.State.Debug {
				l.drawBossLines()
			}
		}
	}
}

func (l *GameLayer) drawBossLines() {
	if l.level.Boss == nil || len(l.level.BossPath) == 0 {
		return
	}
	var (
		points = make([]mgl32.Vec2, len(l.level.BossPath))
		geom   *twodee.LineGeometry
		style  = &twodee.LineStyle{
			Thickness: 0.15,
			Color:     color.RGBA{255, 0, 255, 128},
		}
	)
	for i, gridPoint := range l.level.BossPath {
		points[i] = mgl32.Vec2{
			l.level.BossCollisions.InversePosition(gridPoint.X, 0.5),
			l.level.BossCollisions.InversePosition(gridPoint.Y, 0.5),
		}
	}
	geom = twodee.NewLineGeometry(points, false)
	l.debugLines.Bind()
	l.debugLines.Draw(geom, mgl32.Ident4(), style)
	l.debugLines.Unbind()
}

func (l *GameLayer) Update(elapsed time.Duration) {
	if l.splash != "" {
		return
	}
	if !l.checkJoy() {
		l.checkKeys()
	}
	if l.shake != nil {
		l.shake.Update(elapsed)
	}
	l.updateCamera(0.05)
	if l.level != nil {
		l.level.Update(elapsed)
		if collides, level := l.level.PortalCollides(); collides {
			l.loadLevel(level)
		}
	}
	l.effects.Color = l.level.Color
}

func (l *GameLayer) updateCamera(scale float32) {
	if l.level.Player.Dead || (l.level.Boss != nil && l.level.Boss.Dead) {
		return
	}
	var (
		pPt     = l.level.Player.Pos()
		cRect   = l.camera.WorldBounds
		cWidth  = cRect.Max.X() - cRect.Min.X()
		cHeight = cRect.Max.Y() - cRect.Min.Y()
		cMidX   = cRect.Min.X() + (cWidth / 2.0)
		cMidY   = cRect.Min.Y() + (cHeight / 2.0)
		pVec    = pPt.Vec2
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
		cRect.Min.X()+adj[0],
		cRect.Min.Y()+adj[1],
		cRect.Max.X()+adj[0],
		cRect.Max.Y()+adj[1],
	)
	l.camera.SetWorldBounds(bounds)
}

func (l *GameLayer) shakeCamera(e twodee.GETyper) {
	if event, ok := e.(*ShakeEvent); ok {
		if l.shake == nil || event.Priority > l.shakePriority {
			decay := twodee.SineDecayFunc(
				time.Duration(event.Millis)*time.Millisecond,
				event.Amplitude,
				event.Frequency,
				event.Decay,
				func() {
					l.shake = nil
					l.shakePriority = -1
				},
			)
			l.shake = twodee.NewContinuousAnimation(decay)
			l.shakePriority = event.Priority
		}
	}
}

func (l *GameLayer) bossDied(e twodee.GETyper) {
	if event, ok := e.(*BossDiedEvent); ok {
		if l.level.Boss != nil && !l.level.Boss.Dead {
			bounds := l.camera.WorldBounds
			pt := l.level.Boss.Pos()
			midpoint := bounds.Midpoint()
			adjx := pt.X() - midpoint.X()
			adjy := pt.Y() - midpoint.Y()
			bounds.Min.Vec2[0] += adjx
			bounds.Max.Vec2[0] += adjx
			bounds.Min.Vec2[1] += adjy
			bounds.Max.Vec2[1] += adjy
			l.camera.SetWorldBounds(bounds)
			l.level.Boss.Die()
			l.level.Boss.SetCallback(func() {
				l.checkBosses(event.Name)
				l.loadLevel("main")
			})
		}
	}
}

func (l *GameLayer) playerDied(e twodee.GETyper) {
	if l.app.State.Debug {
		fmt.Printf("Player died\n")
	}
	l.level.Player.Die()
	l.level.Player.SetCallback(func() {
		if l.app.State.Debug {
			fmt.Printf("Done animating death\n")
		}
		l.loadLevel("main")
	})
}

func (l *GameLayer) checkBosses(name string) {
	var (
		ok     bool
		boss   string
		needed = []string{"boss1", "boss2"}
	)
	l.app.State.KilledBosses[name] = true
	for _, boss = range needed {
		if _, ok = l.app.State.KilledBosses[boss]; !ok {
			return
		}
	}
	l.splash = "won"
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
		if l.splash != "" {
			if l.splash == "won" {
				l.app.State.Exit = true
			}
			l.splash = ""
			return false
		}
		switch event.Code {
		case twodee.KeyZ:
			l.level.Player.Roll()
		case twodee.KeyM:
			if twodee.MusicIsPaused() {
				l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(ResumeMusic))
			} else {
				l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PauseMusic))
			}
		case twodee.Key0:
			l.app.State.Debug = !l.app.State.Debug
			fmt.Printf("Debug state: %v\n", l.app.State.Debug)
			l.app.GameEventHandler.Enqueue(NewShakeEvent(3, 200, 3.0, 4.0, 1.0))
		case twodee.Key1:
			if l.app.State.Debug {
				l.loadLevel("boss1")
			}
		case twodee.Key2:
			if l.app.State.Debug {
				l.loadLevel("boss2")
			}
		case twodee.Key3:
			if l.app.State.Debug {
				l.loadLevel("main")
			}
		case twodee.Key9:
			if l.app.State.Debug {
				if l.level.Boss != nil {
					l.app.GameEventHandler.Enqueue(NewBossDiedEvent(l.level.Boss.Name))
				}
			}
		}
	}
	return true
}

func (l *GameLayer) checkJoy() bool {
	var (
		events = l.app.Context.Events
	)
	if !events.JoystickPresent(twodee.Joystick1) {
		return false
	}
	var (
		axes    []float32 = events.JoystickAxes(twodee.Joystick1)
		buttons []byte    = events.JoystickButtons(twodee.Joystick1)
		x                 = float64(axes[0])
		y                 = float64(-axes[1])
	)
	if math.Abs(x) < 0.2 {
		x = 0.0
	}
	if math.Abs(y) < 0.2 {
		y = 0.0
	}
	l.level.Player.MoveX(float32(x))
	l.level.Player.MoveY(float32(y))
	if len(buttons) > 11 && buttons[11] != 0 { // Very much hardcoded to xbox controller
		l.level.Player.Roll()
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

func (l *GameLayer) getSplashSpriteConfig(name string, camera *twodee.Camera) twodee.SpriteConfig {
	var (
		frame = l.spritesheet.GetFrame(name)
		pt    = camera.WorldBounds.Midpoint()
	)
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			pt.X(), pt.Y(), 0,
			0, 0, 0,
			2.0, 2.0, 1.0,
		},
		Frame: frame.Frame,
	}
}
