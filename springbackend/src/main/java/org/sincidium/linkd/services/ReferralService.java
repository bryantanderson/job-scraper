package org.sincidium.linkd.services;

import java.util.List;
import java.util.UUID;

import org.sincidium.linkd.dtos.ReferralDto;
import org.sincidium.linkd.enums.ResourceEnum;
import org.sincidium.linkd.models.Referral;
import org.sincidium.linkd.models.User;
import org.sincidium.linkd.repositories.ReferralRepository;
import org.sincidium.linkd.utils.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class ReferralService implements ICrudService<Referral, ReferralDto, UUID> {

    @Autowired
    private ReferralRepository referralRepository;

    @Autowired
    private UserService userService;

    public List<Referral> getAll() {
        return referralRepository.findAll();
    }

    @Override
    public Referral create(ReferralDto dto) {
        Referral referral = parseDto(dto);
        User candidate = userService.get(dto.getCandidateId());
        User recruiter = userService.get(dto.getRecruiterId());
        referral.setCandidate(candidate);
        referral.setRecruiter(recruiter);
        return referralRepository.save(referral);
    }

    @Override
    public Referral get(UUID id) throws ResourceNotFoundException {
        return referralRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException(ResourceEnum.REFERRAL));
    }

    @Override
    public Referral update(UUID id, ReferralDto dto) {
        Referral referral = get(id);
        copyDtoAttributes(referral, dto);
        return referralRepository.save(referral);
    }

    @Override
    public void delete(UUID id) {
        referralRepository.deleteById(id);
    }

    @Override
    public Referral parseDto(ReferralDto dto) {
        Referral referral = new Referral();
        copyDtoAttributes(referral, dto);
        return referral;
    }

    @Override
    public void copyDtoAttributes(Referral referral, ReferralDto dto) {
        referral.setLink(dto.getLink());
    }
}
