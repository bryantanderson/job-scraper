package org.sincidium.linkd.models.assessment;

import java.io.Serializable;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class Score implements Serializable {
    private int score;
    private String explanation;
}
