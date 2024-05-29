## Order Management System (Microservice)

![Architecture Design](./static/architecture.png)

### Installation
```bash
cd {service}
go mod init github.com/naufalihsan/{service}
go work init ./{service}
go install github.com/cosmtrek/air@latest
```

### FAQs
- Why using gRPC?
    ```
    - Low Latency
    - Strongly typed contracts between the services
    - Protocol buffers > JSON
    ```

- How consul works in a microservice architecture?
    ```
    - Service Registration: Each service instance registers itself with Consul when it starts up. This can be done through a Consul agent that runs alongside the service.
    - Service Discovery: When a service needs to communicate with another service, it queries Consul to find the instances of the target service. Consul responds with the addresses of healthy instances.
    - Health Checks: Consul performs health checks on registered services to ensure they are functioning correctly. If a service fails a health check, it is removed from the list of available services.
    - Configuration Management: Services can read their configuration from the Consul key/value store. This allows for centralized management of configuration data, which can be updated dynamically without restarting the services.
    - Security and Segmentation: Consul enforces access control policies to ensure that only authorized services can communicate with each other, enhancing the security of the microservice architecture.
    ```

- What's the different between direct and fanout exchange?
    | Aspect            | Direct Exchange                            | Fanout Exchange                                  |
|-------------------|--------------------------------------------|--------------------------------------------------|
| Routing Mechanism | Matches routing key exactly                | Ignores routing key, broadcasts to all queues    |
| Granularity       | Fine-grained control based on routing keys | Broad dissemination, no filtering by routing key |
| Use Case          | Targeted message delivery                  | Broadcast messages to multiple consumers         |
| Performance       | Efficient for targeted delivery            | Can cause higher load due to broad dissemination |
| Example Scenario  | Routing logs based on severity             | Broadcasting chat messages                       |
| Pattern           | CQRS, Task Distribution                    | Event-Driven Architecture, Pub/Sub               |