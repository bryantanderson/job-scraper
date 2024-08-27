package org.sincidium.linkd.utils;

import org.sincidium.linkd.enums.ResourceEnum;

public class ResourceNotFoundException extends RuntimeException {
    public ResourceNotFoundException(ResourceEnum resourceName) {
        super(resourceName.toString() + " was not found");
    }
}
