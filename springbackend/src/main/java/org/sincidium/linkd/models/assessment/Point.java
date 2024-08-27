package org.sincidium.linkd.models.assessment;

import java.io.Serializable;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class Point implements Serializable {
    private String explanation;
    private boolean isValid;
}
