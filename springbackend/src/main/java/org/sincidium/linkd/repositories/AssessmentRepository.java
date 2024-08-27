package org.sincidium.linkd.repositories;

import org.sincidium.linkd.models.Assessment;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import jakarta.transaction.Transactional;

@Repository
@Transactional
public interface AssessmentRepository extends JpaRepository<Assessment, String> {
    // DAO for assessments
}
