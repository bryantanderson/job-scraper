package org.sincidium.linkd.repositories;

import java.util.UUID;

import org.sincidium.linkd.models.Experience;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import jakarta.transaction.Transactional;

@Repository
@Transactional
public interface ExperienceRepository extends JpaRepository<Experience, UUID> {
    // DAO for Experience
}
