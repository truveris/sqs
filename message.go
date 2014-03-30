// Copyright 2014, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqs

type Message struct {
	Body          string
	ReceiptHandle string
	UserID        string
}
