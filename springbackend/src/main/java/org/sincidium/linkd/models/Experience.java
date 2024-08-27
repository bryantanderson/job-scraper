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
@Table(name = "experiences")
@EqualsAndHashCode(callSuper = false)
public class Experience extends Base {

    @NotEmpty(message = "Start date cannot be empty.")
    @Column(name = "start_date", nullable = false)
    private Date startDate;

    @Nullable
    @Column(name = "end_date", nullable = true)
    private Date endDate;

    @NotEmpty(message = "Title cannot be empty.")
    @Column(name = "title", nullable = false, length = 128)
    private String title;

    @NotEmpty(message = "Company cannot be empty.")
    @Column(name = "company", nullable = false, length = 128)
    private String company;

    @Column(name = "description", nullable = true)
    private String description;

    @ManyToOne
    @JoinColumn(name = "profile_id", nullable = false)
    @JsonIgnoreProperties({"experiences", "user", "education", "skills"})
    private Profile profile;

}
