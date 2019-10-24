# quarkus-fission-k3s
Demontration using QuarkusIO on Fission K3S

**Deploy a local K3S on Docker:**
```
k3d create --publish 8081:30001@k3d-k3s-default-worker-0 --publish 8080:80@k3d-k3s-default-worker-0  --workers 2
export KUBECONFIG="$(k3d get-kubeconfig --name='k3s-default')"
```

**Initialize Helm:**
```
kubectl apply -f kube-yml/rbac-tiller.yml
helm init --service-account tiller --history-max 200 \n
```

**Install Dashboard:**
```
helm install stable/kubernetes-dashboard --name kubernetes-dashboard --namespace kube-system -f helm-values/dashboard-values.yml
```

**Install Fission:**
```
helm install --name fission --namespace fission --set serviceType=NodePort,routerServiceType=NodePort,prometheusDeploy=false https://github.com/fission/fission/releases/download/1.6.0/fission-all-1.6.0.tgz
```

# Example of function
```
cat > hello.js  <<\EOF
module.exports = function(context, callback) {
    callback(200, "Hello, world!\n");
}
EOF

fission function create --name hello --env nodejs --code hello.js
fission route add --function hello --url /hello
```
