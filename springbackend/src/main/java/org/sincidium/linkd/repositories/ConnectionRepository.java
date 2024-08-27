package org.sincidium.linkd.repositories;

import java.util.UUID;

import org.sincidium.linkd.models.Connection;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import jakarta.transaction.Transactional;

@Repository
@Transactional
public interface ConnectionRepository extends JpaRepository<Connection, UUID> {
    // DAO for Connection
}
