# Forgerock

## Build the container image
First we need to login to docker hub so we can push our container image up to docker hub. From the console run **docker login** and hit enter. Then enter your credentials and you should be logged in. In this example my docker hub repo is called elminestrone(weird name yes, its a long story). Next step is to build the image, 
from the root of the repository run:
```
docker build --tag elminestrone/go-stock-api:1.0 .
```

This will build the docker image and tag it in the format of dockerHubRepoName/imageName:version. Once we have it tagged appropriately we push to our repo in docker hub:

```
docker push elminestrone/go-stock-api:1.0
```

## Do the kubernetes
Now that our image is in docker hub we can deploy our kubernetes manifest files. Typically these would be templatized and leveraged with helm charts but due to the short amount of time I've included just the manifest files. To run the deploy use

```
kubectl apply -f manifests/
```

This will deploy a pod that leverages a configmap for environment variables and also a secret for our api key to communicate with the external stock lookup api. Additionally it will set up a nodeport service to expose the service. I used nodeport in this example since I'm running the service locally but normally you would use loadbalancer or use an ingress controller if hosted in the cloud. When running a nodeport service we use the ip of the node and the port that has been assigned to access the service, in this case since I'm running locally it will be localhost and the port that is generated. We can view this by running:

```
kubectl get service forgerock-service
```
From the output we see
```
NAME                TYPE       CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
forgerock-service   NodePort   10.110.193.124   <none>        9000:32553/TCP   2m20s
```
Looking under ports, we see the container is listening on port 9000 but the nodeport service is on 32553 so we should be able to connect using localhost:32553

```
curl localhost:32553
```
And we get back:
```
The 3 day average closing price of MSFT is 213.73
```
