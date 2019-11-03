FROM maven:3.5-jdk-8 as BUILD
WORKDIR /usr/src/myapp/

# To reuse the build cache, here we split maven dependency
# download and package into two RUN commands to avoid cache invalidation.
COPY pom.xml .
RUN mvn dependency:go-offline

COPY src /usr/src/myapp/src/
RUN mvn package

FROM fabric8/java-alpine-openjdk8-jre
ENV JAVA_OPTIONS="-Dquarkus.http.port=8888 -Dquarkus.http.host=0.0.0.0 -Djava.util.logging.manager=org.jboss.logmanager.LogManager"
ENV AB_ENABLED=jmx_exporter
COPY --from=BUILD /usr/src/myapp/target/lib/* /deployments/lib/
COPY --from=BUILD /usr/src/myapp/target/*-runner.jar /deployments/app.jar
EXPOSE 8888

# run with user 1001 and be prepared for be running in OpenShift too
RUN adduser -G root --no-create-home --disabled-password 1001 \
  && chown -R 1001 /deployments \
  && chmod -R "g+rwX" /deployments \
  && chown -R 1001:root /deployments
USER 1001

ENTRYPOINT [ "/deployments/run-java.sh" ]