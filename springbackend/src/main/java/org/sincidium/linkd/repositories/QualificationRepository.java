package org.sincidium.linkd.repositories;

import java.util.UUID;

import org.sincidium.linkd.models.Qualification;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import jakarta.transaction.Transactional;

@Repository
@Transactional
public interface QualificationRepository extends JpaRepository<Qualification, UUID> {
    // DAO for Qualifications
}
