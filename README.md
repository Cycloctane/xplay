# xplay

A simple media server that searches multimedia files from a directory, parses their metadata, generates xspf playlist as index and serves audio/video over http.

Supported file formats: mp3, flac, ogg, mp4, mkv. Metadata parsing currently does not support mkv.

[![This program produces valid XSPF playlist files.](assets/valid-xspf.png)](https://validator.xspf.org/referrer/)

This program:

- Traverses directories recursively by default. Can be disabled by `--no-recursive` option
- Does not follow symbolic links found in directories
- Excludes files starting with a period (hidden files in linux)

## Usage

Start http media server with `/play.xspf` as index:

```bash
./xplay -b $bind_ipaddr -p $bind_port -d ./music
```

Use `-w` to generate and save xspf to file and exit. `-b` and `-p` options will be ignored:

```bash
./xplay -d . -w > playlist.xspf
```

Metadata parsing can become slow when handling a large number of multimedia files. Use `--no-tag` option to disable metadata/tag parsing if you do not need metadata in xspf playlists.

## Client

Media players with http and xspf support (like VLC) can be used as clients.

```bash
vlc http://$ip:$port/play.xspf
```
