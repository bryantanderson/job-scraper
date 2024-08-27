package org.sincidium.linkd.models;

import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;
import org.sincidium.linkd.models.assessment.Point;
import org.sincidium.linkd.models.assessment.Score;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.validation.constraints.NotNull;
import lombok.Data;

@Data
@Entity
@Table(name = "assessments")
public class Assessment {

    @Id
    private String id;

    @NotNull
    private int score;

    @NotNull
    private boolean experiencePoint;

    @NotNull
    private boolean locationPoint;

    @NotNull
    @JdbcTypeCode(SqlTypes.JSON)
    @Column(name = "responsibility_point", columnDefinition = "jsonb", nullable = false)
    private Score responsibilityPoint;

    @NotNull
    @JdbcTypeCode(SqlTypes.JSON)
    @Column(name = "requirement_point", columnDefinition = "jsonb", nullable = false)
    private Point requirementPoint;

    @NotNull
    @JdbcTypeCode(SqlTypes.JSON)
    @Column(name = "compatibility_point", columnDefinition = "jsonb", nullable = false)
    private Point compatibilityPoint;

}
