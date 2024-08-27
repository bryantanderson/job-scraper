
Style guide:

https://google.github.io/styleguide/javaguide.html

Current project structure, in terms of Layers:

[PRESENTATION]
Controllers: Classes that handle HTTP requests, and define API endpoints.
[BUSINESS]
Services: Classes that implement business logic. They may or may not implement an interface.
[PERSISTENCE]
Repositories: Classes responsible for interacting with the database (data access objects).
[DATABASE]
Models: Database Schemas, in POJO (plain old java objects).

Utils: General purpose utility functions, used across the entire app.

To run the app locally without Docker:

mvn clean
mvn spring-boot:run

To package the application using Maven:

./mvnw clean package -DskipTests

Without -DskipTests, maven will try to run the application's tests (which will break for now).
