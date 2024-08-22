# xplay

A simple media server that reads a directory, generates xspf playlist as index and serves audio/video over http.

Currently supported file formats: mp3, flac, ogg, mp4, mkv

[![This program produces valid XSPF playlist files.](img/valid-xspf.png)](https://validator.xspf.org/referrer/)

This program:

- Follows symlinks
- Excludes files starting with a period (hidden files in linux)
- Ignores files in subdirectories

## Usage

Start http media server with `/play.xspf` as index:

```bash
./xplay -b $bind_ipaddr -p $bind_port -d ./music
```

Use `-w` to generate and save xspf to file and exit. `-b` and `-p` options will be ignored:

```bash
./xplay -d . -w > playlist.xspf
```

## Client

Media players with http and xspf support (like VLC) can be used as clients.

```bash
vlc http://$ip:$port/play.xspf
```
