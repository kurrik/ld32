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

import twodee "../lib/twodee"

type AudioSystem struct {
	app                         *Application
	bgm                         *twodee.Music
	bossMusic                   *twodee.Music
	bossDeathEffect             *twodee.SoundEffect
	colorChangeEffect           *twodee.SoundEffect
	playerDeathEffect           *twodee.SoundEffect
	rollEffect                  *twodee.SoundEffect
	bgmObserverId               int
	bossMusicObserverId         int
	pauseMusicObserverId        int
	resumeMusicObserverId       int
	bossDeathEffectObserverId   int
	colorChangeEffectObserverId int
	playerDeathEffectObserverId int
	rollEffectObserverId        int
	musicToggle                 int32
}

func (a *AudioSystem) PlayBackgroundMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPlaying() {
			twodee.PauseMusic()
		}
		a.bgm.Play(-1)
	}
}

func (a *AudioSystem) PlayBossMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPlaying() {
			twodee.PauseMusic()
		}
		a.bossMusic.Play(-1)
	}
}

func (a *AudioSystem) PauseMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPlaying() {
			twodee.PauseMusic()
		}
	}
}

func (a *AudioSystem) ResumeMusic(e twodee.GETyper) {
	if a.musicToggle == 1 {
		if twodee.MusicIsPaused() {
			twodee.ResumeMusic()
		}
	}
}

func (a *AudioSystem) PlayBossDeathEffect(e twodee.GETyper) {
	if a.bossDeathEffect.IsPlaying(2) == 0 {
		a.bossDeathEffect.PlayChannel(2, 1)
	}
}

func (a *AudioSystem) PlayColorChangeEffect(e twodee.GETyper) {
	if a.colorChangeEffect.IsPlaying(3) == 0 {
		a.colorChangeEffect.PlayChannel(3, 1)
	}
}

func (a *AudioSystem) PlayPlayerDeathEffect(e twodee.GETyper) {
	if a.playerDeathEffect.IsPlaying(4) == 0 {
		a.playerDeathEffect.PlayChannel(4, 1)
	}
}

func (a *AudioSystem) PlayRollEffect(e twodee.GETyper) {
	if a.rollEffect.IsPlaying(5) == 0 {
		a.rollEffect.PlayChannel(5, 1)
	}
}

func (a *AudioSystem) Delete() {
	a.app.GameEventHandler.RemoveObserver(PlayBackgroundMusic, a.bgmObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayBossMusic, a.bossMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayBossDeathEffect, a.bossDeathEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayColorChangeEffect, a.colorChangeEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayPlayerDeathEffect, a.playerDeathEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(PlayRollEffect, a.rollEffectObserverId)
	a.app.GameEventHandler.RemoveObserver(PauseMusic, a.pauseMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(ResumeMusic, a.resumeMusicObserverId)
	a.bgm.Delete()
	a.bossMusic.Delete()
	a.bossDeathEffect.Delete()
	a.colorChangeEffect.Delete()
	a.playerDeathEffect.Delete()
	a.rollEffect.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	var (
		bgm               *twodee.Music
		bossMusic         *twodee.Music
		bossDeathEffect   *twodee.SoundEffect
		colorChangeEffect *twodee.SoundEffect
		playerDeathEffect *twodee.SoundEffect
		rollEffect        *twodee.SoundEffect
	)

	if bgm, err = twodee.NewMusic("resources/music/Shrine_Theme_Rough.ogg"); err != nil {
		return
	}
	if bossMusic, err = twodee.NewMusic("resources/music/Boss_Theme_Rough.ogg"); err != nil {
		return
	}
	if bossDeathEffect, err = twodee.NewSoundEffect("resources/music/BossDeathEffect.ogg"); err != nil {
		return
	}
	if colorChangeEffect, err = twodee.NewSoundEffect("resources/music/ColorChangeEffect.ogg"); err != nil {
		return
	}
	if playerDeathEffect, err = twodee.NewSoundEffect("resources/music/PlayerDeath.ogg"); err != nil {
		return
	}
	if rollEffect, err = twodee.NewSoundEffect("resources/music/RollEffect.ogg"); err != nil {
		return
	}
	audioSystem = &AudioSystem{
		app:               app,
		bgm:               bgm,
		bossMusic:         bossMusic,
		bossDeathEffect:   bossDeathEffect,
		colorChangeEffect: colorChangeEffect,
		playerDeathEffect: playerDeathEffect,
		rollEffect:        rollEffect,
		musicToggle:       1,
	}
	playerDeathEffect.SetVolume(50)
	audioSystem.bgmObserverId = app.GameEventHandler.AddObserver(PlayBackgroundMusic, audioSystem.PlayBackgroundMusic)
	audioSystem.bgmObserverId = app.GameEventHandler.AddObserver(PlayBossMusic, audioSystem.PlayBossMusic)
	audioSystem.pauseMusicObserverId = app.GameEventHandler.AddObserver(PauseMusic, audioSystem.PauseMusic)
	audioSystem.resumeMusicObserverId = app.GameEventHandler.AddObserver(ResumeMusic, audioSystem.ResumeMusic)
	audioSystem.bossDeathEffectObserverId = app.GameEventHandler.AddObserver(PlayBossDeathEffect, audioSystem.PlayBossDeathEffect)
	audioSystem.colorChangeEffectObserverId = app.GameEventHandler.AddObserver(PlayColorChangeEffect, audioSystem.PlayColorChangeEffect)
	audioSystem.playerDeathEffectObserverId = app.GameEventHandler.AddObserver(PlayPlayerDeathEffect, audioSystem.PlayPlayerDeathEffect)
	audioSystem.rollEffectObserverId = app.GameEventHandler.AddObserver(PlayRollEffect, audioSystem.PlayRollEffect)
	return
}
