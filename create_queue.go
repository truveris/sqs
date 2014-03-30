// Copyright 2014, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.
//
// See the API documentation for further information on CreateQueue and its
// attributes.
//
// http://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_CreateQueue.html
//

package sqs

import (
	"fmt"
	"net/url"
)

type CreateQueueAttributes struct {
	MaximumMessageSize            int
	ReceiveMessageWaitTimeSeconds int
	VisibilityTimeout             int
	MessageRetentionPeriod        int
	Policy                        string
}

func (attrs CreateQueueAttributes) Map() map[string]string {
	m := make(map[string]string)

	if attrs.MaximumMessageSize > 0 {
		m["MaximumMessageSize"] = fmt.Sprintf("%d", attrs.MaximumMessageSize)
	}

	if attrs.ReceiveMessageWaitTimeSeconds > 0 {
		m["ReceiveMessageWaitTimeSeconds"] = fmt.Sprintf("%d",
			attrs.ReceiveMessageWaitTimeSeconds)
	}

	if attrs.VisibilityTimeout > 0 {
		m["VisibilityTimeout"] = fmt.Sprintf("%d", attrs.VisibilityTimeout)
	}

	if attrs.MessageRetentionPeriod > 0 {
		m["MessageRetentionPeriod"] = fmt.Sprintf("%d", attrs.MessageRetentionPeriod)
	}

	return m
}

func (attrs CreateQueueAttributes) AddToQuery(query *url.Values) {
	index := 1
	for name, value := range attrs.Map() {
		query.Set(fmt.Sprintf("Attribute.%d.Name", index), name)
		query.Set(fmt.Sprintf("Attribute.%d.Value", index), value)
		index++
	}
}

func buildCreateQueueURL(endpointURL, name string, attributes CreateQueueAttributes) string {
	query := url.Values{}
	query.Set("Action", "CreateQueue")
	query.Set("QueueName", name)
	attributes.AddToQuery(&query)
	query.Set("Version", SQSAPIVersion)
	query.Set("SignatureVersion", SQSSignatureVersion)
	url := endpointURL + "?" + query.Encode()
	return url
}
