package org.sincidium.linkd.config;

import java.util.ArrayList;
import java.util.List;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import io.swagger.v3.oas.models.security.SecurityRequirement;
import io.swagger.v3.oas.models.Components;
import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.security.SecurityScheme;
import io.swagger.v3.oas.models.security.SecurityScheme.Type;
import io.swagger.v3.oas.models.servers.Server;

// http://localhost:8080/swagger-ui/index.html

@Configuration
public class OpenApiConfig {

    private SecurityScheme createAPIKeyScheme() {
        return new SecurityScheme().type(Type.HTTP)
                .bearerFormat("JWT")
                .scheme("bearer");
    }

    @Bean
    public OpenAPI defineOpenAPI() {
        List<Server> servers = new ArrayList<Server>();
        Server server = new Server();
        server.setUrl("http://localhost:8080");
        server.setDescription("Development Set-up");
        servers.add(server);

        // Configure security fields in swagger UI
        SecurityRequirement security = new SecurityRequirement();
        security.addList("Bearer Authentication");

        Contact contact = new Contact();
        contact.setName("Taha");
        contact.setEmail("tahaaansari@hotmail.com");

        Info information = new Info();
        information.setTitle("Linkd");
        information.setVersion("1.0");
        information.setDescription("Swagger Docs for our Spring Boot based API");
        information.setContact(contact);

        return new OpenAPI().addSecurityItem(security)
                .components(new Components().addSecuritySchemes("Bearer Authentication", createAPIKeyScheme()))
                .info(information).servers(servers);
    }
}
