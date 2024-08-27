package org.sincidium.linkd.services;

import java.util.UUID;

import org.sincidium.linkd.dtos.EducationDto;
import org.sincidium.linkd.enums.ResourceEnum;
import org.sincidium.linkd.models.Education;
import org.sincidium.linkd.models.Profile;
import org.sincidium.linkd.repositories.EducationRepository;
import org.sincidium.linkd.utils.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class EducationService implements ICrudService<Education, EducationDto, UUID> {

    @Autowired
    private EducationRepository educationRepository;

    @Autowired
    private ProfileService profileService;

    @Override
    public Education create(EducationDto dto) {
        Education education = parseDto(dto);
        Profile profile = profileService.get(dto.getUserOrProfileId());
        education.setProfile(profile);
        return educationRepository.save(education);
    }

    @Override
    public Education get(UUID id) throws ResourceNotFoundException {
        return educationRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException(ResourceEnum.EDUCATION));
    }

    @Override
    public Education update(UUID id, EducationDto dto) {
        Education education = get(id);
        copyDtoAttributes(education, dto);
        return educationRepository.save(education);
    }

    @Override
    public void delete(UUID id) {
        educationRepository.deleteById(id);
    }

    @Override
    public Education parseDto(EducationDto dto) {
        Education education = new Education();
        copyDtoAttributes(education, dto);
        return education;
    }

    @Override
    public void copyDtoAttributes(Education education, EducationDto dto) {
        education.setTitle(dto.getTitle());
        education.setDescription(dto.getDescription());
        education.setInstitute(dto.getInstitute());
        education.setStartDate(dto.getStartDate());
        education.setEndDate(dto.getEndDate());
    }

}
