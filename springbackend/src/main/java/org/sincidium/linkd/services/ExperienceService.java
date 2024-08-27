package org.sincidium.linkd.services;

import java.util.UUID;

import org.sincidium.linkd.dtos.ExperienceDto;
import org.sincidium.linkd.enums.ResourceEnum;
import org.sincidium.linkd.models.Experience;
import org.sincidium.linkd.models.Profile;
import org.sincidium.linkd.repositories.ExperienceRepository;
import org.sincidium.linkd.utils.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class ExperienceService implements ICrudService<Experience, ExperienceDto, UUID> {

    @Autowired
    public ExperienceRepository experienceRepository;

    @Autowired
    private ProfileService profileService;

    @Override
    public Experience create(ExperienceDto dto) {
        Experience experience = parseDto(dto);
        Profile profile = profileService.get(dto.getUserOrProfileId());
        experience.setProfile(profile);
        return experienceRepository.save(experience);
    }

    @Override
    public Experience get(UUID id) throws ResourceNotFoundException {
        return experienceRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException(ResourceEnum.EXPERIENCE));
    }

    @Override
    public Experience update(UUID id, ExperienceDto dto) {
        Experience experience = get(id);
        return experienceRepository.save(experience);
    }

    @Override
    public void delete(UUID id) {
        experienceRepository.deleteById(id);
    }

    @Override
    public Experience parseDto(ExperienceDto dto) {
        Experience experience = new Experience();
        copyDtoAttributes(experience, dto);
        return experience;
    }

    @Override
    public void copyDtoAttributes(Experience model, ExperienceDto dto) {
        model.setStartDate(dto.getStartDate());
        model.setEndDate(dto.getEndDate());
        model.setTitle(dto.getTitle());
        model.setCompany(dto.getCompany());
        model.setDescription(dto.getDescription());
    }

}
