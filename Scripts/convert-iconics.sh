#!/bin/bash

SOURCE="../Iconics"
TARGET="../../Character Sheets Website/public/images/iconics";
TXT="$TARGET/iconics.txt"
echo -n "" > "$TXT"
cd "$SOURCE"
find . -type f | while read F
do
  F=$(echo "$F" | sed 's/^..//')
  FB=$(dirname "$F")
  FN=$(basename "$(basename "$F" .jpg)" .png)
  if [ "$FB" != "." ]
  then
    NAME="$FB/$FN"
    LOC=$(echo "$FB/$FN" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9\/]\+/-/g' | sed 's/^-\+//' | sed 's/-\+$//')

    if echo "$F" | grep '.txt$'; then
      echo " * $NAME"
      DF="$TARGET/large/$(echo "$LOC" | sed 's/-txt$//').txt"
      echo 
      cp "$F" "$DF"
    else
      LARGE="$TARGET/large/$LOC.png"
      SMALL="$TARGET/small/$LOC.png"

      echo " * $NAME"

      echo "   * Large file: $LARGE"
      mkdir -p "$(dirname "$LARGE")"
      convert -set colorspace sRGB -bordercolor white -border 10x10 +profile '*' "$F" png:- | ../Scripts/magicwand 0,0 -t 4 -r outside -c trans png:- png:- | convert -trim png:- "$LARGE"

      echo "   * Small file: $SMALL"
      mkdir -p "$(dirname "$SMALL")"
      convert -resize 300x400 "$LARGE" "$SMALL"

      echo "$LOC=$NAME" >> "$TXT"
    fi
  fi
done
