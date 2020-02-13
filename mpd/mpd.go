package mpd

import (
	"log"
	"strconv"
	"time"

	"github.com/fhs/gompd/mpd"
)

type PlayState string

const (
	Playing PlayState = "Playing"
	Paused            = "Paused"
	Stopped           = "Stopped"
)

type Song struct {
	Id     string
	Length int64
	Title  string
	Artist string
	Album  string
}

type Notifier interface {
	UpdatePlaybackStatus(PlayState)
	UpdateCurrentSong(Song)
}

type Mpd struct {
	client   *mpd.Client
	watcher  *mpd.Watcher
	notifier Notifier
}

func Dial(network, address string) (*Mpd, error) {
	client, err := mpd.Dial(network, address)
	if err != nil {
		client.Close()
		return nil, err
	}

	watcher, err := mpd.NewWatcher(network, address, "")
	if err != nil {
		client.Close()
		watcher.Close()
		return nil, err
	}

	go func() {
		for _ = range time.Tick(30 * time.Second) {
			client.Ping()
		}
	}()

	return &Mpd{client: client, watcher: watcher}, nil
}

func (w *Mpd) Close() {
	w.client.Close()
	w.watcher.Close()
}

func (w *Mpd) notifyPlaybackStatus() {
	w.notifier.UpdatePlaybackStatus(w.PlayState())
}

func (w *Mpd) notifyCurrentSong() {
	attrs, err := w.client.CurrentSong()
	if err != nil {
		log.Println(err)
		return
	}

	songtime, err := strconv.Atoi(attrs["Time"])
	if err != nil {
		w.notifier.UpdateCurrentSong(Song{})
		return
	}

	w.notifier.UpdateCurrentSong(Song{
		Id:     attrs["Id"],
		Length: int64(songtime) * 1e6,
		Title:  attrs["Title"],
		Artist: attrs["Artist"],
		Album:  attrs["Album"],
	})
}

func (w *Mpd) Update(notifier Notifier) {
	w.notifier = notifier

	go func() {
		w.notifyPlaybackStatus()
		w.notifyCurrentSong()

		for subsystem := range w.watcher.Event {
			if subsystem == "player" {
				w.notifyPlaybackStatus()
				w.notifyCurrentSong()
			}
		}
	}()
}

func (w *Mpd) Play() {
	w.client.Play(0)
	w.notifyPlaybackStatus()
}

func (w *Mpd) Pause() {
	w.client.Pause(true)
	w.notifyPlaybackStatus()
}

func (w *Mpd) PlayPause() {
	if w.PlayState() == Playing {
		w.client.Pause(true)
	} else {
		w.client.Pause(false)
	}
	w.notifyPlaybackStatus()
}

func (w *Mpd) Stop() {
	w.client.Stop()
	w.notifyPlaybackStatus()
}

func (w *Mpd) Next() {
	w.client.Next()
}

func (w *Mpd) Previous() {
	w.client.Previous()
}

func (w *Mpd) PlayState() PlayState {
	attrs, _ := w.client.Status()

	switch attrs["state"] {
	case "play":
		return Playing
	case "pause":
		return Paused
	}
	return Stopped
}
