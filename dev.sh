#!/bin/sh

if which node > /dev/null
    then
        echo "node is installed, skipping..."
    else
        npm install -g npm
    fi

if which serve > /dev/null
    then
        echo "serve is installed, skipping..."
    else
     npm install -g serve 
   fi

path="$(pwd)"
pathBackend=$path"/zebra"
pathUI=$path"/zebra-ui"

(cd $pathBackend; go build; ./zebra; ./zebra-server; ./herd; ./simulator/gen_certs.sh "./simulator" "zebra.inisieme.local" "1.1.1.1"; cd ./simulator && go test; cd $pathBackend; ./zebra; ./zebra-server; ./herd; 	rm -rf ./simulator/simulator-store && ./herd --store ./simulator/simulator-store; rm -f ./simulator/zebra-simulator.json; rm -f ./simulator/admin.yaml; 	./zebra -c ./simulator/admin.yaml config init https://127.0.0.1:6666; 	./zebra -c ./simulator/admin.yaml config email admin@zebra.project-safari.io; 	./zebra -c ./simulator/admin.yaml config user admin; 	./zebra -c ./simulator/admin.yaml config ca-cert ./simulator/zebra-ca.crt; 	./zebra-server -c ./simulator/zebra-simulator.json init --auth-key "AvadaKedavra" --user="./simulator/admin.yaml" --password "Riddikulus" --cert "./simulator/zebra-server.crt" --key "./simulator/zebra-server.key" -a "tcp://127.0.0.1:6666" --store="./simulator/simulator-store"; 	./zebra-server --config ./simulator/zebra-simulator.json ) &
(cd $pathUI; npm run build;  serve -s build) &