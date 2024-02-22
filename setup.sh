#!/bin/bash

minikube start

kubectl create namespace grafana
kubectl apply -f grafana.yaml --namespace=grafana

helm install --values values.yaml loki --namespace=loki grafana/loki --create-namespace

helm upgrade --install promtail grafana/promtail

kubectl apply -f promtail.yaml

kubectl apply -f k8s

helm install my-rabbitmq bitnami/rabbitmq

minikube tunnel # needed to access go load balancer
minikube dashboard # needed for dashboard
kubectl port-forward service/grafana 3000:3000 --namespace=grafana # needed to view logs