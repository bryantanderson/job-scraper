package org.sincidium.linkd.dtos;

import java.util.UUID;

import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class ConnectionDto {
    private int rank;

    @Min(1)
    @Max(3)
    private int degree;
    private UUID userId;
    private UUID connectionUserId;
    private UUID jobId;
    private String profileUrl;
}
