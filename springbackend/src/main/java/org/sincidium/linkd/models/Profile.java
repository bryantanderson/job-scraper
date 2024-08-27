package org.sincidium.linkd.models;

import java.util.ArrayList;
import java.util.Collection;
import java.util.UUID;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import jakarta.persistence.CascadeType;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Index;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.MapsId;
import jakarta.persistence.OneToMany;
import jakarta.persistence.OneToOne;
import jakarta.persistence.Table;
import lombok.Data;

@Data
@Entity
@Table(name = "profiles", indexes = {
    @Index(name = "linked_in_unique_index", columnList = "uniqueIdentifier")
})
public class Profile {

    @Id
    @Column(name = "user_id")
    private UUID id;

    @Column(name = "unique_identifier", nullable = false, unique = true)
    private String uniqueIdentifier;

    @Column(name = "summary", nullable = false)
    private String summary;

    @Column(name = "location", nullable = false)
    private String location;

    // MapsId indicates that the PK values will be copied from the User entity
    @OneToOne
    @MapsId
    @JoinColumn(name = "user_id")
    @JsonIgnoreProperties("profile")
    private User user;

    @OneToMany(cascade = CascadeType.ALL, mappedBy = "profile")
    @JsonIgnoreProperties("profile")
    private Collection<Experience> experiences = new ArrayList<Experience>();

    @OneToMany(cascade = CascadeType.ALL, mappedBy = "profile")
    @JsonIgnoreProperties("profile")
    private Collection<Education> education = new ArrayList<Education>();

    @OneToMany(cascade = CascadeType.ALL, mappedBy = "profile")
    @JsonIgnoreProperties("profile")
    private Collection<Skill> skills = new ArrayList<Skill>();
}
