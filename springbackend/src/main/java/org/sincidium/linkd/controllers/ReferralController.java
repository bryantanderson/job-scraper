package org.sincidium.linkd.controllers;

import java.util.UUID;

import org.sincidium.linkd.dtos.ReferralDto;
import org.sincidium.linkd.models.Referral;
import org.sincidium.linkd.services.ReferralService;
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
@RequestMapping("/referrals")
@Tag(name = "Referrals")
public class ReferralController {

    @Autowired
    private ReferralService referralService;

    @PostMapping("/")
    @ResponseStatus(HttpStatus.CREATED)
    public Referral createReferral(@RequestBody ReferralDto referralDto) {
        return referralService.create(referralDto);
    }

    @GetMapping("/{referralId}")
    @ResponseStatus(HttpStatus.OK)
    public Referral getReferral(@PathVariable UUID referralId) {
        return referralService.get(referralId);
    }

    @PatchMapping("/{referralId}")
    @ResponseStatus(HttpStatus.OK)
    public Referral updateReferral(@PathVariable UUID referralId, @RequestBody ReferralDto dto) {
        return referralService.update(referralId, dto);
    }

    @DeleteMapping("/{referralId}")
    public void deleteReferral(@PathVariable UUID referralId) {
        referralService.delete(referralId);
    }
}
