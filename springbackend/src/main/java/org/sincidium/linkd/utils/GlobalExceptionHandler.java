package org.sincidium.linkd.utils;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.server.ResponseStatusException;

// Declare as global exception handler using this annotation
@ControllerAdvice
public class GlobalExceptionHandler {

    Logger logger = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    @Value("${application.environment}")
    private String environment;

    // Global exception handler for ResourceNotFoundException
    @ExceptionHandler(ResourceNotFoundException.class)
    public ResponseEntity<String> handleResourceNotFoundException(ResourceNotFoundException ex) {
        String errorMessage = ex.getMessage();
        logger.error(errorMessage);
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(errorMessage);
    }

    // Exception handler for re-throwing
    @ExceptionHandler(ResponseStatusException.class)
    public ResponseEntity<String> handleResponseStatusException(ResponseStatusException ex) {
        String errorMessage = ex.getMessage();
        logger.error(errorMessage);
        return ResponseEntity.status(ex.getStatusCode()).body(errorMessage);
    }

    // Global exception handler for all uncaught runtime exceptions
    @ExceptionHandler(RuntimeException.class)
    public ResponseEntity<String> handleUncaughtExceptions(RuntimeException ex) {
        String errorMessage = ex.getMessage();
        logger.error(errorMessage);
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(errorMessage);
    }
}
