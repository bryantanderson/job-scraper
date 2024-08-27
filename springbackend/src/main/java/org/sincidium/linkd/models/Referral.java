package org.sincidium.linkd.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@Entity
@Table(name = "referrals")
@EqualsAndHashCode(callSuper = false)
public class Referral extends Base {

    @Column(name = "link", nullable = false)
    private String link;

    @ManyToOne(fetch = FetchType.EAGER)
    @JoinColumn(name = "candidate_id", nullable = false)
    @JsonIgnoreProperties({"sentReferrals", "receivedReferrals", "connections", "acceptedConnections", "profile", "roles", "password"})
    private User candidate;

    @ManyToOne(fetch = FetchType.EAGER)
    @JoinColumn(name = "recruiter_id", nullable = false)
    @JsonIgnoreProperties({"sentReferrals", "receivedReferrals", "connections", "acceptedConnections", "profile", "roles", "password"})
    private User recruiter;
}
