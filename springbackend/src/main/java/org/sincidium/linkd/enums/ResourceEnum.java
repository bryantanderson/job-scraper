package org.sincidium.linkd.enums;

public enum ResourceEnum {
    EDUCATION("education"),
    EXPERIENCE("experience"),
    DUMMY("Dummy"),
    COMPANY("Company"),
    CONNECTION("Connection"),
    USER("User"),
    JOB("Job"),
    REFERRAL("Referral"),
    PROFILE("Profile");

    private final String description;

    private ResourceEnum(String description) {
        this.description = description;
    }

    @Override
    public String toString() {
        return description;
    }
}
