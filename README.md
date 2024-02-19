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

1. RABBIT_MQ_PASSWORD="BV5QxJAfupW1TZjy" RABBIT_MQ_HOST="localhost" FILE_PATH_PREFIX=./data air
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
