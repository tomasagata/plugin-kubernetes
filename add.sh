kubectl create namespace oakestra-controller-manager
kubectl apply -f oakestra-network/Deployment/oakestra-nodenetmanager/node-netmanager.yaml -n oakestra-system
kubectl apply -f oakestra-agent/Deployment/oakestra-agent.yaml
cd oakestra-controller-manager
make install
make deploy
cd ..