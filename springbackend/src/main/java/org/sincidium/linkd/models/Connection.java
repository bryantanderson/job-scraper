package org.sincidium.linkd.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;
import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotNull;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@Entity
@Table(name = "connections")
@EqualsAndHashCode(callSuper = false)
public class Connection extends Base {

    @NotNull
    @Column(name = "rank", nullable = true)
    private int rank;

    @NotNull
    @Min(1)
    @Max(3)
    @Column(name = "degree", nullable = true)
    private int degree;

    @ManyToOne(fetch = FetchType.EAGER)
    @JoinColumn(name = "user_id", nullable = false)
    @JsonIgnoreProperties({"connections", "acceptedConnections", "sentReferrals", "receivedReferrals", "profile", "roles", "password"})
    private User user;

    @ManyToOne(fetch = FetchType.EAGER)
    @JoinColumn(name = "connection_user_id", nullable = false)
    @JsonIgnoreProperties({"connections", "acceptedConnections", "sentReferrals", "receivedReferrals", "profile", "roles", "password"})
    private User connectionUser;
}
