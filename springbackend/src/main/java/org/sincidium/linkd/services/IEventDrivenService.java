package org.sincidium.linkd.services;

import com.azure.messaging.servicebus.ServiceBusErrorContext;
import com.azure.messaging.servicebus.ServiceBusReceivedMessageContext;

public interface IEventDrivenService {
    void handleMessage(ServiceBusReceivedMessageContext messageContext);
    void handleError(ServiceBusErrorContext errorContext);
}
