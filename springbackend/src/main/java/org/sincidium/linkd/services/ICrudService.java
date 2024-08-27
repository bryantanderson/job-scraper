package org.sincidium.linkd.services;

import org.sincidium.linkd.utils.ResourceNotFoundException;
import org.springframework.stereotype.Service;

@Service
public interface ICrudService<M, D, T> {

    M create(D dto);

    M get(T id) throws ResourceNotFoundException;

    M update(T id, D dto);

    void delete(T id);

    M parseDto(D dto);

    void copyDtoAttributes(M model, D dto);
}
