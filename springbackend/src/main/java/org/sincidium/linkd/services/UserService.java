package org.sincidium.linkd.services;

import java.io.BufferedReader;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.util.Collection;
import java.util.Optional;
import java.util.UUID;
import java.util.stream.Collectors;

import org.sincidium.linkd.dtos.UserDto;
import org.sincidium.linkd.enums.ResourceEnum;
import org.sincidium.linkd.models.Connection;
import org.sincidium.linkd.models.Profile;
import org.sincidium.linkd.models.User;
import org.sincidium.linkd.repositories.ProfileRepository;
import org.sincidium.linkd.repositories.ConnectionRepository;
import org.sincidium.linkd.repositories.UserRepository;
import org.sincidium.linkd.utils.ResourceNotFoundException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.server.ResponseStatusException;

import joinery.DataFrame;

@Service
public class UserService implements ICrudService<User, UserDto, UUID> {

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private ProfileRepository profileRepository;

    @Autowired
    private ConnectionRepository connectionRepository;

    @Override
    public User create(UserDto dto) {
        // Check whether the unique identifier is already taken
        if (dto.getEmail() != null) {
            Optional<Profile> existingProfileWithUniqueIdentifier = profileRepository
                    .findByUniqueIdentifier(dto.getUri());
            if (existingProfileWithUniqueIdentifier.isPresent()) {
                throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Unique profile identifier already taken");
            }
        }
        User user = parseDto(dto);
        Profile profile = new Profile();
        profile.setSummary(dto.getSummary());
        profile.setLocation(dto.getLocation());
        profile.setUniqueIdentifier(dto.getUri());
        user.setProfile(profile);
        return userRepository.save(user);
    }

    public User query(String email) {
        // Check whether the email is already taken
        return userRepository.findByEmail(email).orElseThrow(
                () -> new ResourceNotFoundException(ResourceEnum.USER));
    }

    @Override
    public User get(UUID id) throws ResourceNotFoundException {
        return userRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException(ResourceEnum.USER));
    }

    @Override
    public User update(UUID id, UserDto dto) {
        User user = get(id);
        copyDtoAttributes(user, dto);
        return userRepository.save(user);
    }

    @Override
    public void delete(UUID id) {
        userRepository.deleteById(id);
    }

    @Override
    public User parseDto(UserDto dto) {
        User user = new User();
        copyDtoAttributes(user, dto);
        return user;
    }

    @Override
    public void copyDtoAttributes(User user, UserDto dto) {
        user.setFirstName(dto.getFirstName());
        user.setLastName(dto.getLastName());
        user.setEmail(dto.getEmail());
        user.setContactNumber(dto.getContactNumber());
    }

    public Collection<Connection> parseCSV(MultipartFile file, UUID userId) throws IOException {
        User curUser = get(userId);
        // Delete the note lines that LinkedIn leaves in their connection files
        InputStream inputStream = removeFirstThreeLines(file);

        // Parse the remaining content into a df
        DataFrame<Object> df = DataFrame.readCsv(inputStream);

        // Create new User, Profile and Connection
        for (int i = 0; i < df.length(); i++) {

            // New User
            String firstName = df.col("First Name").get(i) != null ? df.col("First Name").get(i).toString() : null;
            String lastName = df.col("Last Name").get(i) != null ? df.col("Last Name").get(i).toString() : null;
            String email = df.col("Email Address").get(i) != null ? df.col("Email Address").get(i).toString() : null;
            String uri = df.col("URL").get(i) != null ? df.col("URL").get(i).toString() : null;

            // LinkedIn CSV has some empty/null rows
            if (firstName == null || lastName == null || uri == null) {
                continue;
            }

            UserDto userDto = new UserDto();
            userDto.setFirstName(firstName);
            userDto.setLastName(lastName);
            if (email != null) {
                userDto.setEmail(email);
            }
            userDto.setUri(uri);

            // Both User and Profile are created through create()
            User user = create(userDto);

            // Create a connection between the current User and their uploaded connection(s)
            Connection connection = new Connection();
            connection.setUser(curUser);
            connection.setConnectionUser(user);
            connection.setDegree(1);
            connectionRepository.save(connection);
        }

        return curUser.getConnections();
    }

    private InputStream removeFirstThreeLines(MultipartFile file) throws IOException {
        BufferedReader reader = new BufferedReader(new InputStreamReader(file.getInputStream()));

        // Skip the first 3 lines and collect the rest
        String csvContent = reader.lines().skip(3).collect(Collectors.joining("\n"));

        // Create a new InputStream from the modified content
        InputStream inputStream = new ByteArrayInputStream(csvContent.getBytes());
        return inputStream;
    }
}
