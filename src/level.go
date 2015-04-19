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
	"io/ioutil"
	"path/filepath"
	"time"

	twodee "../lib/twodee"
	"github.com/kurrik/tmxgo"
)

type Level struct {
	Player     *Player
	Boss       *Boss
	Props      PropList
	Background *twodee.Batch
	Sheet      *twodee.Spritesheet
	Collisions *twodee.Grid
	Portals    []Portal
	Width      float32
	Height     float32
}

type Portal struct {
	Rect  twodee.Rectangle
	Level string
}

func NewLevel(mapPath string, sheet *twodee.Spritesheet, events *twodee.GameEventHandler) (level *Level, err error) {
	level = &Level{
		Player: NewPlayer(events),
		Props:  NewPropList(),
		Sheet:  sheet,
	}
	level.Props = append(level.Props, level.Player)
	if err = level.loadMap(mapPath); err != nil {
		return
	}
	return
}

func (l *Level) Update(elapsed time.Duration) {
	l.Player.Update(elapsed, l)
}

func (l *Level) loadMap(path string) (err error) {
	var (
		data        []byte
		m           *tmxgo.Map
		tiles       []*tmxgo.Tile
		textiles    []twodee.TexturedTile
		texturepath string
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
	for i, t := range tiles {
		textiles[i] = t
	}
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
			case "sprite":
				l.Props = append(l.Props, NewStaticProp(
					x, y,
					l.Sheet,
					obj.Type,
				))
			case "start":
				l.Player.MoveTo(twodee.Pt(x, y))
			case "boss":
				l.Boss = BossMap[obj.Type]()
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
