// Copyright 2014-2015, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqschan

import (
	"github.com/truveris/sqs"
	"sort"
)

// BySentTimestamp is used to sort SQS Messages by their timestamp if multipled
// are fetched at the same time.
type BySentTimestamp []*sqs.Message

func (a BySentTimestamp) Len() int {
	return len(a)
}
func (a BySentTimestamp) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a BySentTimestamp) Less(i, j int) bool {
	return a[i].SentTimestamp < a[j].SentTimestamp
}

// IncomingFromURL creates a receiving channel from the provided queue URL.  An
// error channel is created as well in case any error is encountered.
func IncomingFromURL(client *sqs.Client, url string) (<-chan *sqs.Message, <-chan error, error) {
	ch := make(chan *sqs.Message)
	errch := make(chan error)

	go func() {
		for {
			req := sqs.NewReceiveMessageRequest(url)
			req.Set("WaitTimeSeconds", "20")
			req.Set("MaxNumberOfMessages", "10")
			req.Set("AttributeName", "All")

			msgs, err := client.GetMessagesFromRequest(req)
			if err != nil {
				errch <- err
				continue
			}

			sort.Sort(BySentTimestamp(msgs))

			for _, msg := range msgs {
				ch <- msg
			}
		}
	}()

	return ch, errch, nil
}

// Incoming creates a receiving channel from the provided queue name as well as
// an error channel in case any error occurs receiving messages.
func Incoming(client *sqs.Client, name string) (<-chan *sqs.Message, <-chan error, error) {
	url, err := client.CreateQueue(name)
	if err != nil {
		return nil, nil, err
	}

	return IncomingFromURL(client, url)
}

// OutgoingFromURL creates a channel from a queue URL that will forward all the
// data to the provided queue.  An error channel is also provided in case any
// error occurs during the transmission.
func OutgoingFromURL(client *sqs.Client, url string) (chan string, chan error, error) {
	ch := make(chan string)
	errch := make(chan error)

	go func() {
		for line := range ch {
			err := client.SendMessage(url, line)
			if err != nil {
				errch <- err
				continue
			}
		}
	}()

	return ch, errch, nil
}

// Outgoing creates a channel from a queue name that will forward all the data
// to the provided queue.  An error channel is also provided in case any error
// occurs during the transmission.
func Outgoing(client *sqs.Client, name string) (chan string, chan error, error) {
	url, err := client.CreateQueue(name)
	if err != nil {
		return nil, nil, err
	}

	return OutgoingFromURL(client, url)
}
