// Copyright 2014, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqs

import (
	neturl "net/url"
	"strconv"
)

const (
	SQSAPIVersion       = "2012-11-05"
	SQSSignatureVersion = "4"
	SQSContentType      = "application/x-www-form-urlencoded"
)

type SQSRequest struct {
	Values   *neturl.Values
	QueueURL string
}

func (req *SQSRequest) Set(key, value string) {
	req.Values.Set(key, value)
}

func (req *SQSRequest) Query() string {
	return req.Values.Encode()
}

func (req *SQSRequest) URL() string {
	return req.QueueURL + "?" + req.Query()
}

// Create a basic req for all SQS requests.
func NewRequest(url, action string) *SQSRequest {
	req := &SQSRequest{
		Values:   &neturl.Values{},
		QueueURL: url,
	}
	req.Set("Version", SQSAPIVersion)
	req.Set("SignatureVersion", SQSSignatureVersion)
	req.Set("Action", action)
	return req
}

// Build the data portion of a Message POST API call.
func NewSendMessageRequest(url, body string) *SQSRequest {
	req := NewRequest(url, "SendMessage")
	req.Set("MessageBody", body)
	return req
}

// Build the URL to conduct a ReceiveMessage GET API call.
func NewReceiveMessageRequest(url string) *SQSRequest {
	req := NewRequest(url, "ReceiveMessage")
	return req
}

// Build the URL to conduct a ReceiveMessage GET API call.
func NewLongPollingReceiveSingleMessageRequest(url string, waitTimeSeconds int64) *SQSRequest {
	req := NewRequest(url, "ReceiveMessage")
	req.Set("AttributeName", "SenderId")
	req.Set("WaitTimeSeconds", strconv.FormatInt(waitTimeSeconds, 10))
	req.Set("MaxNumberOfMessages", "1")
	return req
}

// Build the URL to conduct a DeleteMessage GET API call.
func NewDeleteMessageRequest(url, receipt string) *SQSRequest {
	req := NewRequest(url, "DeleteMessage")
	req.Set("ReceiptHandle", receipt)
	return req
}
