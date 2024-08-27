package org.sincidium.linkd.controllers;

import java.util.Arrays;
import java.util.Random;

import org.sincidium.linkd.dtos.ConnectionDto;
import org.sincidium.linkd.dtos.JobDto;
import org.sincidium.linkd.dtos.UserDto;
import org.sincidium.linkd.models.Connection;
import org.sincidium.linkd.models.Job;
import org.sincidium.linkd.models.User;
import org.sincidium.linkd.services.ConnectionService;
import org.sincidium.linkd.services.JobService;
import org.sincidium.linkd.services.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

import io.swagger.v3.oas.annotations.tags.Tag;

/*
 * This controller is a "convenience" controller for testing common development workflows.
 * It will be disabled in non-development environments.
 */

@RestController
@RequestMapping("/test")
@Tag(name = "Test")
public class TestController {

    @Autowired
    private ConnectionService connectionService;

    @Autowired
    private UserService userService;

    @Autowired
    private JobService jobService;

    @PostMapping("/connection")
    @ResponseStatus(HttpStatus.CREATED)
    public Connection createConnection() {
        Random rand = new Random();
        UserDto testUserDto = new UserDto(
            "string", "string", String.valueOf(rand.nextDouble()), "string", "string"
        );
        UserDto testConnectingUserDto = new UserDto(
            "string", "string", String.valueOf(rand.nextDouble()), "string", "string"
        );
        JobDto testJobDto = new JobDto(
            "string", "string", "string", 5, "onsite",
            Arrays.asList("foo", "bar"),
            Arrays.asList("foo", "bar")
        );

        User testUser = userService.create(testUserDto);
        User testConnectingUser = userService.create(testConnectingUserDto);
        Job testJob = jobService.create(testJobDto);

        ConnectionDto testConnectionDto = new ConnectionDto(
            10,
            3,
            testUser.getId(),
            testConnectingUser.getId(),
            testJob.getId(),
            "https://www.linkedin.com/in/fake-user/"
        );

        return connectionService.create(testConnectionDto);
    }
}
