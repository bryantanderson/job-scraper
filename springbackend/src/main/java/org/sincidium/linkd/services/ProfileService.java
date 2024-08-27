package org.sincidium.linkd.services;

import java.util.UUID;

import org.sincidium.linkd.dtos.ProfileDto;
import org.sincidium.linkd.enums.ResourceEnum;
import org.sincidium.linkd.models.Profile;
import org.sincidium.linkd.repositories.ProfileRepository;
import org.sincidium.linkd.utils.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class ProfileService implements ICrudService<Profile, ProfileDto, UUID> {

    @Autowired
    private ProfileRepository profileRepository;

    @Override
    public Profile get(UUID id) throws ResourceNotFoundException {
        return profileRepository.findById(id).orElseThrow(() -> new ResourceNotFoundException(ResourceEnum.PROFILE));
    }

    public Profile query(String uri) {
        return profileRepository.findByUniqueIdentifier(uri).orElseThrow(() -> new ResourceNotFoundException(ResourceEnum.PROFILE));
    }

    @Override
    public Profile create(ProfileDto dto) {
        Profile profile = parseDto(dto);
        return profileRepository.save(profile);
    }

    @Override
    public Profile update(UUID id, ProfileDto dto) {
        Profile profile = get(id);
        copyDtoAttributes(profile, dto);
        return profileRepository.save(profile);
    }

    @Override
    public void delete(UUID id) {
        profileRepository.deleteById(id);
    }

    @Override
    public Profile parseDto(ProfileDto dto) {
        Profile profile = new Profile();
        return profile;
    }

    @Override
    public void copyDtoAttributes(Profile model, ProfileDto dto) {
        return;
    }

}
