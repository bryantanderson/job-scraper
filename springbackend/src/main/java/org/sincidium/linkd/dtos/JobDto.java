package org.sincidium.linkd.dtos;

import java.util.List;

import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class JobDto {
    private String title;
    private String description;
    private String location;
    private int yearsOfExperience;
    private String locationType;
    private List<String> responsibilities;
    private List<String> qualifications;
}
