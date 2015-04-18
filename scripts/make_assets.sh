#!/usr/bin/env bash

mkdir -p tmp

build_aesprite() {
  aseprite \
    --batch assets/${1}.ase \
    --save-as tmp/${1}_01.png
}


TexturePacker \
  --format json-array \
  --trim-sprite-names \
  --size-constraints POT \
  --disable-rotation \
  --data src/resources/spritesheet.json \
  --sheet src/resources/spritesheet.png \
  tmp

rm -rf tmp
