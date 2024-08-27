FROM jelastic/maven:3.9.4-openjdk-22.ea-b17 AS build

WORKDIR /app/src/springbackend

COPY ./pom.xml ./pom.xml

# go-offline using the pom.xml
RUN mvn dependency:go-offline

COPY ./src ./src

RUN mvn clean install -DskipTests

FROM openjdk:22 AS dev

COPY --from=build /app/src/springbackend/target/*.jar /app/src/springbackend/app.jar

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "-Dspring.devtools.restart.enabled=true", "-Dspring.profiles.active=dev", "/app/src/springbackend/app.jar"]
