#!/bin/bash
set -eu

DIR="$1"

if [ ! -d $DIR ]; then
	exit 1
fi
cd $DIR

echo pid,bid,frame,side,width,height

find . -type f -path "*/images/*.jpg" | while read jpg; do
  if [[ $jpg =~ ^\./([0-9]{9})/images/([0-9]{9}_([0-9]{5})(_([12]))?)\.jpg$ ]]; then
    bid=${BASH_REMATCH[1]}
    pid=${BASH_REMATCH[2]}
    frame=${BASH_REMATCH[3]}
    side=${BASH_REMATCH[5]}
    wh=($(identify -format "%w %h" $jpg))
    echo $pid,$bid,$frame,$side,${wh[0]},${wh[1]}
  else
    echo ng: $jpg >&2
fi
done
exit 0
