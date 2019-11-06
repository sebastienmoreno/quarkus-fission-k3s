# Fission: Go Environment

This is the Go environment for Fission.

It's a Docker image containing a Go runtime, along with a dynamic loader.

Looking for ready-to-run examples? See the [Go examples directory](../../examples/go).

## Build this image

```
docker build -t smoreno/quarkus-runtime . && docker push smoreno/quarkus-runtime
```

Note that if you build the runtime, you must also build the go-builder
image, to ensure that it's at the same version of go:

```
cd builder && docker build -t smoreno/quarkus-native-builder . && docker push smoreno/quarkus-native-builder
```

## Using the image in fission

You can add this customized image to fission with "fission env
create":

```
fission env create --name go --image smoreno/quarkus-runtime --builder smoreno/quarkus-native-builder --version 2
```

Or, if you already have an environment, you can update its image:

```
fission env update --name go --image smoreno/quarkus-runtime --builder smoreno/quarkus-native-builder
```

After this, fission functions that have the env parameter set to the
same environment name as this command will use this environment.
