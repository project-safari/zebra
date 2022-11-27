#!/bin/sh\

(cd /Users/eachim/desktop/zebra-all/zebras/zebra; make simulator) &
(cd /Users/eachim/desktop/zebra-all/zebras/zebra-ui; npm start) &