# mpd-sound-menu

Tool to display/control mpd from Ubuntu's sound menubar item.

Install with:

``` bash
$ go get hawx.me/code/mpd-sound-menu
```

Then add a `.desktop` entry for it:

``` bash
$ cat > ~/.local/share/applications/mpd-sound-menu.desktop
[Desktop Entry]
Name=mpd-sound-menu
Comment=Sound menu plugin to display mpd
Exec=/path/to/mpd-sound-menu
Icon=mpf-sound-menu
Type=Application
```

Replace the value of "Exec" with the correct path to _mpd-sound-menu_.

Finally start _mpd-sound-menu_ (I use upstart), and you should see it in the
sound menubar item.
