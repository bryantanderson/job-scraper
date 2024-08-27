package org.sincidium.linkd.dtos;

import java.util.UUID;

import lombok.Data;

@Data
public class ProfileDto {
    private UUID userId;
    private String uniqueIdentifier;
}
