#!/bin/bash
set -eu

dir="$1"
if [ ! -d "$dir" ]; then
  if [ -d "../../data/dataset/" ] ; then
    dir="../../data/dataset/"
  else
    echo "USAGE: $0 <dataset root dir>"
    exit 1
  fi
fi

echo "register page"
TMP="/tmp/tmp.csv"
TMPSORTED="/tmp/tmpsorted.csv"
bash ./pagesize2list.sh $dir > $TMP
(head -n+1 $TMP && tail -n+2 $TMP | sort) > $TMPSORTED
curl -u admin:jikei-admin -X POST http://localhost:8080/jikei/api/admin/pagelistupload -F "file=@$TMPSORTED"
rm $TMP $TMPSORTED

echo "register jikei"
find $dir -type f -path */[12]*_coordinate.csv -exec curl -u admin:jikei-admin -X POST http://localhost:8080/jikei/api/admin/jikeilistupload -F "file=@{}" \;
