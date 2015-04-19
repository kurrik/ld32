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
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"time"

	"../lib/twodee"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/tmxgo"
)

type Level struct {
	Player          *Player
	Boss            *Boss
	Props           PropList
	Background      *twodee.Batch
	Sheet           *twodee.Spritesheet
	Collisions      *twodee.Grid
	Portals         []Portal
	Plates          PropList
	Width           float32
	Height          float32
	Color           mgl32.Vec3
	events          *twodee.GameEventHandler
	colorObserverId int
}

type Portal struct {
	Rect  twodee.Rectangle
	Level string
}

func NewLevel(mapPath string, sheet *twodee.Spritesheet, events *twodee.GameEventHandler) (level *Level, err error) {
	level = &Level{
		Player: NewPlayer(events, sheet),
		Props:  NewPropList(),
		Plates: NewPropList(),
		Sheet:  sheet,
		events: events,
	}
	level.Props = append(level.Props, level.Player)
	if err = level.loadMap(mapPath); err != nil {
		return
	}
	if level.Boss != nil {
		level.Props = append(level.Props, level.Boss)
	}
	level.colorObserverId = events.AddObserver(ChangeColor, level.changeColor)
	return
}

func (l *Level) changeColor(e twodee.GETyper) {
	if event, ok := e.(*ColorEvent); ok {
		var (
			sentEvent = false
		)
		if event.Add {
			l.Color = l.Color.Add(event.Color)
		} else {
			l.Color = l.Color.Sub(event.Color)
		}
		if l.Boss != nil {
			if l.Color.Sub(l.Boss.Color).Len() < 0.1 {
				l.Boss.NextColor()
				l.events.Enqueue(NewShakeEvent(2, 1000, 1.0, 10.0, 1.0))
				sentEvent = true
			}
		}
		if !sentEvent {
			l.events.Enqueue(NewShakeEvent(1, 200, 0.4, 2.0, 1.0))
		}
	}
}

func (l *Level) Update(elapsed time.Duration) {
	// TODO: Probably this should update a slice of Mobs or other
	// updateable things in the level.
	if l.Boss != nil {
		l.Boss.Update(elapsed)
		l.Boss.ExamineWorld(l)
	}
	l.Player.UpdateLevel(elapsed, l)
	l.Plates.Update(elapsed)
	l.Plates.CheckCollision(l.Player)
}

func (l *Level) Delete() {
	if l.colorObserverId != 0 {
		l.events.RemoveObserver(ChangeColor, l.colorObserverId)
	}
}

func (l *Level) loadMap(path string) (err error) {
	var (
		data        []byte
		m           *tmxgo.Map
		tiles       []*tmxgo.Tile
		textiles    []twodee.TexturedTile
		texturepath string
		colorbytes  []byte
	)
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if m, err = tmxgo.ParseMapString(string(data)); err != nil {
		return
	}
	l.Collisions = twodee.NewGrid(m.Width, m.Height)
	l.Width = float32(m.Width*m.TileWidth) / PxPerUnit
	l.Height = float32(m.Height*m.TileHeight) / PxPerUnit
	if tiles, err = m.TilesFromLayerName("collision"); err == nil {
		// Able to find collision tiles
		for i, t := range tiles {
			if t != nil {
				l.Collisions.SetIndex(int32(i), true)
			}
		}
	}
	if tiles, err = m.TilesFromLayerName("ground"); err != nil {
		return
	}
	if texturepath, err = tmxgo.GetTexturePath(tiles); err != nil {
		return
	}
	textiles = make([]twodee.TexturedTile, len(tiles))
	j := 0
	for i, t := range tiles {
		if t != nil {
			textiles[i] = t
			j++
		}
	}
	textiles = textiles[:j]
	var (
		tilem = twodee.TileMetadata{
			Path:      filepath.Join(filepath.Dir(path), texturepath),
			PxPerUnit: PxPerUnit,
		}
	)
	for _, objgroup := range m.ObjectGroups {
		for _, obj := range objgroup.Objects {
			x, y := l.getObjectMiddle(m, obj)
			switch obj.Name {
			case "portal":
				l.Portals = append(l.Portals, Portal{
					Rect:  l.getObjectBounds(m, obj),
					Level: obj.Type,
				})
			case "plate":
				if colorbytes, err = hex.DecodeString(obj.Type); err != nil {
					return
				}
				color := mgl32.Vec3{
					float32(colorbytes[0]) / 255.0,
					float32(colorbytes[1]) / 255.0,
					float32(colorbytes[2]) / 255.0,
				}
				l.Plates = append(l.Plates, NewPlate(x, y, color, l.Sheet, l.events))
			case "sprite":
				l.Props = append(l.Props, NewStaticProp(
					x, y,
					l.Sheet,
					obj.Type,
				))
			case "start":
				l.Player.MoveTo(twodee.Pt(x, y))
			case "boss":
				l.Boss = BossMap[obj.Type](x, y, l.events)
				l.Boss.MoveTo(twodee.Pt(x, y))
			}
		}
	}
	l.Background, err = twodee.LoadBatch(textiles, tilem)
	return
}

func (l *Level) getObjectMiddle(m *tmxgo.Map, obj tmxgo.Object) (x float32, y float32) {
	x = float32(obj.X+(obj.Width/2.0)) / PxPerUnit
	y = float32(m.Height*m.TileHeight-obj.Y-(obj.Height/2.0)) / PxPerUnit // Height is reversed
	return
}

func (l *Level) getObjectBounds(m *tmxgo.Map, obj tmxgo.Object) twodee.Rectangle {
	var (
		x = float32(obj.X) / PxPerUnit
		y = float32(m.Height*m.TileHeight-obj.Y) / PxPerUnit // Height is reversed
		w = float32(obj.Width) / PxPerUnit
		h = float32(obj.Height) / PxPerUnit
	)
	return twodee.Rect(x, y-h, x+w, y)
}

func (l *Level) PortalCollides() (bool, string) {
	for _, portal := range l.Portals {
		if portal.Rect.Overlaps(l.Player.Bounds()) {
			return true, portal.Level
		}
	}
	return false, ""
}
