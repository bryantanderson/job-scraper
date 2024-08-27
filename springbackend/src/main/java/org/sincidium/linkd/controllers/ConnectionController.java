package org.sincidium.linkd.controllers;

import java.util.UUID;

import org.sincidium.linkd.dtos.ConnectionDto;
import org.sincidium.linkd.models.Connection;
import org.sincidium.linkd.services.ConnectionService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PatchMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

import io.swagger.v3.oas.annotations.tags.Tag;

@RestController
@RequestMapping("/connections")
@Tag(name = "Connections")
public class ConnectionController {

    @Autowired
    private ConnectionService connectionService;

    @PostMapping("/")
    @ResponseStatus(HttpStatus.CREATED)
    public Connection createConnection(@RequestBody ConnectionDto dto) {
        return connectionService.create(dto);
    }

    @GetMapping("/{connectionId}")
    @ResponseStatus(HttpStatus.OK)
    public Connection getConnection(@PathVariable UUID connectionId) {
        return connectionService.get(connectionId);
    }

    @PatchMapping("/{connectionId}")
    @ResponseStatus(HttpStatus.OK)
    public Connection updateConnection(@PathVariable UUID connectionId, @RequestBody ConnectionDto dto) {
        return connectionService.update(connectionId, dto);
    }

    @DeleteMapping("/{connectionId}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteConnection(@PathVariable UUID connectionId) {
        connectionService.delete(connectionId);
    }

}
