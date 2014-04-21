// Copyright 2014, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqs

// struct defining what to extract from the XML document received in response
// to the GetMessage API call.
type ReceiveMessageResult struct {
	Bodies         []string `xml:"ReceiveMessageResult>Message>Body"`
	ReceiptHandles []string `xml:"ReceiveMessageResult>Message>ReceiptHandle"`
	Values         []string `xml:"ReceiveMessageResult>Message>Attribute>Value"`
}

type CreateQueueResult struct {
	QueueURL string `xml:"CreateQueueResult>QueueUrl"`
}

type GetQueueURLResult struct {
	QueueURL string `xml:"GetQueueUrlResult>QueueUrl"`
}
