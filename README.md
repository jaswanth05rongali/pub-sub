# pub-sub
A publisher-subscriber system using **Apache Kafka**, **Docker** and **GO**.

Pub-Sub is a concept where you want to broadcast the data to N number of consumers at the
same time. Consumers can use these data to facilitate their needs.
For Example, In a Payment Gateway System once a transaction is completed you might want to
send notifications via Email/SMS to the customer that the transaction is completed.
So for the above use case, we have 2 different consumers. One for consuming email messages and sending emails and the other for consuming phone messages and sending SMS.

This consists of a Producer API, two consumer services along with their respective clients(i.e. Email and Phone client). Here, we are not actually sending emails and SMS, we are just mocking the clients.

## Setting up Kafka Cluster
Here, we start the kafka cluster using the docker image `docker-compose.yml`.
In the docker file, change the IP address(which is `${MY_IP}` here) in kafka brokers configuration to your local system's IP where you want to run these services.
Also, If you want, you can have an environmental variable `MY_IP` set to your IP address on your local system. Then it is not required to change `${MY_IP}` in `docker-compose.yml`.

*Run the below commands with root i.e. use `sudo`.
For starting the Kafka Cluster -
You can run your docker using `docker-compose up -d`

You can stop docker using `docker-compose stop`.

You can create a new kafka topic using `docker run --net=host --rm confluentinc/cp-kafka:3.0.1 kafka-topics --create --topic foo --partitions 4 --replication-factor 2 --if-not-exists --zookeeper localhost:32181`.
Follow other docker commands to handle kafka topics.

## Setting up paths in the code
In `./config/config.go` setup the path variable `jsonPath` so that it recognizes json config file at `./config/config.json`. You can use relative path which is considers the path from your `go/src` folder or you can use absolute path i.e. from home directory directly to maklke it easier.
Also, in `config.json`, set your broker urls, topic names, etc.

For logging, you need to set the path variable `jsonPath` in `./logger/loggerconfig/loggerConfig.go` so that it recognizes json loggerConfig file (similar to the above case).

After installing all the required go packages in the code, You can follow the below steps to run the producer and consumer services.

## Producer
You can run your producer API by the following steps - 
1. cd to pub-sub folder on your terminal
2. Use `go run ./api/main.go` to run your producer API

## Consumer
You can run your consumer service by the following steps - 
1. cd to consumer in pub-sub folder on your terminal.
2. In consumer folder, Use `go run consumer.go --topic topic_name` to subscribe to that particular topic given by topic_name.

## Use POSTMAN
For hitting the producer API, you can use postman software and in the address bar, change type to POST and use the URL `0.0.0.0:9000/api/v1/data` to call the API.
For the input, in `Body` option change it to `raw` type and choose `JSON (application/json)` to give input in JSON format.

The input should be like this - 
```
{
	"request_id":"1",
	"topic_name":"foo",
	"message_body":"Transaction Successful",
	"transaction_id":"987456321",
	"email":"kafka@gopostman.com",
	"phone":"9876543210",
	"customer_id":"1",
	"key":"123",
	"pubMessageType":"0",
	"pubPartition":"1"
}
```
