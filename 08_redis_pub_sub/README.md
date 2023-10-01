# Network Telemetry with Redis Pub/Sub

This project demonstrates a simple network telemetry system using Redis as a Pub/Sub mechanism. The code within `main.go` contains both a publisher and a subscriber connected to Redis. We simulate network telemetry events by sending data to the publisher, which then determines if an alert should be triggered based on certain conditions. If the conditions are met, the subscriber will log and alert accordingly.

## Running the Project

Before running the main application, ensure you have Redis running. You can easily spin up a Redis instance using Docker:

```bash
docker run --name some-redis -d -p 6379:6379 redis
```

In a seperate terminal window we need to do `go run main.go` to initialize our pub/sub app

```
INFO[0000] Server starting...                           

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.11.1
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
[watermill] 2023/10/01 13:33:57.399664 main.go:86:      level=INFO  msg="Starting subscriber for topic" topic=packet-counter-errors 
[watermill] 2023/10/01 13:33:57.399733 main.go:86:      level=INFO  msg="Starting subscriber for topic" topic=packet-counter-errors 
[watermill] 2023/10/01 13:33:57.399779 subscriber.go:168:       level=INFO  msg="Subscribing to redis stream topic" consumer_group=packet-error-logs consumer_uuid=2AWoaHQFk4farNyUThPXM provider=redis topic=packet-counter-errors 
[watermill] 2023/10/01 13:33:57.399791 subscriber.go:189:       level=INFO  msg="Starting consuming" consumer_group=packet-error-logs consumer_uuid=2AWoaHQFk4farNyUThPXM provider=redis topic=packet-counter-errors 
⇨ http server started on [::]:8080
```


Next open up a terminal window and send the following CURL to the gateway on `:8080`
This is some JSON body which has some logic inside the publisher based upon the thresehold for the Packet count. If the thresehold is breached then we will send to Redis and the subscriber will pick up the message and alert and log it.

```
curl -X POST \
  http://localhost:8080/telemetry-data \
  -H 'Content-Type: application/json' \
  -d '{
    "hostname": "dc-router-1",
    "interface": "GigabitEthernet0/1/0",
    "input_errors": 1533300
}'
```
All being well then you should see the alert and log messages being generated
```
[watermill] 2023/10/01 13:45:57.088633 main.go:110:     level=INFO  msg="Received message for topic" hostname=dc-router-1 input_errors=1533300 interface=GigabitEthernet0/1/0 topic=packet-counter-errors 
WARN[0004] High Input errors detected on router: dc-router-1 on interface GigabitEthernet0/1/0 with 1533300 errors 
[watermill] 2023/10/01 13:45:57.088726 main.go:110:     level=INFO  msg="Received message for topic" hostname=dc-router-1 input_errors=1533300 interface=GigabitEthernet0/1/0 topic=packet-counter-errors 
ERRO[0004] ALERT! High Input errors on router: dc-router-1 on interface GigabitEthernet0/1/0 with 1533300 errors
```
