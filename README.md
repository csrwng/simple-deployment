## Simple Kubernetes Deployment Image
Provides a simple deployment script that takes the following env variables:
- KUBERNETES_MASTER - the URL to the Kubernetes master (defaults to http://localhost:8080)
- NAME - Name for the deployment - tags replication controllers with the same name
- TARGET_IMAGE - the image to deploy
- TARGET_REPLICAS - the number of replicas to create
- PORT_MAPPING - the port mapping for each replica

### Building the binary
Get and build the Kubernetes source
Edit build-binary.sh to point to the Kubernetes source
run ./build-binary.sh

### Building the image
Edit build-image.sh and set your DOCKER_USER name
run ./build-image.sh

### Pushing the image to DockerHub
```
sudo docker push <yourname>/simple-deploy
```

### Running a deployment
```
cd <KUBERNETES_SOURCE>
cluster/kubecfg.sh -c <SIMPLE_DEPLOY_SOURCE>/json/blue8081.json create deployments
```
This will create a set of simple nginx pods mapped to 8081 and respond with a data.json with the color blue
If you want to deploy a new set of pods, that serve the color red do
```
cd <KUBERNETES_SOURCE>
cluster/kubecfg.sh -c <SIMPLE_DEPLOY_SOURCE/json/red8082.json create deployments
```
