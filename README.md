# Building container image and running it locally
- docker-machine start default
- eval $(docker-machine env default)
- go build -o main .
- docker build -t scratch-app -f Dockerfile .
- docker run --publish 7080:7080 -it scratch-app
- Test from docker host:
  - docker-machine ssh default
  - curl localhost:7080/Arya
- OR test from your laptop running docker-machine host
  - docker ps | grep scratch-app
  - docker-machine ip default
  - curl IP_RETRIVED_IN_PREVIOS_STEP:7080/Arya

# Tagging the image and uploading to hub.docker.com
- docker images
- docker tag 0cf5c5dbd51f ideepesh/k8s-demo:v2
- docker push ideepesh/k8s-demo
- docker tag 0adflkadfjlf ideepesh/k8s-demo:latest
- docker push ideepesh/k8s-demo

# Running the image in kubernetes for single pod
- minikube start
- cd basics
- kubectl create -f config/k8sDemoOnePod.yml
- kubectl get pods
- kubectl get pod scratch-app
- kubectl describe pod scratch-app
- Expose service:
  - kubectl expose pod scratch-app --type=NodePort --name scratch-app-pod-svc
  - kubectl get service
  - kubectl describe service scratch-app-pod-svc
  - minikube service scratch-app-pod-svc --url
  - curl URL_RETRIEVED_WITH_PREVIOUS_COMMAND/Arya
- kubectl exec scratch-app -- ANY_COMMAND
- kubectl delete service scratch-app-pod-svc
- kubectl delete pod scratch-app

# Running the image in kubernetes with replication controller
- minikube start
- cd basics
- kubectl create -f config/k8sDemoReplicas.yml
- kubectl get replicationcontrollers
- kubectl describe rc scratch-app-controller
- kubectl get pods
- kubectl scale --replicas=3 -f config/k8sDemoReplicas.yml
- kubectl scale --replicas=4 replicationcontroller scratch-app-controller
- Expose service
  - kubectl expose rc scratch-app-controller --type=NodePort --name scratch-app-rc-svc
  - kubectl get service
  - kubectl describe service scratct-app-rc-svc
  - minikube service scratch-app-rc-svc --url
  - curl URL_RETRIEVED_WITH_PREVIOUS_COMMAND/Arya
- kubectl delete service scratch-app-rc-svc
- kubectl delete replicationcontroller scratch-app-controller

# Running the image in kubernetes with deployment
- minikube start
- cd basics
- kubectl create -f config/k8sDemoDeployment.yml
- kubectl get deployments
- kubectl describe deployment scratch-app-deployment
- kubectl get replicaset
- kubectl describe rs XXXXX
- kubectl get pods
- kubectl 
- Expose service
  - kubectl expose deployment scratch-app-deployment --type=NodePort --name scratch-app-dep-svc
  - kubectl get service
  - kubectl describe service scratct-app-dep-svc
  - minikube service scratch-app-dep-svc --url
  - curl URL_RETRIEVED_WITH_PREVIOUS_COMMAND/Arya
- kubectl set image deployment/scratch-app-deployment k8s-demo=ideepesh/k8s-demo:v1
- kubectl rollout status deployment/scratch-app-deployment
- curl URL_RETRIEVED_WITH_PREVIOUS_COMMAND/Arya
- kubectl rollout undo deployment/scratch-app-deployment
- kubectl rollout status deployment/scratch-app-deployment
- curl URL_RETRIEVED_WITH_PREVIOUS_COMMAND/Arya
- kubectl delete service scratch-app-dep-svc
- kubectl delete deployment scratch-app-deployment

# Creating service from config file
- minikube start
- cd basics
- kubectl create -f config/k8sDemoDeployment.yml
- kubectl get deployments
- kubectl create -f config/k8sDemoSvc.yml
- kubectl get service
- minikube service scratch-app-dep-svc --url
- curl URL_RETRIEVED_WITH_PREVIOUS_COMMAND/Arya
- kubectl delete service scratch-app-dep-svc
- kubectl delete deployment scratch-app-deployment

# Using Node selectors:
- minikube start
- cd basics
- kubectl get nodes
- kubectl get nodes --show-labels
- kubectl create -f config/k8sDemoDepNodeSelector.yml
- kubectl get po
- kubectl get deploy
- kubectl describe po
- kubectl label nodes minikube hardware=commodity
- kubectl get nodes --show-labels
- kubectl get po
- kubectl label nodes minikube hardware=commodity
- kubectl get po
- kubectl delete po ANY_ONE_POD
- kubectl get po
- kubectl delete deployment scratch-app-dep-node-selector

# Adding secrets with volumes
- minikube start
- cd basics
- create -f config/secrets/my-secrets.yml 
- kubectl create -f config/secrets/my-secrets.yml
- kubectl create -f config/k8sDemoDepSecretsVolume.yml
- kubectl get po
- kubectl describe pod scratch-app-deployment-5bf9dc78f7-2hrbh
- kubectl exec scratch-app-deployment-5bf9dc78f7-2hrbh -it -- /bin/bash
  - cd /etc/credentials/
  - cat username 
  - cat password 
  - mount
  - exit
- kubectl delete secrets user-secrets
- kubectl delete deployment scratch-app-dep-secretes-volume

# Running redis in docker container
- docker run --name my-redis-container -p 7001:6379 -d redis
- docker ps | grep my-redis-container
- docker-machine ip default
- redis-cli -h IP_OF_THE_DOCKER_MACHINE_HOST -p 7001
- docker rm -f my-redis-container

# Running the redis and app as separate containers in same pod in kubernetes
- minikube start
- cd appAndRedis
- kubectl create -f config/k8sRedisPod.yml
- kubectl get pods
- kubectl get pod redis-app
- kubectl describe pod redis-app
- kubectl create -f config/k8sRedisSvc.yml
- kubectl get service
- minikube service redis-app-service --url
- redis-cli -h IP_OF_THE_SERVICE -p 32002
- curl APP_URL_RETRIEVED_IN_PREVIOUS_COMMAND:32001/Arya
- kubectl delete service redis-app-service
- kubectl delete pod redis-app

# Running the redis and app in separate pods in kubernetes. Sevice discovery via env variables.
- minikube start
- cd svcDiscEnvVar
- kubectl create -f config/k8sRedisPod.yml
- kubectl describe pod redis-pod
- kubectl create -f config/k8sRedisSvc.yml
- kubectl get service
- kubectl create -f config/k8sRedisAppPod.yml
- kubectl describe pod redis-app-svc-disc
- kubectl logs redis-app-svc-disc 
- kubectl create -f config/k8sRedisAppSvc.yml
- kubectl get service
- minikube service redis-app-service --url
- curl APP_URL_RETRIEVED_IN_PREVIOUS_COMMAND:32001/Arya
- kubectl delete service redis-app-service
- kubectl delete service redis--service
- kubectl delete pod redis-app-svc-disc
- kubectl delete pod redis-pod

# Config map
- minikube start
- cd basics
- kubectl create -f config/configMap/user-config.yaml 
- kubectl create configmap prop-config --from-file=config/configMaps/prop-config.properties 
- kubectl create -f config/k8sDemoDepConfigMapVolume.yml
- kubectl get po
- kubectl describe pod scratch-app-deployment-5bf9dc78f7-2hrbh
- kubectl exec scratch-app-deployment-5bf9dc78f7-2hrbh -it -- /bin/bash
  - cd /etc/config/
  - cat username 
  - cat email 
  - cd /etc/properties/
  - cat prop-config.properties
  - mount
  - exit
- kubectl delete deployment scratch-app-dep-config-map-volume
- kubectl delete configmap user-config
- kubectl delete configmap prop-config

# Using nginx as reverse proxy with configMap/volume specifying nginx configuration
- minikube start
- cd nginx 
- kubectl create -f config/k8sRedisPod.yml
- kubectl describe pod redis-pod
- kubectl create -f config/k8sRedisSvc.yml
- kubectl get service
- kubectl create configmap nginx-config --from-file=config/reverseproxy.conf
- kubectl create -f config/k8sNginxDep.yml
- kubectl create -f config/k8sNginxSvc.yml
- kubectl get service
- minikube service nginx-service --url
- curl APP_URL_RETRIEVED_IN_PREVIOUS_COMMAND:32080/Arya
- kubectl delete service nginx-service
- kubectl delete deployment nginx-dep
- kubectl delete configmap nginx-config
- kubectl delete service redis-service
- kubectl delete pod redis-pod

# Ingress controller with [Nginx Ingress controller](github.com/kubernetes/ingress-nginx)
- cd ingressCongroller
- minikube start
- minikube addons list
- minikube addons enable ingress
- kubectl create -f config/k8sRedisPod.yml
- kubectl create -f config/k8sRedisSvc.yml
- kubectl create -f [config/mandatory.yaml](https://github.com/kubernetes/ingress-nginx/blob/master/deploy/mandatory.yaml)
- kubectl create -f [config/service-nodeport.yaml](https://github.com/kubernetes/ingress-nginx/blob/master/deploy/provider/baremetal/service-nodeport.yaml)
- kubectl create -f config/k8sIngress.yml
- minikube ip
- curl IP_OF_CLUSTER
- curl IP_OF_CLUSTER/Sansa -H 'Host: foo.bar.com'
- curl IP_OF_CLUSTER/Sansa -H 'Host: bar.baz.com'
- curl IP_OF_CLUSTER/Arya -H 'Host: bar.baz.com'
- curl IP_OF_CLUSTER/Arya -H 'Host: foo.bar.com'
- kubectl delete namesapce ingress-nginx
- kubectl delete ingress ingress-resource
- kubectl delete service redis-app-service
- kubectl delete pod redis-app

# Debugging tricks
- k8s: Simple port forwarding on localhost to pod:
  - kubectl port-forward FULL_POD_NAME 6080:7080
  - localhost:6080/Arya
  - Checks if the pod is working
- ssh into the docker-machine host
  - docker-machine ssh default
- Service discovery
  - kubectl get po --namespace=kube-system
  - kubectl run -i --tty busybox --image=busybox ---restart=Never -- sh
  - nslookup NAME_OF_ANY_SERVICE_CLUSTER