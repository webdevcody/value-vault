THIS PROJECT IS A MESS, DO NOT JUDGE ME

I'M WORKING ON THIS TO FORCE MYSELF TO LEARN MORE ABOUT:

- Go
- K8s
- Distributed Databases

## Running Locally

GOOD LUCK, TODO UPDATE ME

`./start-node.sh 0 primary 1 8080 8080`

^ the above command will start a primary node, assuming a cluster size of 1, with id of 0, on port 8080

# Prereq

- must have a DO cluster setup
- must have a DO registry created called `wdc-registry`

# Deploy

1. build container `docker build --platform=linux/amd64 -t registry.digitalocean.com/wdc-registry/key-value-app .`
2. push container `docker push registry.digitalocean.com/wdc-registry/key-value-app`
3. apply changes `kubectl apply -f k8s`

## Running Minikube

1. eval $(minikube docker-env)
1. build container `docker build --platform=linux/amd64 -t registry.digitalocean.com/wdc-registry/key-value-app .`
1. apply changes `kubectl apply -f k8s`
1. setup tunnel `minikube tunnel`
1. access localhost:80 for your service
1. dashboard `minikube dashboard` useful for debugging

## Local Development

<!-- 1. RABBIT_MQ_PASSWORD="BV5QxJAfupW1TZjy" RABBIT_MQ_HOST="localhost" FILE_PATH_PREFIX=./data air -->

1. CONFIG_VERSION=1 FILE_PATH_PREFIX=./data/1 NODES=2 HOSTNAME=localhost:8080 PORT=8080 IS_LOCAL=true air
2. POST@http://localhost:8080/keys/hello
3. GET@http://localhost:8080/keys/hello

## RabbitMQ

1. install rabbitmq

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-rabbitmq bitnami/rabbitmq
```

1. forward ports for dashboard `kubectl port-forward svc/my-rabbitmq 15672:15672`
1. forward this port too `kubectl port-forward --namespace default svc/my-rabbitmq 5672:5672`
1. open http://localhost:15672/

1. username: user
1. password: `$(kubectl get secret --namespace default my-rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 -d)`

## Loadtests

1. install k6
2. k6 run loadtest.js

## Build Container

`eval $(minikube docker-env) && docker build --platform=linux/amd64 -t registry.digitalocean.com/wdc-registry/key-value-app:1 .`

## Deployment Process

1. increase replicas and apply
2. increase nodes and apply
3. do a get request on all keys after new version of pods all up
4. done

## Stuff

CONFIG_VERSION=1 FILE_PATH_PREFIX=./data/1 NODES=1 HOSTNAME=api-primary-0 LOCAL_HOSTNAME=localhost:8080 PORT=8080 IS_LOCAL=true RABBIT_MQ_HOST=localhost RABBIT_MQ_PASSWORD=nuq8W2xD7Xm3lawk MODE=primary air
