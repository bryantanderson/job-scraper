package org.sincidium.linkd.models;

import java.util.Date;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import io.micrometer.common.lang.Nullable;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;
import jakarta.validation.constraints.NotEmpty;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@Entity
@Table(name = "educations")
@EqualsAndHashCode(callSuper = false)
public class Education extends Base {

    @Column(name = "start_date", nullable = false)
    private Date startDate;

    @Nullable
    @Column(name = "end_date", nullable = true)
    private Date endDate;

    @NotEmpty(message = "Title cannot be empty.")
    @Column(name = "title", nullable = false, length = 256)
    private String title;

    @NotEmpty(message = "Institute cannot be empty.")
    @Column(name = "institute", nullable = false)
    private String institute;

    @Column(name = "description", nullable = true)
    private String description;

    @ManyToOne
    @JoinColumn(name = "profile_id", nullable = false)
    @JsonIgnoreProperties({"education", "user", "experiences", "skills"})
    private Profile profile;

}
