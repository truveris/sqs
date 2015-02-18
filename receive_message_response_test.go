package sqs

import (
	"testing"
)

func TestGetMessages(t *testing.T) {
	data := `<ReceiveMessageResponse>
	<Things>stuff</Things>
  <ReceiveMessageResult>
    <Message>
      <MessageId>
        5fea7756-0ea4-451a-a703-a558b933e274
      </MessageId>
      <ReceiptHandle>
        MbZj6wDWli+JvwwJaBV+3dcjk2YW2vA3+STFFljTM8tJJg6HRG6PYSasuWXPJB+Cw
        Lj1FjgXUv1uSj1gUPAWV66FU/WeR4mq2OKpEGYWbnLmpRCJVAyeMjeU5ZBdtcQ+QE
        auMZc8ZRv37sIW2iJKq3M9MFx1YvV11A2x/KSbkJ0=
      </ReceiptHandle>
      <MD5OfBody>
        fafb00f5732ab283681e124bf8747ed1
      </MD5OfBody>
      <Body>This is a test message</Body>
      <Attribute>
        <Name>SenderId</Name>
        <Value>195004372649</Value>
      </Attribute>
      <Attribute>
        <Name>SentTimestamp</Name>
        <Value>1238099229000</Value>
      </Attribute>
      <Attribute>
        <Name>ApproximateReceiveCount</Name>
        <Value>5</Value>
      </Attribute>
      <Attribute>
        <Name>ApproximateFirstReceiveTimestamp</Name>
        <Value>1250700979248</Value>
      </Attribute>
    </Message>
  </ReceiveMessageResult>
  <ResponseMetadata>
    <RequestId>
      b6633655-283d-45b4-aee4-4e84e0ae6afa
    </RequestId>
  </ResponseMetadata>
</ReceiveMessageResponse>`
	r, err := NewReceiveMessageResponse([]byte(data))
	if err != nil {
		t.Fatal(err)
		return
	}

	msgs, err := r.GetMessages()
	if err != nil {
		t.Fatal(err)
		return
	}

	if len(msgs) != 1 {
		t.Fatal("should have one message")
	}
	msg := msgs[0]

	if msg.MessageID != "5fea7756-0ea4-451a-a703-a558b933e274" {
		t.Errorf("wrong MessageID, found %v", msg.MessageID)
	}
	if msg.SenderID != "195004372649" {
		t.Errorf("wrong SenderID, found %v", msg.SenderID)
	}
}
