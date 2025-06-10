cd oakestra-controller-manager
make undeploy
cd ..
kubectl delete -f oakestra-agent/Deployment/oakestra-agent.yaml
kubectl delete -f oakestra-network/Deployment/oakestra-nodenetmanager/node-netmanager.yaml