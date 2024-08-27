package org.sincidium.linkd.services;

import java.util.List;
import java.util.UUID;

import org.sincidium.linkd.dtos.DummyDto;
import org.sincidium.linkd.enums.ResourceEnum;
import org.sincidium.linkd.models.Dummy;
import org.sincidium.linkd.repositories.DummyRepository;
import org.sincidium.linkd.utils.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class DummyService implements ICrudService<Dummy, DummyDto, UUID> {

    @Autowired
    private final DummyRepository dummyRepository;

    public List<Dummy> getAllDummies() {
        return this.dummyRepository.findAll();
    }

    public Dummy create(DummyDto dto) {
        return dummyRepository.save(parseDto(dto));
    }

    public Dummy get(UUID id) throws ResourceNotFoundException {
        return dummyRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException(ResourceEnum.DUMMY));
    }

    public Dummy update(UUID id, DummyDto dto) {
        Dummy dummy = get(id);
        copyDtoAttributes(dummy, dto);
        return dummyRepository.save(dummy);
    }

    public void delete(UUID id) {
        dummyRepository.deleteById(id);
    }

    public Dummy parseDto(DummyDto dto) {
        Dummy dummy = new Dummy();
        dummy.setName(dto.getName());
        return dummy;
    }

    public void copyDtoAttributes(Dummy dummy, DummyDto dto) {
        dummy.setName(dto.getName());
    }

}
