// Copyright 2014, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqs

type Message struct {
	// This is not part of the data received from the SQS servers, it is
	// assigned internally to keep track of where this message come from.
	QueueURL      string

	Body          string
	ReceiptHandle string
	UserID        string
}
