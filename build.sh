cd oakestra-network/oakestra-nodeNetManager
docker build -t tomasagata/nodenetmanager:latest .
docker push tomasagata/nodenetmanager:latest
cd ../oakestra-cni
docker build -t tomasagata/oakestra-cni:latest .
docker push tomasagata/oakestra-cni:latest
cd ../../..
