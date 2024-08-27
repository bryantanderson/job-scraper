package org.sincidium.linkd.repositories;

import java.util.UUID;

import org.sincidium.linkd.models.Education;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import jakarta.transaction.Transactional;

@Repository
@Transactional
public interface EducationRepository extends JpaRepository<Education, UUID> {
    // DAO for Education
}
