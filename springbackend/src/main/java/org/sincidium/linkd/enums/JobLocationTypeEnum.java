package org.sincidium.linkd.enums;

public enum JobLocationTypeEnum {
    onsite("onsite"),
    hybrid("hybrid"),
    remote("remote");

    private final String description;

    private JobLocationTypeEnum(String description) {
        this.description = description;
    }

    public String toString() {
        return this.description;
    }

}
