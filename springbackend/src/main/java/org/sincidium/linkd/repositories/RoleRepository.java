package org.sincidium.linkd.repositories;

import java.util.UUID;

import org.sincidium.linkd.models.Role;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import jakarta.transaction.Transactional;

@Repository
@Transactional
public interface RoleRepository extends JpaRepository<Role, UUID> {
    // DAO for Role
    Role findByName(String name);
}
