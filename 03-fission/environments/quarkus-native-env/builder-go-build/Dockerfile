## Fission builder base image
ARG BUILDER_IMAGE=fission/builder:latest
FROM ${BUILDER_IMAGE}

ARG GRAAL_VERSION=19.2.1
ARG MAVEN_VERSION=3.5.4
ARG MAVEN_SHA=ce50b1c91364cb77efe3776f756a6d92b76d9038b0a0782f7d53acf1e997a14d
ARG MAVEN_BASE_URL=https://apache.osuosl.org/maven/maven-3/${MAVEN_VERSION}/binaries
ARG USER_HOME_DIR="/root"

WORKDIR /

RUN apk update && apk add --no-cache curl tar bash procps ca-certificates curl

# INSTALL GCC
RUN apk add build-base

# INSTALL GRAALVM
ENV LANG=en_US.UTF-8

ENV GRAALVM_PKG=https://github.com/oracle/graal/releases/download/vm-$GRAAL_VERSION/graalvm-ce-linux-amd64-$GRAAL_VERSION.tar.gz \
    JAVA_HOME=/opt/graalvm-ce-$GRAAL_VERSION/ \
    GRAALVM_HOME=/opt/graalvm-ce-$GRAAL_VERSION/ \
    PATH=$PATH:/opt/rh/llvm-toolset-7/root/usr/bin \
    PATH=${GRAALVM_HOME}/bin:$PATH \
    LD_LIBRARY_PATH=/opt/rh/llvm-toolset-7/root/usr/lib64 \
    MANPATH=/opt/rh/llvm-toolset-7/root/usr/share/man \
    PKG_CONFIG_PATH=/opt/rh/llvm-toolset-7/root/usr/lib64/pkgconfig \
    PYTHONPATH=/opt/rh/llvm-toolset-7/root/usr/lib/python2.7/site-packages \
    X_SCLS=llvm-toolset-7

RUN set -eux \
    && curl --fail --silent --location --retry 3 ${GRAALVM_PKG} \
    | gunzip | tar x -C /opt/ \
    # Set alternative links
    && mkdir -p "/usr/java" \
    && ln -sfT "$JAVA_HOME" /usr/java/default \
    && ln -sfT "$JAVA_HOME" /usr/java/latest

# INSTALL NATIVE IMAGE
RUN ${GRAALVM_HOME}/bin/gu install native-image

# INSTALL MAVEN
RUN mkdir -p /usr/share/maven /usr/share/maven/ref \
    && curl -fsSL -o /tmp/apache-maven.tar.gz ${MAVEN_BASE_URL}/apache-maven-${MAVEN_VERSION}-bin.tar.gz \
    && echo "${MAVEN_SHA}  /tmp/apache-maven.tar.gz" | sha256sum -c - \
    && tar -xzf /tmp/apache-maven.tar.gz -C /usr/share/maven --strip-components=1 \
    && rm -f /tmp/apache-maven.tar.gz \
    && ln -s /usr/share/maven/bin/mvn /usr/bin/mvn

ENV MAVEN_HOME /usr/share/maven
ENV MAVEN_CONFIG "$USER_HOME_DIR/.m2"

## Fission builder specific section
ADD build.sh /usr/local/bin/build
EXPOSE 8001
