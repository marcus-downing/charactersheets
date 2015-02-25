#!/bin/bash
source /etc/init.d/functions.sh
eval $(eval_ecolors)

SOURCE="../Iconics"
TARGET="../../Character Sheets Website/public/images/iconics";
TXT="$TARGET/iconics.txt"
echo -n "" > "$TXT"
cd ../Iconics
find . -type f | while read F
do
  F=$(echo "$F" | sed 's/^..//')
  FB=$(dirname "$F")
  FN=$(basename "$(basename "$F" .jpg)" .png)
  if [ "$FB" != "." ]
  then
    NAME="$FB/$FN"
    LOC=$(echo "$FB/$FN" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9\/]\+/-/g')
    LARGE="$TARGET/large/$LOC.png"
    SMALL="$TARGET/small/$LOC.png"

    einfo "$NAME"
    eindent

    ebegin "Large file: $LARGE"
    mkdir -p "$(dirname "$LARGE")"
    convert -bordercolor white -border 10x10 "$F" png:- | ../Scripts/magicwand 0,0 -t 4 -r outside -c trans png:- png:- | convert -trim png:- "$LARGE"
    eend $?

    ebegin "Small file: $SMALL"
    mkdir -p "$(dirname "$SMALL")"
    convert -resize 300x400 "$LARGE" "$SMALL"
    eend $?

    echo "$LOC=$NAME" >> "$TXT"
    eoutdent
  fi
done
