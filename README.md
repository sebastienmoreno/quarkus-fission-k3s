QUARKUS-FISSION-K3S
===================

Demontration using QuarkusIO on Fission K3S

**Deploy a local K3S on Docker:**
```
k3d create --publish 8081:30001@k3d-k3s-default-worker-0 --publish 8080:80@k3d-k3s-default-worker-0  --workers 2
export KUBECONFIG="$(k3d get-kubeconfig --name='k3s-default')"
```

**Initialize Helm:**
```
kubectl apply -f 03-k3s/kube-yml/rbac-tiller.yml
helm init --service-account tiller --history-max 200
```

**Install Dashboard:**
```
helm install stable/kubernetes-dashboard --name kubernetes-dashboard --namespace kube-system -f 03-k3s/helm-values/dashboard-values.yml
```

**Install local storage:**
```
helm install 03-k3s/helm-charts/local-path-provisioner/ --name local-path-storage --namespace local-path-storage
```

> See: https://github.com/rancher/local-path-provisioner

**Install Fission:**

```
helm repo add fission-charts https://fission.github.io/fission-charts/
helm repo update
helm install --name fission --namespace fission -f 03-k3s/helm-values/fission-values.yml fission-charts/fission-all --version 1.6.0
```
> The charts are not available yet in 1.6.0, so use the variant:

```
helm install --name fission --namespace fission -f 03-k3s/helm-values/fission-values.yml 03-k3s/helm-charts/fission-all/
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
