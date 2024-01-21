# golang-users-api
Golang's application that ingests initial data from a CSV, populates a queue which is consumed and then structures the data in a relational database and then cached.

## Tech
- Docker
- Golang
- Postgres
- Redis
- RabbitMQ

## Steps to run
- Setup all containers `docker-compose up --build`

## Requirements
- Core
 - [x] CSV File Ingestion and Queue: Create an application that ingests a CSV file and sends its content to a message queue in the RabbitMQ broker.
 - [ ] Message Processing and Storage: Develop an application that consumes messages from the queue, structures them, and stores the data in both a PostgreSQL database and a Redis cache.
 - [x] API Development: Extend the consumer application to also serve as an API. The API should expose data via the HTTP protocol. Ensure that the GET method includes filters for queries, allowing for customizable data retrieval.
- Security
 - [ ] Implement secure communication channels (HTTPS) for the API.
 - [ ] Ensure data encryption in transit and at rest.
 - [x] Apply validation to prevent security vulnerabilities like SQL injection or cross-site scripting.
- Scalability
 - [ ] Address scalability with mechanisms for handling high data volume or traffic.
 - [ ] Implement load balancing, caching optimizations, or database sharding.
- Error Handling and Logging
  - [ ] Develop comprehensive error-handling to manage exceptions gracefully.
  - [ ] Implement robust logging for effective troubleshooting and debugging. 
- Performance Optimization
  - [x] Optimize API response times and resource utilization using techniques like query optimization, indexing, and caching.
- Documentation
  - [x] Provide thorough documentation, including code comments, API endpoint details, data models, and overall architecture. 
- Container Orchestration
  - [ ] Deploy the application using Kubernetes.
  - [ ] Showcase container management and scaling in a production environment. 

## Documentation files
- [API](docs/api.md)
- [Comments](docs/comments.md)
- [Next Steps](docs/next-steps.md)
- [Arch](docs/arch.md)