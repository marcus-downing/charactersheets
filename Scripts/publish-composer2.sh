#!/bin/bash

rsync -r --exclude=About --exclude=.svn --exclude=*.ai --exclude=*.zip Pathfinder/* Composer2/public/pdf/pathfinder/
rsync -r --exclude=About --exclude=.svn --exclude=*.ai --exclude=*.zip 3.5/* Composer2/public/pdf/dnd35/

#rsync --exclude=About --exclude=*.ai Pathfinder/* Composer2/public/pdf/pathfinder/
#rsync --exclude=About --exclude=*.ai 3.5/* Composer2/public/pdf/dnd35/
