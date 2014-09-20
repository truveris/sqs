// Copyright 2014, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

// TODO: This only covers 0x01 to get it working on IRC. Clearly this should
// cover all the characters.

package sqs

import (
	"strings"
)

// Encode a string for SQS transit. Make sure you pass it through SQSDecode.
func SQSEncode(s string) string {
	s = strings.Replace(s, "\x01", "\\x01", -1)
	return s
}

func SQSDecode(s string) string {
	s = strings.Replace(s, "\\x01", "\x01", -1)
	return s
}
