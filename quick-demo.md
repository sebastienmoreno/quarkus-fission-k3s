QUARKUS-FISSION-K3S
===================

Quick demontration using QuarkusIO on Fission K3S
> Note: Prepare K3D cluster with same yml files

# QUARKUS

```
mvn clean package

java -jar target/rockstar-1.0-SNAPSHOT-runner.jar

curl -w "\n" http://localhost:8080/

curl -X POST -H "Content-Type: application/json;charset=UTF-8" "http://localhost:8080/" -d '{"name":"Jaco Pastorius"}'

```

# K3S

## Deploy Remote K3S cluster

```
export SERVER_IP=54.229.95.135
export AGENT1_IP=34.241.125.164
export AGENT2_IP=54.154.214.146

k3sup install --ip $SERVER_IP --ssh-key ~/.ssh/ec2-training-keypair.pem --user ec2-user

k3sup join --server-ip $SERVER_IP --ip $AGENT1_IP --ssh-key ~/.ssh/ec2-training-keypair.pem --user ec2-user

k3sup join --server-ip $SERVER_IP --ip $AGENT1_IP --ssh-key ~/.ssh/ec2-training-keypair.pem --user ec2-user

export KUBECONFIG=$(PWD)/kubeconfig

# Verification
kubectl get nodes -o wide

```

## Post-installation of cluster

**Initialize Helm:**
```
kubectl apply -f 02-k3s/kube-yml/rbac-tiller.yml
helm init --service-account tiller --history-max 200
```

**Install local storage:**
```
helm install 02-k3s/helm-charts/local-path-provisioner/ --name local-path-storage --namespace local-path-storage
```

**Install Dashboard:**
```
helm install stable/kubernetes-dashboard --name kubernetes-dashboard --namespace kube-system -f 02-k3s/helm-values/dashboard-values.yml

open "https://${SERVER_IP}:30001"
```

# FISSION

**Install Fission:**
```
helm install --name fission --namespace fission -f 03-fission/helm-values/fission-values.yml 03-fission/helm-charts/fission-all/
```

# Quarkus Fission use case

**Install Postgres:**
```
helm install --name postgres --namespace fission-function --set postgresqlPassword=quarkus_test,postgresqlDatabase=quarkus_test,postgresqlUsername=quarkus_test,persistence.storageClass=local-path,fullnameOverride=postgres stable/postgresql
```

**Create environment:**
```
fission env create --name quarkus-native --image smoreno/quarkus-native-env --version 2 --keeparchive=true --builder smoreno/quarkus-native-builder

# Check pods
kubectl -n fission-function get pods
```

**Package and build:**
```
zip -r java-src-pkg.zip *
fission package create --sourcearchive java-src-pkg.zip --env quarkus-native

fission package info --name java-src-pkg-zip-wcok

kubectl -n fission-builder get pods
```

> *Build is too slow, switch to K3D demo*

## Deploy a local K3D cluster

**Deploy a local K3S on Docker:**
```
k3d create --publish 8081:30001@k3d-k3s-default-worker-0 --publish 9080:80@k3d-k3s-default-worker-0  --workers 2
```

**Get Kubeconfig:**
```
export KUBECONFIG="$(k3d get-kubeconfig --name='k3s-default')"
```

**Rockstar tests**
```
fission fn create --name rockstar --pkg java-src-pkg-zip-gxc2 --env quarkus-native

fission fn test --name rockstar

fission route add --function rockstar --url /rockstars --createingress --method GET

fission route add --function rockstar --url /rockstar --createingress --method POST

curl -X GET "http://localhost:9080/rockstars"

curl -X POST -H "Content-Type: application/json;charset=UTF-8" -d '{"name":"Jaco Pastorius"}' "http://localhost:9080/rockstar"
```
