package mpris

import (
	"github.com/guelfey/go.dbus"
	"github.com/guelfey/go.dbus/prop"
)

// http://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html

type Controller interface {
	Next()
	Previous()
	Play()
	Pause()
	PlayPause()
	Stop()
}

type Player struct {
	Control Controller
}

func (p Player) Next() *dbus.Error {
	p.Control.Next()
	return nil
}

func (p Player) Previous() *dbus.Error {
	p.Control.Previous()
	return nil
}

func (p Player) Pause() *dbus.Error {
	p.Control.Pause()
	return nil
}

func (p Player) PlayPause() *dbus.Error {
	p.Control.PlayPause()
	return nil
}

func (p Player) Stop() *dbus.Error {
	p.Control.Stop()
	return nil
}

func (p Player) Play() *dbus.Error {
	p.Control.Play()
	return nil
}

func (p Player) Seek(x int64) *dbus.Error {
	return nil
}

func (p Player) SetPosition(o dbus.ObjectPath, x int64) *dbus.Error {
	return nil
}

func (p Player) OpenUri(s string) *dbus.Error {
	return nil
}

func PlayerProps() map[string]*prop.Prop {
	return map[string]*prop.Prop{
		"PlaybackStatus": &prop.Prop{
			Value:    "Stopped",
			Writable: true,
			Emit:     prop.EmitTrue,
		},
		"LoopStatus": &prop.Prop{
			Value:    "None", // one of: "None", "Track", "Playlist"
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"Rate": &prop.Prop{
			Value:    1.0,
			Writable: true,
			Emit:     prop.EmitFalse,
		},
		"Shuffle": &prop.Prop{
			Value:    false,
			Writable: true,
			Emit:     prop.EmitFalse,
		},
		"Metadata": &prop.Prop{
			Value:    map[string]dbus.Variant{},
			Writable: true,
			Emit:     prop.EmitTrue,
		},
		"Volume": &prop.Prop{
			Value:    0.0,
			Writable: true,
			Emit:     prop.EmitFalse,
			Callback: func(change *prop.Change) *dbus.Error {
				return nil // do nothing
			},
		},
		"Position": &prop.Prop{
			Value:    0.0,
			Writable: false,
			Emit:     prop.EmitTrue,
		},
		"MinimumRate": &prop.Prop{
			Value:    1.0,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"MaximumRate": &prop.Prop{
			Value:    1.0,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"CanGoNext": &prop.Prop{
			Value:    false,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"CanGoPrevious": &prop.Prop{
			Value:    false,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"CanPlay": &prop.Prop{
			Value:    true,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"CanPause": &prop.Prop{
			Value:    true,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"CanSeek": &prop.Prop{
			Value:    false,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"CanControl": &prop.Prop{
			Value:    true,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
	}
}
