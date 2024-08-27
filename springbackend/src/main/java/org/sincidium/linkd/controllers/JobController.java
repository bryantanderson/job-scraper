package org.sincidium.linkd.controllers;

import java.util.UUID;

import org.sincidium.linkd.dtos.JobDto;
import org.sincidium.linkd.models.Job;
import org.sincidium.linkd.services.JobService;
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
@RequestMapping("/jobs")
@Tag(name = "Jobs")
public class JobController {

    @Autowired
    private JobService jobService;

    @PostMapping("/")
    @ResponseStatus(HttpStatus.CREATED)
    public Job createJob(@RequestBody JobDto dto) {
        return jobService.create(dto);
    }

    @GetMapping("/{jobId}")
    @ResponseStatus(HttpStatus.OK)
    public Job getJob(@PathVariable UUID jobId) {
        return jobService.get(jobId);
    }

    @PatchMapping("/{jobId}")
    @ResponseStatus(HttpStatus.OK)
    public Job updateJob(@PathVariable UUID jobId, @RequestBody JobDto dto) {
        return jobService.update(jobId, dto);
    }

    @DeleteMapping("/{jobId}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteJob(@PathVariable UUID jobId) {
        jobService.delete(jobId);
    }
}
