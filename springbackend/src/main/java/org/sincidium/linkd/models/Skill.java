package org.sincidium.linkd.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;
import jakarta.validation.constraints.NotBlank;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@Entity
@Table(name = "skills")
@EqualsAndHashCode(callSuper = false)
public class Skill extends Base {

    @NotBlank(message = "Description cannot be empty")
    @Column(name = "description", nullable = false)
    private String description;

    @ManyToOne
    @JoinColumn(name = "profile_id", nullable = false)
    @JsonIgnoreProperties({"experiences", "user", "education", "skills"})
    private Profile profile;

}
