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


```
zip -r quarkus-function-pkg.zip lib vote-1.0-SNAPSHOT-runner.jar

fission env create --name java --image fission/jvm-env --version 2 --keeparchive=true
fission package create --sourcearchive quarkus-function-pkg.zip --env java

fission fn create --name hello --pkg quarkus-function-pkg-zip-vmug --env java --entrypoint fr.ippon.vote.GreetingResource

fission fn test --name hello2

```

CUSTOM
```
fission env create --name jvm --image smoreno/jvm-env --version 2 --keeparchive --builder smoreno/jvm-builder
2736  fission env create --name jvm --image smoreno/jvm-env --version 2 --keeparchive --builder smoreno/jvm-builder
fission env list
cd 01-quarkus/vote
zip -Tv java-src-pkg.zip
fission package create --sourcearchive java-src-pkg.zip --env jvm
fission package info --name java-src-pkg-zip-iwyg
fission fn create --name hello2 --pkg java-src-pkg-zip-iwyg --env jvm --entrypoint fr.ippon.vote.GreetingResource

fission fn test --name hello2

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
helm repo add fission-charts https://fission.github.io/fission-charts/
helm repo update
helm install --name fission --namespace fission -f 03-fission/helm-values/fission-values.yml fission-charts/fission-all --version 1.6.0
```
> The charts are not available yet in 1.6.0, so use the variant:

```
helm install --name fission --namespace fission -f 03-fission/helm-values/fission-values.yml 03-fission/helm-charts/fission-all/
```

OR

```
helm install --name fission --namespace fission --set serviceType=NodePort,routerServiceType=NodePort,prometheusDeploy=false,persistence.storageClass=local-path https://github.com/fission/fission/releases/download/1.6.0/fission-all-1.6.0.tgz
```


# Example of function
```
cat > hello.js  <<\EOF
module.exports = function(context, callback) {
    callback(200, "Hello, world!\n");
}
EOF

fission function create --name hello --env nodejs --code hello.js
fission route add --function hello --url /ihello --createingress
curl http://localhost:8080/ihello
```
