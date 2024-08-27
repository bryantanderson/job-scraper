package org.sincidium.linkd.models;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotEmpty;
import lombok.Data;
import lombok.EqualsAndHashCode;

@Data
@Entity
@Table(name = "dummy")
@EqualsAndHashCode(callSuper = false)
public class Dummy extends Base {

    @NotEmpty(message = "Name cannot be empty.")
    @Column(name = "name", nullable = false)
    private String name;

}
