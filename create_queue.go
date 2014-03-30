// Copyright 2014, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.
//
// See the API documentation for further information on CreateQueue and its
// attributes.

// http://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_CreateQueue.html
//

package sqs

import (
	"fmt"
	"net/url"
)

type CreateQueueAttributes struct {
	index int
	MaximumMessageSize            int
	ReceiveMessageWaitTimeSeconds int
	VisibilityTimeout             int
	MessageRetentionPeriod        int
	Policy			      string
}

func (attrs CreateQueueAttributes) addAttributesToQuery(query *url.Values) {
	index := 1

	// 4096
	if attrs.MaximumMessageSize > 0 {
		attr := fmt.Sprintf("Attribute.%d", index)
		query.Set(attr+".Name", "MaximumMessageSize")
		query.Set(attr+".Value", fmt.Sprintf("%d", attrs.MaximumMessageSize))
		index++
	}

	// 20
	if attrs.ReceiveMessageWaitTimeSeconds > 0 {
		attr := fmt.Sprintf("Attribute.%d", index)
		query.Set(attr+".Name", "ReceiveMessageWaitTimeSeconds")
		query.Set(attr+".Value", fmt.Sprintf("%d", attrs.ReceiveMessageWaitTimeSeconds))
		index++
	}

	// 10
	if attrs.VisibilityTimeout > 0 {
		attr := fmt.Sprintf("Attribute.%d", index)
		query.Set(attr+".Name", "VisibilityTimeout")
		query.Set(attr+".Value", fmt.Sprintf("%d", attrs.VisibilityTimeout))
		index++
	}

	// 300
	if attrs.MessageRetentionPeriod > 0 {
		attr := fmt.Sprintf("Attribute.%d", index)
		query.Set(attr+".Name", "MessageRetentionPeriod")
		query.Set(attr+".Value", fmt.Sprintf("%d", attrs.MessageRetentionPeriod))
		index++
	}
}

// Build the URL to conduct a CreateMessage GET API call.
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
