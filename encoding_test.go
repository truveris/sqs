// Copyright 2014, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqs

import (
	"testing"
)

func TestSQSEncode(t *testing.T) {
	s := SQSEncode("PRIVMSG ##truveris :\x01ACTION doesn't like Java.\x01")
	if s != "PRIVMSG ##truveris :\\x01ACTION doesn't like Java.\\x01" {
		t.Fail()
	}
}

func TestSQSDecode(t *testing.T) {
	s := SQSDecode("PRIVMSG ##truveris :\\x01ACTION doesn't like Java.\\x01")
	if s != "PRIVMSG ##truveris :\x01ACTION doesn't like Java.\x01" {
		t.Fail()
	}
}
