# Simple Go SQS wrappers

Simple SQS helper package with an optional channel abstraction.

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

## Example: Simplified string channel
This example returns a simple string channel and an error channel:

```go
import (
	"github.com/truveris/sqs"
	"github.com/truveris/sqs/sqschan"
)

func main() {
	client := sqs.Client(...)

	ch, errch, err := sqschan.ReadBody(client, "my-queue-name")
	if err != nil {
		// do your thing
	}

    var body string

	for {
		select {
			case err = <-errch:
				// do your error thing
			case body = <-ch:
				// do your thing, body is a string
				fmt.Printf("body: " + body)
		}
	}
}
```

Using ReadBody() does not allow you to acknowledge the reception of the message,
the message is deleted automatically after being delivered to the channel.  If
you want to allow the message to go back in the queue automatically in case of
failure, use the sqschan.Incoming().

## Example: Message channel
This example allows the use of the entire response with all its meta-data:

```go
import (
	"github.com/truveris/sqs"
	"github.com/truveris/sqschan"
)

func main() {
	client := sqs.Client(...)

	ch, errch, err := sqschan.Incoming(client, "jimmy-queue")
	if err != nil {
		// do your thing
	}

    var msg *sqs.Message

	for {
		select {
			case err = <-errch:
				// do your error thing
			case msg = <-ch:
				// do your thing.
				client.DeleteMessage(msg)
		}
	}
}
```

To the contrary of ReadBody(), you are obligated to acknowledge the reception of
the messages by deleting them. You also have access to all the meta-data
shipped with the message: msg.UserID, msg.ReceiptHandle, etc.

## References
 * Amazon SQS API Documentation
   http://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/Welcome.html
