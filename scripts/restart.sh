#!/usr/bin/env bash
kubectl --kubeconfig="/Users/kelly/.kube/k8s-sprintbot-kubeconfig.yaml" scale deployment sprintbot --replicas=0 -n sprintbot
sleep 1
kubectl --kubeconfig="/Users/kelly/.kube/k8s-sprintbot-kubeconfig.yaml" scale deployment sprintbot --replicas=1 -n sprintbot
