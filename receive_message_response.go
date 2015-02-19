package sqs

import (
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
)

var (
	errBadMD5Sum = errors.New("Body MD5 is invalid")
)

type ReceiveMessageResponse struct {
	XMLName              xml.Name `xml:"ReceiveMessageResponse"`
	Things               string
	ReceiveMessageResult struct {
		Message []struct {
			MessageId     string
			ReceiptHandle string
			MD5OfBody     string
			Body          string
			Attribute     []struct {
				Name  string
				Value string
			}
		}
	}
	ResponseMetadata struct {
		RequestId string
	}
}

func (r *ReceiveMessageResponse) GetMessages(QueueURL string) ([]*Message, error) {
	var err error
	var messages []*Message

	for _, rawmsg := range r.ReceiveMessageResult.Message {
		msg := &Message{
			MessageID:     strings.Trim(rawmsg.MessageId, " \t\n\r"),
			ReceiptHandle: strings.Trim(rawmsg.ReceiptHandle, " \t\n\r"),
			MD5:           strings.Trim(rawmsg.MD5OfBody, " \t\n\r"),
			Body:          rawmsg.Body,
			QueueURL:      QueueURL,
		}

		// That should never be a problem unless some underlying
		// implementation or dependency is doing something unhealthy.
		if !msg.CheckMD5() {
			return nil, errBadMD5Sum
		}

		for _, attr := range rawmsg.Attribute {
			err = nil
			switch attr.Name {
			case "SenderId":
				msg.SenderID = attr.Value
			case "SentTimestamp":
				msg.SentTimestamp, err = strconv.ParseUint(attr.Value, 10, 64)
			case "ApproximateReceiveCount":
				msg.ApproximateReceiveCount, err = strconv.ParseUint(attr.Value, 10, 64)
			case "ApproximateFirstReceiveTimestamp":
				msg.ApproximateFirstReceiveTimestamp, err = strconv.ParseUint(attr.Value, 10, 64)
			}

			if err != nil {
				return nil, err
			}
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func NewReceiveMessageResponse(data []byte) (*ReceiveMessageResponse, error) {
	r := &ReceiveMessageResponse{}

	err := xml.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
