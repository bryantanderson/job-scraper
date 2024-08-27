package org.sincidium.linkd.controllers;

import java.util.UUID;

import org.sincidium.linkd.dtos.EducationDto;
import org.sincidium.linkd.models.Education;
import org.sincidium.linkd.services.EducationService;
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
@RequestMapping("/educations")
@Tag(name = "Educations")
public class EducationController {

    @Autowired
    private EducationService educationService;

    @PostMapping("/")
    @ResponseStatus(HttpStatus.CREATED)
    public Education createEducation(@RequestBody EducationDto dto) {
        return educationService.create(dto);
    }

    @GetMapping("/{educationId}")
    @ResponseStatus(HttpStatus.OK)
    public Education getEducation(@PathVariable UUID educationId) {
        return educationService.get(educationId);
    }

    @PatchMapping("/{educationId}")
    @ResponseStatus(HttpStatus.OK)
    public Education updateEducation(@PathVariable UUID educationId, @RequestBody EducationDto dto) {
        return educationService.update(educationId, dto);
    }

    @DeleteMapping("/{educationId}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteEducation(@PathVariable UUID educationId) {
        educationService.delete(educationId);
    }
}
