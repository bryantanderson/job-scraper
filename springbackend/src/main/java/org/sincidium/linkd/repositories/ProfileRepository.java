package org.sincidium.linkd.repositories;

import java.util.Optional;
import java.util.UUID;

import org.sincidium.linkd.models.Profile;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import jakarta.transaction.Transactional;

@Repository
@Transactional
public interface ProfileRepository extends JpaRepository<Profile, UUID> {
    // DAO for Profile
    Optional<Profile> findByUniqueIdentifier(String uri);
}
