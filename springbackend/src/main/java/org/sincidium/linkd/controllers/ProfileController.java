package org.sincidium.linkd.controllers;

import java.util.UUID;

import org.sincidium.linkd.models.Profile;
import org.sincidium.linkd.services.ProfileService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

import io.swagger.v3.oas.annotations.tags.Tag;

@RestController
@RequestMapping("/profiles")
@Tag(name = "Profiles")
public class ProfileController {

    @Autowired
    private ProfileService profileService;

    @GetMapping("/")
    @ResponseStatus(HttpStatus.OK)
    public Profile queryProfile(@RequestParam String uri) {
        return profileService.query(uri);
    }

    @GetMapping("/{userOrProfileId}")
    @ResponseStatus(HttpStatus.OK)
    public Profile getProfile(@PathVariable UUID userOrProfileId) {
        return profileService.get(userOrProfileId);
    }

}
