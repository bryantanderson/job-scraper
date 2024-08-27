package org.sincidium.linkd.controllers;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

import io.swagger.v3.oas.annotations.tags.Tag;

@RestController
@RequestMapping("/health")
@Tag(name = "Health")
public class HealthController {

    @GetMapping("/")
    @ResponseStatus(HttpStatus.OK)
    public ResponseEntity<String> getHealthCheck() {
        return ResponseEntity.ok("OK");
    }
}
