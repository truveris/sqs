// Copyright 2014-2015, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqs

type CreateQueueResult struct {
	QueueURL string `xml:"CreateQueueResult>QueueUrl"`
}

type GetQueueURLResult struct {
	QueueURL string `xml:"GetQueueUrlResult>QueueUrl"`
}
