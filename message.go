// Copyright 2014-2015, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqs

import (
	"crypto/md5"
	"fmt"
	"io"
)

type Message struct {
	// This is not part of the data received from the SQS servers, it is
	// assigned internally to keep track of where this message come from.
	QueueURL string

	MessageID                        string
	MD5                              string
	Body                             string
	ReceiptHandle                    string
	SenderID                         string
	SentTimestamp                    uint64
	ApproximateReceiveCount          uint64
	ApproximateFirstReceiveTimestamp uint64
}

func (msg Message) CheckMD5() bool {
	h := md5.New()
	io.WriteString(h, msg.Body)
	md5 := fmt.Sprintf("%x", h.Sum(nil))
	if md5 == msg.MD5 {
		return true
	}
	return false
}
