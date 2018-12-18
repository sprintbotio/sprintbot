#!/usr/bin/env bash
pod=$(kubectl --kubeconfig="/Users/kelly/.kube/k8s-sprintbot-kubeconfig.yaml"  get pods --field-selector=status.phase!=Terminating -n sprintbot|grep sprintbot|awk '{print $1}'|head -n 1)
sleep 5
kubectl --kubeconfig="/Users/kelly/.kube/k8s-sprintbot-kubeconfig.yaml" logs -f $pod -n sprintbot