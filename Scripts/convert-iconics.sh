source /etc/init.d/functions.sh
eval $(eval_ecolors)

TARGET="../Composer/public/images/iconics";
cd ../Iconics
find . -type f | while read F
do
  F=$(echo "$F" | sed 's/^..//')
  FB=$(dirname "$F")
  FN=$(basename "$(basename "$F" .jpg)" .png)
  if [ "$FB" != "." ]
  then
    ebegin "$FB/${GOOD}$FN"
    mkdir -p "../Composer/public/images/iconics/$FB/Large"
    mkdir -p "../Composer/public/images/iconics/$FB/Small"
    convert -bordercolor white -border 10x10 "$F" png:- | ../Scripts/magicwand 0,0 -t 2 -r outside -c trans png:- png:- | convert -trim png:- "$TARGET/$FB/Large/$FN.png"
    eend $?
    convert -resize 500x500 "$TARGET/$FB/Large/$FN.png" "$TARGET/$FB/Small/$FN.png"
  fi
done
