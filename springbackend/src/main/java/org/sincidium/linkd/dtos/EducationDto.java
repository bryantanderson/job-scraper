package org.sincidium.linkd.dtos;

import java.util.Date;
import java.util.UUID;

import jakarta.annotation.Nullable;
import lombok.Data;

@Data
public class EducationDto {
    private String title;
    private String description;
    private String institute;

    private Date startDate;

    @Nullable
    private Date endDate;
    private UUID userOrProfileId;
}
