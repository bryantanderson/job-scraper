package org.sincidium.linkd.dtos;

import java.util.Date;
import java.util.UUID;

import io.micrometer.common.lang.Nullable;
import lombok.Data;

@Data
public class ExperienceDto {
    private Date startDate;

    @Nullable
    private Date endDate;

    private String title;
    private String company;

    @Nullable
    private String description;
    private UUID userOrProfileId;
}
