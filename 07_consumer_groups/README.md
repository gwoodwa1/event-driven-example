# Subscribing Using Consumer Groups for a Network Telemetry use case with Redis

The code is designed to handle topics like `Packet Counter Errors` as a solid example of things we could be subscribing to from a message broker like redis.
We could have something like Prometheus scraping the data.

## Code Overview

The main application sets up subscribers to process telemetry data. The data is processed based on the topic and is either logged or alerts are sent based on the telemetry readings.

- **Subscribers**: The application sets up two subscribers for telemetry data. One subscriber logs the data, and the other sends alerts based on the data.
- **Redis as Message Broker**: The application uses Redis as a message broker to handle the telemetry data messages.
- **Consumer Groups**: The application leverages the concept of consumer groups in Redis. Consumer groups allow multiple subscribers to divide the message processing load among them. Each message is delivered to a single subscriber within the group, ensuring that messages are processed once and only once, even if multiple instances of the application are running. This ensures efficient and reliable processing of telemetry data.

## How Consumer Groups Work

In the context of Redis Streams, a consumer group is a way to ensure that multiple subscribers can read messages from a topic (or stream) without overlapping. Each message in the stream is delivered to one and only one member of the consumer group. This allows for:

1. **Load Balancing**: Messages are distributed among all subscribers in the group, spreading the processing load.
2. **Fault Tolerance**: If a subscriber fails to process a message, the message can be reprocessed by another subscriber in the group.
   
![consumer_groups](https://github.com/gwoodwa1/event-driven-example/assets/63735312/577a62c9-ee4e-4878-90f8-2e613919640a)

## Setup

To run the application, you need to have a Redis instance running. You can easily set up a Redis container using Docker with the following command:

```bash
docker run --name redis-container -p 6379:6379 -d redis
```

Additionally, set the environment variable for the Redis connection:

```bash
export REDIS_ADDR=127.0.0.1:6379
```

This ensures that the application can connect to the Redis instance for message processing.

