QUARKUS-FISSION-K3S
===================

Demontration using QuarkusIO on Fission K3S


# QUARKUS

**Initiate the project**
```
mvn io.quarkus:quarkus-maven-plugin:0.27.0:create \
    -DprojectGroupId=fr.ippon \
    -DprojectArtifactId=vote \
    -DclassName="fr.ippon.vote.GreetingResource" \
    -Dpath="/hello"
```

**Build**
```
./mvnw package
```

**Test**
```
java -jar target/vote-1.0-SNAPSHOT-runner.jar
curl -w "\n" http://localhost:8080/hello
```


# K3S

**Deploy a local K3S on Docker:**
```
k3d create --publish 8081:30001@k3d-k3s-default-worker-0 --publish 8080:80@k3d-k3s-default-worker-0  --workers 2
export KUBECONFIG="$(k3d get-kubeconfig --name='k3s-default')"
```

**Initialize Helm:**
```
kubectl apply -f 02-k3s/kube-yml/rbac-tiller.yml
helm init --service-account tiller --history-max 200
```

**Install Dashboard:**
```
helm install stable/kubernetes-dashboard --name kubernetes-dashboard --namespace kube-system -f 02-k3s/helm-values/dashboard-values.yml
```

**Install local storage:**
```
helm install 02-k3s/helm-charts/local-path-provisioner/ --name local-path-storage --namespace local-path-storage
```

> See: https://github.com/rancher/local-path-provisioner

# FISSION

**Install Fission:**
```
helm install --name fission --namespace fission -f 03-fission/helm-values/fission-values.yml 03-fission/helm-charts/fission-all/
```

> The charts are not available yet in 1.6.0 in official repos, when it does you could use:
```
helm repo add fission-charts https://fission.github.io/fission-charts/
helm repo update
helm install --name fission --namespace fission -f 03-fission/helm-values/fission-values.yml fission-charts/fission-all --version 1.6.0
```

Other type of install with binaries:
```
helm install --name fission --namespace fission --set serviceType=NodePort,routerServiceType=NodePort,prometheusDeploy=false,persistence.storageClass=local-path https://github.com/fission/fission/releases/download/1.6.0/fission-all-1.6.0.tgz
```

# Quarkus Fission use case

helm install --name postgres --namespace fission-function --set postgresqlPassword=quarkus_test,postgresqlDatabase=quarkus_test,postgresqlUsername=quarkus_test,persistence.storageClass=local-path,fullnameOverride=postgres stable/postgresql

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

**Create function:**
```
fission fn create --name hello --pkg java-src-pkg-zip-ludc --env quarkus-native

fission fn test --name hello

fission route add --function hello --url /hello --createingress

curl http://localhost:8080/hello
```

**Rockstar**
```
fission fn create --name rockstar --pkg java-src-pkg-zip-ludc --env quarkus-native

fission fn test --name rockstar

fission route add --function rockstar --url /rockstars --createingress --method GET

fission route add --function rockstar --url /rockstar --createingress --method POST

curl -X GET "http://localhost:8080/rockstars"

curl -X POST -H "Content-Type: application/json;charset=UTF-8" "http://localhost:8080/rockstar" -d '{"name":"Led Zeppelin"}'

```
