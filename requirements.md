
REQUIREMENTS
============

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

## optional: native build

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

