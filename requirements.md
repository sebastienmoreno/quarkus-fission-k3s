
REQUIREMENTS
============

# For Quarkus

Install JDK:
```
brew cask install java
```
> Note: add the variable  JAVA_HOME=/Library/Java/JavaVirtualMachines/openjdk-13.0.1.jdk

Install Maven:
```
brew install maven
```

## For Fission

## Install Fission CLI

```
curl -Lo fission https://github.com/fission/fission/releases/download/1.6.0/fission-cli-osx \
    && chmod +x fission && sudo mv fission /usr/local/bin/
```

# For K3S

## Install Kubectl
```
brew install kubernetes-cli
```

> https://kubernetes.io/fr/docs/tasks/tools/install-kubectl/#installer-kubectl-sur-macos

## Install Helm
```
brew install kubernetes-helm
```

> https://helm.sh/docs/using_helm/#installing-helm

