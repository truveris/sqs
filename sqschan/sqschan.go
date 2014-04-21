// Copyright 2014, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqschan

import (
	"github.com/truveris/sqs"
)


func IncomingFromURL(client *sqs.Client, url string) (<-chan *sqs.Message, <-chan error, error) {
	ch := make(chan *sqs.Message)
	errch := make(chan error)

	go func() {
		for {
			req := sqs.NewReceiveMessageRequest(url)
			req.Set("WaitTimeSeconds", "20")

			msg, err := client.GetSingleMessageFromRequest(req)
			if err != nil {
				errch <- err
				continue
			}

			if msg == nil {
				continue
			}

			ch <- msg
		}
	}()

	return ch, errch, nil
}

func Incoming(client *sqs.Client, name string) (<-chan *sqs.Message, <-chan error, error) {
	url, err := client.CreateQueue(name)
	if err != nil {
		return nil, nil, err
	}

	return IncomingFromURL(client, url)
}

// Returns a channel that will forward all the data to the provided queue.
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

func Outgoing(client *sqs.Client, name string) (chan string, chan error, error) {
	url, err := client.CreateQueue(name)
	if err != nil {
		return nil, nil, err
	}

	return OutgoingFromURL(client, url)
}
