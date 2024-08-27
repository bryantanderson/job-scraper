package org.sincidium.linkd.controllers;

import java.io.IOException;
import java.util.Collection;
import java.util.UUID;

import org.sincidium.linkd.dtos.UserDto;
import org.sincidium.linkd.models.Connection;
import org.sincidium.linkd.models.User;
import org.sincidium.linkd.services.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PatchMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RequestPart;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

import io.swagger.v3.oas.annotations.tags.Tag;

@RestController
@RequestMapping("/users")
@Tag(name = "Users")
public class UserController {

    @Autowired
    private UserService userService;

    @GetMapping("/")
    @ResponseStatus(HttpStatus.OK)
    public User queryUser(@RequestParam String email) {
        return userService.query(email);
    }

    @PostMapping("/")
    @ResponseStatus(HttpStatus.CREATED)
    public User createUser(@RequestBody UserDto UserDto) {
        return userService.create(UserDto);
    }

    @GetMapping("/{userId}")
    @ResponseStatus(HttpStatus.OK)
    public User getUser(@PathVariable UUID userId) {
        return userService.get(userId);
    }

    @PatchMapping("/{userId}")
    @ResponseStatus(HttpStatus.OK)
    public User updateUser(@PathVariable UUID userId, @RequestBody UserDto dto) {
        return userService.update(userId, dto);
    }

    @DeleteMapping("/{userId}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteUser(@PathVariable UUID userId) {
        userService.delete(userId);
    }


    @PostMapping(value = "/{userId}/upload_connections", consumes = {"multipart/form-data"})
    @ResponseStatus(HttpStatus.CREATED)
    public Collection<Connection> uploadConnections(@RequestPart("file") MultipartFile file, @PathVariable UUID userId) throws IOException {
        return userService.parseCSV(file, userId);
    }
}
