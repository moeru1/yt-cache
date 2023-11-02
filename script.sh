#!/bin/sh

[ -d "/tmp/dv" ] || mkdir /tmp/dv
maxsize=`df /tmp | awk '/^\//{print ($4/2)*1000}'`

[ "${1%%:*}" = "https" -o "${1%%:*}" = "http" ] && link="$1" || link=`xclip -sel c -o`

link="${link%%&*}"

yt-dlp -ic --output "/tmp/dv/dvout-%(playlist-index)s-%(id)s.%(ext)s" --print after_move:"/tmp/dv/dvout-%(playlist-index)s-%(id)s.%(ext)s" "$link" -N8 -f "bv*[filesize<${maxsize}]+ba / b[filesize<${maxsize}] / w" --add-metadata


