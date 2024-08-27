package org.sincidium.linkd.dtos;

import jakarta.annotation.Nullable;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;

@Data
@NoArgsConstructor
@RequiredArgsConstructor
public class UserDto {
    @NonNull
    private String firstName;

    @NonNull
    private String lastName;

    @Nullable
    private String email;

    @Nullable
    private String contactNumber;

    @NonNull
    private String uri;

    @NonNull
    private String summary;

    @NonNull
    private String location;
}
