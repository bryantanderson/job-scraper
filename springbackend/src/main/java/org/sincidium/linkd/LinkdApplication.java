package org.sincidium.linkd;

import org.sincidium.linkd.services.AssessmentService;
import org.sincidium.linkd.services.JobService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.domain.EntityScan;

@EntityScan(basePackages = "org.sincidium.linkd.models")
@SpringBootApplication
public class LinkdApplication implements CommandLineRunner {

	@Autowired
	private AssessmentService assessmentService;

	@Autowired
	private JobService jobService;

	public static void main(String[] args) {
		SpringApplication.run(LinkdApplication.class, args);
	}

	@Override
	public void run(String... args) throws Exception {
		jobService.subscribe();
		assessmentService.subscribe();
	}

}
