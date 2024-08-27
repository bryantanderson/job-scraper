package org.sincidium.linkd.services;

import java.util.UUID;
import org.sincidium.linkd.dtos.AssessmentDto;
import org.sincidium.linkd.models.Assessment;
import org.sincidium.linkd.models.Job;
import org.sincidium.linkd.models.Profile;
import org.sincidium.linkd.repositories.AssessmentRepository;
import org.sincidium.linkd.utils.IdUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import com.azure.messaging.servicebus.ServiceBusErrorContext;
import com.azure.messaging.servicebus.ServiceBusReceivedMessageContext;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

@Service
public class AssessmentService implements IEventDrivenService {

    Logger logger = LoggerFactory.getLogger(AssessmentService.class);

    @Value("${service.bus.assessment.tasks.topic}")
    private String assessmentTasksTopic;

    @Value("${service.bus.assessment.results.topic}")
    private String assessmentResultsTopic;

    @Autowired
    private AssessmentRepository assessmentRepository;

    @Autowired
    private JobService jobService;

    @Autowired
    private ProfileService profileService;

    @Autowired
    private EventService eventService;

    @Autowired
    private ObjectMapper objectMapper;

    public void create(Assessment assessment) {
        assessmentRepository.save(assessment);
    }

    public void publish(UUID userId, UUID jobId, String profile_url) {
        Job job = jobService.get(jobId);
        Profile profile = profileService.get(userId);
        AssessmentDto dto = new AssessmentDto(userId, profile, job);
        try {
            String data = objectMapper.writeValueAsString(dto);
            eventService.publish(assessmentTasksTopic, data);

        } catch (JsonProcessingException ex) {
            logger.error("Failed to dispatch assessment", ex);
        }
    }

    public void subscribe() {
        eventService.subscribe(
            assessmentResultsTopic, IdUtils.topicNameToSubscriptionName(assessmentResultsTopic), this
        );
    }

    @Override
    public void handleMessage(ServiceBusReceivedMessageContext messageContext) {
        try {
            ObjectMapper mapper = new ObjectMapper();
            Assessment assessment = mapper.readValue(messageContext.getMessage().getBody().toString(), Assessment.class);
            create(assessment);

        } catch (JsonProcessingException ex) {
            logger.error("Error parsing JSON: ", ex);

        } catch (Exception ex) {
            logger.error("Error creating Assessment: ", ex);
        }
    }

    @Override
    public void handleError(ServiceBusErrorContext errorContext) {
        logger.error(errorContext.getException().getMessage());
    }

}
