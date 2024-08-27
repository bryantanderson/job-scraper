package org.sincidium.linkd.config;

import org.springframework.stereotype.Component;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;

@Component
public class CustomMetricsService {

    private final Counter customMetricCounter;

    /*
     * By injecting MeterRegistry, we register the custom metric and make it
     * available for scraping by Prometheus
     */
    public CustomMetricsService(MeterRegistry meterRegistry) {
        customMetricCounter = Counter.builder("custom_metric_name")
                .description("Custom metric for Linkd Spring Boot backend")
                .tags("environment", "development")
                .register(meterRegistry);
    }

    public void incrementCustomMetric() {
        customMetricCounter.increment();
    }

}
