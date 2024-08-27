package org.sincidium.linkd.dtos;

import java.util.UUID;

import lombok.Data;

@Data
public class ReferralDto {
    private String link;
    private UUID candidateId;
    private UUID recruiterId;
}
