# Simple Go SQS wrappers

Simple SQS helper module wrapping github.com/mikedewar/aws4. This module is a
thin layer it does not abstract anything from the SQS API.

## Example: Sending a message
```go
client, err := sqs.NewClient("JimmysAccessKeyId", "JimmysSecretAccessKey", "us-east-1")
if err != nil {
	// wrong region?
}

queueURL, err := client.CreateQueue("jimmys-queue")
if err != nil {
	// can't create or get the queue
}

err = client.SendMessage(queueURL, "Jimmy's message")
if err != nil {
	// unable to sent the message
}
```

## References
 * Amazon SQS API Documentation
   http://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/Welcome.html
