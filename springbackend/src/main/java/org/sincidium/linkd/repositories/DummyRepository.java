package org.sincidium.linkd.repositories;

import java.util.UUID;

import org.sincidium.linkd.models.Dummy;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import jakarta.transaction.Transactional;

@Repository
@Transactional
public interface DummyRepository extends JpaRepository<Dummy, UUID> {
    // DAO for Dummy
}
