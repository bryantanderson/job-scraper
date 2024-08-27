package org.sincidium.linkd.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;
import jakarta.validation.constraints.NotBlank;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@Entity
@Table(name = "responsibilities")
@AllArgsConstructor
@EqualsAndHashCode(callSuper = false)
public class Responsibility extends Base {

    @NotBlank(message = "Description cannot be empty")
    @Column(name = "description", nullable = false)
    private String description;

    @ManyToOne
    @JoinColumn(name = "job_id", nullable = false)
    @JsonIgnoreProperties({ "qualifications", "responsibilities" })
    private Job job;

}
