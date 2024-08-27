package org.sincidium.linkd.controllers;

import java.util.UUID;

import org.sincidium.linkd.dtos.ExperienceDto;
import org.sincidium.linkd.models.Experience;
import org.sincidium.linkd.services.ExperienceService;
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
@RequestMapping("/experiences")
@Tag(name = "Experiences")
public class ExperienceController {

    @Autowired
    private ExperienceService experienceService;

    @PostMapping("/")
    @ResponseStatus(HttpStatus.CREATED)
    public Experience createExperience(@RequestBody ExperienceDto dto) {
        return experienceService.create(dto);
    }

    @GetMapping("/{experienceId}")
    @ResponseStatus(HttpStatus.OK)
    public Experience getExperience(@PathVariable UUID experienceId) {
        return experienceService.get(experienceId);
    }

    @PatchMapping("/{experienceId}")
    @ResponseStatus(HttpStatus.OK)
    public Experience updateExperience(@PathVariable UUID experienceId, @RequestBody ExperienceDto dto) {
        return experienceService.update(experienceId, dto);
    }

    @DeleteMapping("/{experienceId}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteExperience(@PathVariable UUID experienceId) {
        experienceService.delete(experienceId);
    }

}
