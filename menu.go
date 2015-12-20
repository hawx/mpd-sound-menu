package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/guelfey/go.dbus"
	"github.com/guelfey/go.dbus/introspect"
	"github.com/guelfey/go.dbus/prop"
	"hawx.me/code/mpd-sound-menu/mpd"
	"hawx.me/code/mpd-sound-menu/mpris"
)

const (
	MprisPath       = "/org/mpris/MediaPlayer2"
	RootInterface   = "org.mpris.MediaPlayer2"
	PlayerInterface = RootInterface + ".Player"
	Name            = "org.mpris.MediaPlayer2."
)

var (
	network   = flag.String("net", "tcp", "")
	address   = flag.String("addr", ":6600", "")
	localName = flag.String("name", "mpd-sound-menu", "")
)

func Start(client *mpd.Mpd) (*dbus.Conn, error) {
	root := mpris.Root{}
	player := mpris.Player{client}

	conn, err := dbus.SessionBus()
	if err != nil {
		return conn, fmt.Errorf("Failed to connect to session bus: %s", err)
	}

	reply, err := conn.RequestName(Name+*localName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return conn, fmt.Errorf("Failed to get name: %s", err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return conn, fmt.Errorf("Name already taken")
	}

	conn.Export(root, MprisPath, RootInterface)
	conn.Export(player, MprisPath, PlayerInterface)

	propsSpec := map[string]map[string]*prop.Prop{
		RootInterface:   mpris.RootProps(*localName),
		PlayerInterface: mpris.PlayerProps(),
	}

	props := prop.New(conn, MprisPath, propsSpec)

	n := &introspect.Node{
		Name: MprisPath,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name:       RootInterface,
				Methods:    introspect.Methods(root),
				Properties: props.Introspection(RootInterface),
			},
			{
				Name:       PlayerInterface,
				Methods:    introspect.Methods(player),
				Properties: props.Introspection(PlayerInterface),
				Signals: []introspect.Signal{{
					Name: "Seeked",
					Args: []introspect.Arg{{
						Name: "Position",
						Type: "x",
					}},
					Annotations: []introspect.Annotation{},
				}},
			},
		},
	}

	conn.Export(introspect.NewIntrospectable(n), MprisPath, "org.freedesktop.DBus.Introspectable")

	client.Update(&propertyAdapter{props})

	return conn, nil
}

type propertyAdapter struct {
	props *prop.Properties
}

func (u *propertyAdapter) UpdatePlaybackStatus(state mpd.PlayState) {
	u.props.Set(PlayerInterface, "PlaybackStatus", dbus.MakeVariant(string(state)))
}

func (u *propertyAdapter) UpdateCurrentSong(song mpd.Song) {
	u.props.Set(PlayerInterface, "Metadata", dbus.MakeVariant(
		map[string]dbus.Variant{
			"mrpis:trackid": dbus.MakeVariant(song.Id),
			"mpris:length":  dbus.MakeVariant(song.Length),
			"xesam:title":   dbus.MakeVariant(song.Title),
			"xesam:artist":  dbus.MakeVariant(song.Artist),
			"xesam:album":   dbus.MakeVariant(song.Album),
		},
	))
}

func main() {
	flag.Parse()

	mpdclient, err := mpd.Dial(*network, *address)
	if err != nil {
		log.Fatal(err)
	}
	defer mpdclient.Close()

	conn, err := Start(mpdclient)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	log.Printf("caught %s: shutting down", s)
}
