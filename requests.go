// Copyright 2014, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqs

import (
	"net/url"
)

const (
	SQSAPIVersion       = "2012-11-05"
	SQSSignatureVersion = "4"
	SQSContentType      = "application/x-www-form-urlencoded"
)

// Build the data portion of a Message POST API call.
func BuildSendMessageData(msg string) string {
	query := url.Values{}
	query.Set("Action", "SendMessage")
	query.Set("Version", SQSAPIVersion)
	query.Set("SignatureVersion", SQSSignatureVersion)
	query.Set("MessageBody", msg)
	return query.Encode()
}

// Build the URL to conduct a ReceiveMessage GET API call.
func BuildReceiveMessageURL(queueURL string) string {
	query := url.Values{}
	query.Set("Action", "ReceiveMessage")
	query.Set("AttributeName", "SenderId")
	query.Set("Version", SQSAPIVersion)
	query.Set("SignatureVersion", SQSSignatureVersion)
	query.Set("WaitTimeSeconds", "20")
	query.Set("MaxNumberOfMessages", "1")
	url := queueURL + "?" + query.Encode()
	return url
}

// Build the URL to conduct a DeleteMessage GET API call.
func BuildDeleteMessageURL(queueURL, receipt string) string {
	query := url.Values{}
	query.Set("Action", "DeleteMessage")
	query.Set("ReceiptHandle", receipt)
	query.Set("Version", SQSAPIVersion)
	query.Set("SignatureVersion", SQSSignatureVersion)
	url := queueURL + "?" + query.Encode()
	return url
}
