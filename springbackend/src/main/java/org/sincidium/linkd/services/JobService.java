package org.sincidium.linkd.services;

import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

import org.sincidium.linkd.dtos.JobDto;
import org.sincidium.linkd.enums.JobLocationTypeEnum;
import org.sincidium.linkd.enums.ResourceEnum;
import org.sincidium.linkd.models.Job;
import org.sincidium.linkd.models.Responsibility;
import org.sincidium.linkd.models.Qualification;
import org.sincidium.linkd.repositories.JobRepository;
import org.sincidium.linkd.repositories.QualificationRepository;
import org.sincidium.linkd.repositories.ResponsibilityRepository;
import org.sincidium.linkd.utils.IdUtils;
import org.sincidium.linkd.utils.ResourceNotFoundException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import com.azure.messaging.servicebus.ServiceBusErrorContext;
import com.azure.messaging.servicebus.ServiceBusReceivedMessageContext;
import com.fasterxml.jackson.databind.ObjectMapper;

@Service
public class JobService implements ICrudService<Job, JobDto, UUID>, IEventDrivenService {

    Logger logger = LoggerFactory.getLogger(JobService.class);

    @Value("${service.bus.job.tasks.topic}")
    private String jobTasksTopic;

    @Value("${service.bus.job.results.topic}")
    private String jobResultsTopic;

    @Autowired
    private JobRepository jobRepository;

    @Autowired
    private ResponsibilityRepository responsibilityRepository;

    @Autowired
    private QualificationRepository qualificationRepository;

    @Autowired
    private EventService eventService;

    @Override
    @Transactional
    public Job create(JobDto dto) {
        List<Responsibility> responsibilities = new ArrayList<Responsibility>();
        List<Qualification> qualifications = new ArrayList<Qualification>();
        Job job = parseDto(dto);
        Job savedJob = jobRepository.save(job);
        // Save responsibilities
        for (String resp : dto.getResponsibilities()) {
            Responsibility newResp = new Responsibility(resp, savedJob);
            responsibilities.add(newResp);
        }
        List<Responsibility> savedResponsibilities = responsibilityRepository.saveAllAndFlush(responsibilities);
        // Save qualifications
        for (String qual : dto.getQualifications()) {
            Qualification newQual = new Qualification(qual, savedJob);
            qualifications.add(newQual);
        }
        List<Qualification> savedQualifications = qualificationRepository.saveAllAndFlush(qualifications);
        // Set the saved entities within the main job entity
        savedJob.setResponsibilities(savedResponsibilities);
        savedJob.setQualifications(savedQualifications);
        return savedJob;
    }

    public void createFromUnstructuredData(String data) {
        eventService.publish(jobTasksTopic, data);
    }

    @Override
    public Job get(UUID id) throws ResourceNotFoundException {
        return jobRepository.findById(id).orElseThrow(() -> new ResourceNotFoundException(ResourceEnum.JOB));
    }

    @Override
    public Job update(UUID id, JobDto dto) {
        Job job = get(id);
        copyDtoAttributes(job, dto);
        return jobRepository.save(job);
    }

    @Override
    public void delete(UUID id) {
        jobRepository.deleteById(id);
    }

    @Override
    public Job parseDto(JobDto dto) {
        Job job = new Job();
        copyDtoAttributes(job, dto);
        return job;
    }

    @Override
    public void copyDtoAttributes(Job model, JobDto dto) {
        model.setTitle(dto.getTitle());
        model.setDescription(dto.getDescription());
        model.setLocation(dto.getLocation());
        model.setYearsOfExperience(dto.getYearsOfExperience());
        model.setLocationType(parseJobLocationType(dto.getLocationType()));
    }

    public void subscribe() {
        eventService.subscribe(
            jobResultsTopic, IdUtils.topicNameToSubscriptionName(jobResultsTopic), this
        );
    }

    @Override
    public void handleMessage(ServiceBusReceivedMessageContext messageContext) {
        try {
            ObjectMapper mapper = new ObjectMapper();
            JobDto dto = mapper.readValue(messageContext.getMessage().getBody().toString(), JobDto.class);
            create(dto);
        } catch (Exception ex) {
            logger.error(ex.getMessage());
        }
    }

    @Override
    public void handleError(ServiceBusErrorContext errorContext) {
        logger.error(errorContext.getException().getMessage());
    }

    private JobLocationTypeEnum parseJobLocationType(String jobLocationType) throws ResponseStatusException {
        try {
            return JobLocationTypeEnum.valueOf(jobLocationType);
        } catch (Exception ex) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Invalid job location type");
        }
    }


}
