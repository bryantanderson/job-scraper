package org.sincidium.linkd.repositories;

import java.util.List;
import java.util.UUID;

import org.sincidium.linkd.models.Referral;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import jakarta.transaction.Transactional;

@Repository
@Transactional
public interface ReferralRepository extends JpaRepository<Referral, UUID> {
    List<Referral> getAllReferralsByRecruiterId(UUID recruiterId);
    List<Referral> getAllReferralsByCandidateId(UUID candidateId);
}
