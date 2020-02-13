package mpris

import (
	dbus "github.com/guelfey/go.dbus"
	"github.com/guelfey/go.dbus/prop"
)

// http://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html

type Root struct{}

func (m Root) Raise() *dbus.Error {
	return nil
}

func (m Root) Quit() *dbus.Error {
	return nil
}

func RootProps(name string) map[string]*prop.Prop {
	return map[string]*prop.Prop{
		"CanQuit": &prop.Prop{
			Value:    false,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"CanRaise": &prop.Prop{
			Value:    false,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"HasTrackList": &prop.Prop{
			Value:    false,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"Identity": &prop.Prop{
			Value:    name,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"DesktopEntry": &prop.Prop{
			Value:    name,
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"SupportedUriSchemes": &prop.Prop{
			Value:    []string{"file"}, // get from mpd?
			Writable: false,
			Emit:     prop.EmitFalse,
		},
		"SupportedMimeTypes": &prop.Prop{
			Value:    []string{"audio/mpeg"}, // get from mpd?
			Writable: false,
			Emit:     prop.EmitFalse,
		},
	}
}
