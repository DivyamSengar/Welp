# Kubernetes Tutorial

In this tutorial, we will guide you through a simple kubernetes example.

## Prerequisites
Make sure to install `kubeadm` and `kubectl` using the provided installation script.

## Kubernetes
1. 
    Create a simple deployment consisting of a containerized application and define the desired state of your application. Use the following `nginx-deployment.yaml` file as an example:

    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
    name: nginx-deployment
    spec:
    replicas: 3
    selector:
        matchLabels:
        app: nginx
    template:
        metadata:
        labels:
            app: nginx
        spec:
        containers:
            - name: nginx
              image: nginx:latest
              ports:
                - containerPort: 80
    ```
    

    This YAML code defines a Kubernetes Deployment object for running multiple replicas of an NGINX web server. This YAML code, when applied using `kubectl apply -f nginx-deployment.yaml`, will create a Deployment named nginx-deployment with 3 replicas of the NGINX web server, each running on port 80.

2. Apply the deployment
    ```bash
    kubectl apply -f nginx-deployment.yaml
    ```
3. Verify Deployment and Pod Status
    ```bash
    kubectl get deployment nginx-deployment
    kubectl get pods -l app=nginx
    ```
4. Scale Application
    ```bash
    kubectl scale deployment nginx-deployment --replicas=5
    ```
5. Confirm scaling 
    ```bash
    kubectl get deployment nginx-deployment
    kubectl get pods -l app=nginx
    ```
6. Delete a pod of your choice
    ```bash
    kubectl delete pod <pod-name>
    ```
7. Verify the pod deletion
    ```bash
    kubectl get pods -l app=nginx
    ```
