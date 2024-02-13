# Prereq

- must have a DO cluster setup
- must have a DO registry created called `wdc-registry`

# Deploy

1. build container `docker build --platform=linux/amd64 -t registry.digitalocean.com/wdc-registry/key-value-app .`
2. push container `docker push registry.digitalocean.com/wdc-registry/key-value-app`
3. apply changes `kubectl apply -f k8s`
