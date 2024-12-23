# Payment Gateway System

This project is a robust and scalable **Payment Gateway System** designed to handle multi-gateway configurations, region-based gateway selection, and efficient payment processing. The system integrates multiple payment gateways, ensuring flexibility and reliability for international transactions.

---

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [How It Meets the Assessment Requirements](#how-it-meets-the-assessment-requirements)
4. [Architecture](#architecture)
5. [Assumptions](#assumptions)
7. [How It Works](#how-it-works)
8. [Region-Based Gateway Selection](#region-based-gateway-selection)
9. [Gateway Configurations](#gateway-configurations)
10. [Setup Instructions](#setup-instructions)
10. [API Documentation](#api-documentation)
11. [Testing](#testing)
12. [Future Improvements](#future-improvements)

---

## Introduction

The Payment Gateway System is designed to simplify and optimize payment processing for applications requiring seamless integration with multiple payment providers. It incorporates fault-tolerant designs, a database-driven health and priority management system, and a flexible architecture.

---

## Features

- **Dynamically Select Payment Gateways**: Automatically chooses the appropriate gateway based on the user’s country or region.
- **Gateway Priority and Fallbacks**: Supports gateway priority configurations and automatic fallback mechanisms in case of failures.
- **Compliance with Regional Regulations**: Handles transaction data securely with encryption and ensures secure storage to comply with regulations.
- **Asynchronous Callback Handling**: Manages gateway callbacks asynchronously and updates transaction statuses.
- **Health Monitoring with Background Cron**: Periodically checks the health status of gateways and updates the database to ensure accurate gateway availability.
- **Resilience with a Circuit-Breaker-Like Implementation**: Implements fallback mechanisms and fault tolerance similar to a circuit breaker to ensure system stability and resilience.

---

## How It Meets the Assessment Requirements

This project directly addresses the required features as outlined in the assessment file:

1. **Dynamically Select a Payment Gateway**:
   - Based on the user's country or region, the system selects the appropriate payment gateway from a list of supported gateways.

2. **Configure Gateway Priority and Fallbacks**:
   - The system allows configuration of gateway priority.
   - Automatic fallback to alternative gateways is implemented in case of failures.

3. **Ensure Compliance with Regional Regulations**:
   - Transaction data is encrypted and stored securely to meet regulatory requirements.

4. **Manage Callbacks Asynchronously**:
   - The system handles callbacks from gateways asynchronously and updates transaction statuses accordingly.

5. **Implement Resilience and Scalability**:
   - The implementation ensures resilience and fault tolerance by managing gateway health and dynamically selecting the next available gateway in case of failures.

---

## Architecture

The system follows a modular architecture to ensure flexibility and maintainability. The key components and their interactions are illustrated in the architecture diagram below:

https://drive.google.com/file/d/1VC8E66q2ZbjyA9r0VCGRD7ILS1VBu_ar/view?usp=sharing

- **User**: Initiates requests through the `/deposit` endpoint.
- **Application Layer**: Processes requests and communicates with the database and message broker.
- **Database**: Stores transaction information and gateway statuses.
- **Message Broker (e.g., Kafka)**: Handles asynchronous transaction processing.
- **Consumer**: Consumes messages from the queue and processes transactions.
- **Third-Party**: Processes the transactions sent from the application.
- **Cron Job**: Periodically monitors and updates the health status of gateways in the database.
- **Callback Listener**: Handles callbacks from third-party gateways to update transaction statuses.
SUCCESS
---

## Assumptions

1. **Gateway Health Check**:
   - The health of a gateway is determined based on periodic checks (via cron jobs) or transaction failures.
   - A gateway marked as `unhealthy` will not be considered for further transactions until it's explicitly recovered.

2. **Third-Party Gateway Integration**:
   - Third-party gateways are expected to return consistent and well-documented responses.
   - Callback mechanisms provided by the gateways are assumed to be reliable and secure.

3. **Message Broker**:
   - The message broker (e.g., Kafka) is assumed to be reliable and scalable to handle the volume of transactions.

4. **Database Consistency**:
   - The database is assumed to be highly available and supports transactional consistency for critical operations.

5. **Security**:
   - Sensitive data (e.g., user details, transaction data) is encrypted in transit and at rest.
   - Secrets like API keys and encryption keys are securely managed (though the implementation of secrets management is pending).

6. **Scalability**:
   - The current architecture assumes the system can scale horizontally to handle high traffic by adding more consumers and processing nodes.

7. **Cron Job Frequency**:
   - The cron job frequency for health checks is assumed to be sufficient to ensure real-time or near-real-time updates for gateway statuses.

8. **Single Repository for Components**:
   - All components (REST API, producer, gateway selector, consumer, third-party service integration, and cron job for gateway status) are currently built within a single repository. This design is temporary and focused on ensuring functionality for the project.


## How It Works

1. **User Request**:
   - The user initiates a deposit request by sending a `POST` request to the `/deposit` endpoint.

2. **Store Transaction**:
   - The system stores the transaction in the database with a status of `PENDING`.

3. **Produce the Transaction**:
   - The transaction is published to a Kafka topic (or a similar message broker) for asynchronous processing.

4. **Consumer Process**:
   - A Kafka consumer (or similar) consumes the transaction from the queue.

5. **Gateway Selection**:
   - The system retrieves the most prioritized and healthy gateway for the transaction based on the `countryID` associated with the user.

6. **Send to Third Party**:
   - The transaction is sent to the third-party payment gateway for processing.

7. **Handle Gateway Failures**:
   - If the transaction fails due to an unhealthy gateway:
     - The system marks the gateway as `unhealthy` in the database.
     - The transaction is republished to the queue to allow the next prioritized gateway to handle it.

8. **Successful Transaction**:
   - Once the transaction is successfully processed by a gateway:
     - The system updates the `gateway_id` in the transaction table to reflect the successful gateway.
     - The system waits for the callback from the third-party gateway.

9. **Callback Handling**:
   - Upon receiving the callback from the external gateway:
     - The system updates the status of the corresponding transaction (e.g., `COMPLETED`, `FAILED`, etc.) in the database.

10. **Health Monitoring (Cron Job)**:
    - A background cron job periodically runs to check the health status of all gateways.
    - Gateways that fail the health check are marked as `unhealthy`.
    - The system ensures that only `healthy` gateways are considered for transactions, maintaining high availability.

---

## Region-Based Gateway Selection

The system uses the following logic for selecting the appropriate gateway:

1. Retrieve gateways available for the user’s region.
2. Check the health status of each gateway.
3. Sort gateways by priority (if applicable).
4. Select the first gateway with a `healthy` status.

This logic ensures that the user’s transactions are always routed through the most prioritized and available gateway for their specific region, maximizing efficiency and reliability.

---

## Gateway Configurations

The `gateways` table is designed to store information about all available gateways. Each gateway has the following attributes:

- **name**: The unique name of the gateway (e.g., `Stripe`).
- **data_format_supported**: The data format supported by the gateway (e.g., `JSON` or `XML`).
- **health_status**: Tracks the health of the gateway (`healthy` or `unhealthy`).
- **last_checked_at**: The timestamp of the last health check.

To add a new gateway, insert a new record into the `gateways` table:
```sql
INSERT INTO gateways (name, data_format_supported) 
VALUES ('Stripe', 'JSON');
```

---

## Setup Instructions

### Prerequisites

- **Docker** and **Docker Compose** installed on your machine.

### Installation Steps

1. **Install Docker**:
   Ensure Docker and Docker Compose are installed on your machine. You can download Docker from [Docker's official website](https://www.docker.com/).

2. **Clone the Repository**:
   Clone the project repository to your local machine:
   ```bash
   git clone <repository-url>
   ```
   
3. **Navigate to the Project Directory**:
   Move into the cloned project directory:
   ```bash
   cd payment-gateways
   ```

4. **Run the Application**:
   Build and start the application using Docker Compose:
   ```bash
   docker-compose up --build -d
   ```

---

## API Documentation

The API documentation for this project is available in the `docs/swagger.yaml` file. Use any OpenAPI-compatible tool (e.g., Swagger UI, Postman) to view and test the API endpoints.

To view the API documentation locally:
1. Install an OpenAPI viewer such as Swagger UI or can check with https://editor.swagger.io.
2. Load the `docs/swagger.yaml` file to explore and test the available endpoints.

---

## Testing

### Unit Tests

Unit tests ensure the correctness of all critical modules in the project. They use mocking for dependencies such as repositories, Kafka producers, and other services to isolate the functionality under test.

#### Running Unit Tests
To run all unit tests across the project:
```bash
go test -v ./...
```

---

## End-to-End Testing

Simulate real-world payment processing scenarios to validate the system's reliability:

1. **Test Multiple Gateways and Regions**:
   - Configure multiple gateways with different regions in the database.
   - Ensure the system selects the appropriate gateway dynamically based on the user's region.

2. **Simulate Gateway Failures**:
   - Mark a gateway as `unhealthy` in the database.
   - Observe that the system automatically falls back to the next prioritized and healthy gateway.

3. **Test Asynchronous Callback Handling**:
   - Simulate callbacks from third-party gateways with different transaction outcomes.
   - Ensure the system updates the transaction statuses in the database correctly (e.g., `COMPLETED`, `FAILED`).

4. **Validate Health Monitoring**:
   - Test the cron job that monitors gateway health.
   - Ensure it updates the health status of gateways in the database accurately and timely.

These tests confirm the resilience and scalability of the system under various operational scenarios.

---

## Future Improvements

1. **Enhanced Monitoring and Logging**:
   - Integrate tools like **Prometheus** and **Grafana** for real-time monitoring and alerts.

2. **Support for Additional Payment Gateways**:
   - Add configurations for more payment gateways and allow dynamic region-to-gateway mapping.

3. **Load Testing**:
   - Perform detailed load testing to assess and improve system scalability under high transaction volumes.

4. **Secrets Management**:
   - Use a dedicated secrets management tool like **AWS Chamber** to securely manage sensitive data (e.g., API keys, encryption keys).

5. **Dynamic Gateway Prioritization**:
   - Implement an algorithm to dynamically adjust gateway priority based on metrics like success rate or latency.

These improvements aim to enhance the robustness, scalability, and security of the system for long-term reliability.
