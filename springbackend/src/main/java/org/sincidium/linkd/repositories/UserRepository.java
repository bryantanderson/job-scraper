package org.sincidium.linkd.repositories;

import java.util.Optional;
import java.util.UUID;

import org.sincidium.linkd.models.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import jakarta.transaction.Transactional;

@Repository
@Transactional
public interface UserRepository extends JpaRepository<User, UUID> {
    // DAO for user
    Optional<User> findByEmail(String email);
}
