# Undeploy
cd oakestra-controller-manager
make undeploy
cd ..
kubectl delete -f oakestra-agent/Deployment/oakestra-agent.yaml
kubectl delete -f oakestra-network/Deployment/oakestra-nodenetmanager/node-netmanager.yaml

# Redeploy
kubectl create namespace oakestra-controller-manager
kubectl apply -f oakestra-network/Deployment/oakestra-nodenetmanager/node-netmanager.yaml -n oakestra-system
kubectl apply -f oakestra-agent/Deployment/oakestra-agent.yaml
cd oakestra-controller-manager
make deploy
cd ..