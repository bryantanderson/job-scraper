package org.sincidium.linkd.models;

import java.util.ArrayList;
import java.util.Collection;

import org.sincidium.linkd.enums.JobLocationTypeEnum;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import jakarta.persistence.CascadeType;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.OneToMany;
import jakarta.persistence.Table;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@Entity
@Table(name = "jobs")
@EqualsAndHashCode(callSuper = false)
public class Job extends Base {

    @NotEmpty(message = "Title cannot be empty.")
    @Column(name = "title", nullable = false, length = 256)
    private String title;

    @NotBlank(message = "Description cannot be empty")
    @Column(name = "description", nullable = false)
    private String description;

    @NotBlank(message = "Location cannot be empty")
    @Column(name = "location", nullable = false)
    private String location;

    @Column(name = "location_type", nullable = false)
    @Enumerated(EnumType.STRING)
    private JobLocationTypeEnum locationType;

    @Column(name = "years_of_experience", nullable = false)
    private Integer yearsOfExperience;

    @OneToMany(cascade = CascadeType.ALL, mappedBy = "job")
    @JsonIgnoreProperties({"job"})
    private Collection<Qualification> qualifications = new ArrayList<>();

    @OneToMany(cascade = CascadeType.ALL, mappedBy = "job")
    @JsonIgnoreProperties({"job"})
    private Collection<Responsibility> responsibilities = new ArrayList<>();

}
