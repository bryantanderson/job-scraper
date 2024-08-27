package org.sincidium.linkd.repositories;

import java.util.UUID;

import org.sincidium.linkd.models.Job;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import jakarta.transaction.Transactional;

@Repository
@Transactional
public interface JobRepository extends JpaRepository<Job, UUID> {
    // DAO for Jobs
}
