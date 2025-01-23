# xplay-compatible

This branch is an alternative version of xplay. Compared with the main branch, this branch is compatible with a wider range of media players (like potplayer), but does not have basic authentication feature.

[![This program produces valid XSPF playlist files.](assets/valid-xspf.png)](https://validator.xspf.org/referrer/)

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

To secure the media server, activate https with `--ssl-cert` `--ssl-key` and restrict host header by setting expected server hostnames in `--allowed-hosts`:

```bash
./xplay -b 0.0.0.0 -p 8443 -d ./music\
    --ssl-cert example.com.crt --ssl-key example.com.pem\
    --allowed-hosts example.com,www.example.com,127.0.0.1
```

## Client

Media players with http and xspf support (like VLC) can be used as clients. Tested with potplayer, Clementine, audacious...

```bash
vlc http://$ip:$port/play.xspf
```

With https and http basic auth:

```bash
vlc https://$username:$password@$ip:$port/play.xspf
```
