package org.sincidium.linkd.services;

import java.util.UUID;

import org.sincidium.linkd.dtos.ConnectionDto;
import org.sincidium.linkd.enums.ResourceEnum;
import org.sincidium.linkd.models.Connection;
import org.sincidium.linkd.models.User;
import org.sincidium.linkd.repositories.ConnectionRepository;
import org.sincidium.linkd.utils.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;


@Service
public class ConnectionService implements ICrudService<Connection, ConnectionDto, UUID> {

    @Autowired
    private ConnectionRepository connectionRepository;

    @Autowired
    private UserService userService;

    @Autowired
    private AssessmentService assessmentService;

    @Override
    public Connection create(ConnectionDto dto) {
        Connection connection = parseDto(dto);
        User user = userService.get(dto.getUserId());
        User connectionUser = userService.get(dto.getConnectionUserId());
        connection.setUser(user);
        connection.setConnectionUser(connectionUser);
        assessmentService.publish(dto.getConnectionUserId(), dto.getJobId(), dto.getProfileUrl());
        return connectionRepository.save(connection);
    }

    @Override
    public Connection get(UUID id) throws ResourceNotFoundException {
        return connectionRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException(ResourceEnum.CONNECTION));
    }

    @Override
    public Connection update(UUID id, ConnectionDto dto) {
        Connection connection = get(id);
        copyDtoAttributes(connection, dto);
        return connectionRepository.save(connection);
    }

    @Override
    public void delete(UUID id) {
        connectionRepository.deleteById(id);
    }

    @Override
    public Connection parseDto(ConnectionDto dto) {
        Connection connection = new Connection();
        connection.setDegree(dto.getDegree());
        connection.setRank(dto.getRank());
        return connection;
    }

    @Override
    public void copyDtoAttributes(Connection model, ConnectionDto dto) {
        model.setDegree(dto.getDegree());
        model.setRank(dto.getRank());
    }
}
