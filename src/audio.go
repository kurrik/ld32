package main

import twodee "../lib/twodee"

type AudioSystem struct {
	app                   *Application
	bgm                   *twodee.Music
	bgmObserverId         int
	pauseMusicObserverId  int
	resumeMusicObserverId int
}

func (a *AudioSystem) PlayMusic(e twodee.GETyper) {
	a.bgm.Play(-1)
}

func (a *AudioSystem) PauseMusic(e twodee.GETyper) {
	if twodee.MusicIsPlaying() {
		twodee.PauseMusic()
	}
}

func (a *AudioSystem) ResumeMusic(e twodee.GETyper) {
	if twodee.MusicIsPaused() {
		twodee.ResumeMusic()
	}
}

func (a *AudioSystem) Delete() {
	a.app.GameEventHandler.RemoveObserver(PlayMusic, a.bgmObserverId)
	a.app.GameEventHandler.RemoveObserver(PauseMusic, a.pauseMusicObserverId)
	a.app.GameEventHandler.RemoveObserver(ResumeMusic, a.resumeMusicObserverId)
	a.bgm.Delete()
}

func NewAudioSystem(app *Application) (audioSystem *AudioSystem, err error) {
	var bgm *twodee.Music

	if bgm, err = twodee.NewMusic("resources/music/bgm1.ogg"); err != nil {
		return
	}
	audioSystem = &AudioSystem{
		app: app,
		bgm: bgm,
	}
	audioSystem.bgmObserverId = app.GameEventHandler.AddObserver(PlayMusic, audioSystem.PlayMusic)
	audioSystem.pauseMusicObserverId = app.GameEventHandler.AddObserver(PauseMusic, audioSystem.PauseMusic)
	audioSystem.resumeMusicObserverId = app.GameEventHandler.AddObserver(ResumeMusic, audioSystem.ResumeMusic)
	return
}
