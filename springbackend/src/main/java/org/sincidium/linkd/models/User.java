package org.sincidium.linkd.models;

import java.util.ArrayList;
import java.util.Collection;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import jakarta.persistence.CascadeType;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.OneToMany;
import jakarta.persistence.OneToOne;
import jakarta.persistence.PrimaryKeyJoinColumn;
import jakarta.persistence.Table;
import jakarta.validation.constraints.NotEmpty;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@Entity
@Table(name = "users")
@EqualsAndHashCode(callSuper = false)
public class User extends Base {

    @NotEmpty(message = "First name cannot be empty.")
    @Column(name = "first_name", nullable = false, length = 128)
    private String firstName;

    @NotEmpty(message = "Last name cannot be empty.")
    @Column(name = "last_name", nullable = false, length = 128)
    private String lastName;

    @NotEmpty(message = "Email cannot be empty.")
    @Column(name = "email", nullable = true, length = 128)
    private String email;

    @Column(name = "contact_number", nullable = true, length = 128)
    private String contactNumber;

    // PK of the User entity is used as the FK for the associated Profile entity
    @OneToOne(mappedBy = "user", cascade = CascadeType.ALL)
    @PrimaryKeyJoinColumn
    @JsonIgnoreProperties({ "user", "experiences", "education" })
    private Profile profile;

    // Connections initiated by the user.
    @OneToMany(mappedBy = "user", cascade = CascadeType.ALL)
    @JsonIgnoreProperties({ "user" })
    private Collection<Connection> connections = new ArrayList<Connection>();

    // Connections initiated by other users, but accepted by this user.
    @OneToMany(mappedBy = "connectionUser", cascade = CascadeType.ALL)
    @JsonIgnoreProperties({ "connectionUser" })
    private Collection<Connection> acceptedConnections = new ArrayList<Connection>();

    // Referrals initiated by this recruiter, sent to candidates.
    @OneToMany(mappedBy = "recruiter", cascade = CascadeType.ALL)
    @JsonIgnoreProperties({ "recruiter" })
    private Collection<Referral> sentReferrals = new ArrayList<Referral>();

    // Referrals received by this candidate.
    @OneToMany(mappedBy = "candidate", cascade = CascadeType.ALL)
    @JsonIgnoreProperties({ "candidate" })
    private Collection<Referral> receivedReferrals = new ArrayList<Referral>();

    public void setProfile(Profile profile) {
        this.profile = profile;
        // setting the parent class as the value for the child instance
        this.profile.setUser(this);
    }
}
