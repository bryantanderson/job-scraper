package org.sincidium.linkd.controllers;

import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

import io.swagger.v3.oas.annotations.tags.Tag;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;

import java.util.List;
import java.util.UUID;

import org.sincidium.linkd.dtos.DummyDto;
import org.sincidium.linkd.models.Dummy;
import org.sincidium.linkd.services.DummyService;

@RestController
@RequestMapping("/dummies")
@Tag(name = "Dummies")
public class DummyController {

    @Autowired
    private DummyService dummyService;

    @GetMapping("/")
    @ResponseStatus(HttpStatus.OK)
    public List<Dummy> getAllDummies() {
        return dummyService.getAllDummies();
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public Dummy createDummy(@RequestBody DummyDto dummyDto) {
        return dummyService.create(dummyDto);
    }

    @GetMapping("/{dummyId}")
    @ResponseStatus(HttpStatus.OK)
    public Dummy getDummy(@PathVariable UUID dummyId) {
        return dummyService.get(dummyId);
    }

    @DeleteMapping("/{dummyId}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteDummy(@PathVariable UUID dummyId) {
        dummyService.delete(dummyId);
    }
}
