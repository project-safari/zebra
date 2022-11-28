#!/bin/sh

if which node > /dev/null
    then
        echo "node is installed, skipping..."
    else
        npm install -g npm
    fi

path="$(pwd)"
pathBackend=$path"/zebra"
pathUI=$path"/zebra-ui"

(cd $pathBackend; make simulator) &
(cd $pathUI; npm i --legacy-peer-deps; npm i mongoose; npm start) &