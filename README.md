# xplay

A simple media server that searches media files from a directory, generates xspf playlist as index and serves audio/video over http.

Currently supported file formats: mp3, flac, ogg, mp4, mkv

[![This program produces valid XSPF playlist files.](img/valid-xspf.png)](https://validator.xspf.org/referrer/)

This program:

- Does not follow symbolic links found in directories
- Excludes files starting with a period (hidden files in linux)

## Usage

Start http media server with `/play.xspf` as index. Use `-r` option to traverse the directory recursively.

```bash
./xplay -b $bind_ipaddr -p $bind_port -d ./music -r
```

Use `-w` to generate and save xspf to file and exit. `-b` and `-p` options will be ignored:

```bash
./xplay -d . -r -w > playlist.xspf
```

## Client

Media players with http and xspf support (like VLC) can be used as clients.

```bash
vlc http://$ip:$port/play.xspf
```
