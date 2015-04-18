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
	"github.com/kurrik/tmxgo"
	"io/ioutil"
	"path/filepath"
	"time"
)

type Level struct {
	Player     *Player
	Background *twodee.Batch
}

func NewLevel(mapPath string) (level *Level, err error) {
	level = &Level{
		Player: NewPlayer(),
	}
	if err = level.loadMap(mapPath); err != nil {
		return
	}
	return
}

func (l *Level) Update(elapsed time.Duration) {
	l.Player.Update(elapsed)
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
			if obj.Name == "start" {
				l.Player.MoveTo(twodee.Pt(
					float32(obj.X)/PxPerUnit,
					float32(m.Height*m.TileHeight-obj.Y)/PxPerUnit, // Height is reversed
				))
			}
		}
	}
	l.Background, err = twodee.LoadBatch(textiles, tilem)
	return
}
