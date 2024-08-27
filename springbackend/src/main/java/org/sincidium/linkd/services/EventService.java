package org.sincidium.linkd.services;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import com.azure.messaging.servicebus.ServiceBusClientBuilder;
import com.azure.messaging.servicebus.ServiceBusMessage;
import com.azure.messaging.servicebus.ServiceBusProcessorClient;
import com.azure.messaging.servicebus.ServiceBusSenderClient;

@Service
public class EventService {

    Logger logger = LoggerFactory.getLogger(EventService.class);

    @Value("${service.bus.connection.string}")
    public String connectionString;

    public void publish(String topicName, String messageBody) {
        ServiceBusSenderClient sender = new ServiceBusClientBuilder()
                .connectionString(connectionString)
                .sender()
                .topicName(topicName)
                .buildClient();
        ServiceBusMessage message = new ServiceBusMessage(messageBody);
        try {
            sender.sendMessage(message);
        } catch (Exception ex) {
            logger.error("Failed to publish message to topic " + topicName);
        }
        sender.close();
        logger.info("Successfully published to topic " + topicName);
    }

    public void subscribe(String topicName, String subscriptionName, IEventDrivenService service) {
        ServiceBusProcessorClient processor = new ServiceBusClientBuilder()
                .connectionString(connectionString)
                .processor()
                .topicName(topicName)
                .subscriptionName(subscriptionName)
                .processMessage(service::handleMessage)
                .processError(service::handleError)
                .buildProcessorClient();
        processor.start();
    }

}
