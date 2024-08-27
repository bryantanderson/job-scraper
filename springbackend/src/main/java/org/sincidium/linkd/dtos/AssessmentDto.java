package org.sincidium.linkd.dtos;

import java.util.UUID;

import org.sincidium.linkd.models.Job;
import org.sincidium.linkd.models.Profile;

import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class AssessmentDto {
    private UUID userId;
    private Profile candidate;
    private Job job;
}
