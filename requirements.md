
REQUIREMENTS
============

> In order to play the demo you will have to install the following tools:

# For Quarkus

**Install JDK:**
```
brew cask install java
```

> Note: add the variable: export JAVA_HOME=/Library/Java/JavaVirtualMachines/openjdk-13.0.1.jdk/Contents/Home

**Install Maven:**
```
brew install maven
```

## Optional: native build On Mac Os

**Install XCode:**
```
xcode-select --install
```
**Install GraalVM**:
```
brew cask install graalvm/tap/graalvm-ce
export GRAALVM_HOME=/Library/Java/JavaVirtualMachines/graalvm-ce-19.2.1/Contents/Home
export JAVA_HOME=/Library/Java/JavaVirtualMachines/graalvm-ce-19.2.1/Contents/Home
export PATH=/Library/Java/JavaVirtualMachines/graalvm-ce-19.2.1/Contents/Home/bin:"$PATH"
```
**Install native Image**:
```
${GRAALVM_HOME}/bin/gu install native-image
```

> Note: the Maven profile "native" will now build additionally a "runner" native binary 

# For Fission

## Install Fission CLI

```
curl -Lo fission https://github.com/fission/fission/releases/download/1.6.0/fission-cli-osx \
    && chmod +x fission && sudo mv fission /usr/local/bin/
```

# For K3S

## Install Kubectl

**Linux:**
```
curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl

chmod +x ./kubectl
```

**Mac OS:**
```
brew install kubernetes-cli
```

> https://kubernetes.io/fr/docs/tasks/tools/install-kubectl/#installer-kubectl-sur-macos

## Install Helm

**Linux:**
```
curl -L https://git.io/get_helm.sh | bash
```

**Mac OS:**
```
brew install kubernetes-helm
```

> https://helm.sh/docs/using_helm/#installing-helm

