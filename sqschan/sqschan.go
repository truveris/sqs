// Copyright 2014, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqschan

import (
	"github.com/truveris/sqs"
)


func ReadBody(client sqs.Client, name string) (<-chan string, <-chan error, error) {
	ch := make(chan string)
	errch := make(chan error)

	url, err := client.CreateQueue(name)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		for {
			msg, err := client.GetMessage(url)
			if err != nil {
				errch <- err
				continue
			}

			err = client.DeleteMessage(msg.QueueURL, msg.ReceiptHandle)
			if err != nil {
				errch <- err
				continue
			}

			ch <- msg.Body
		}
	}()

	return ch, errch, nil
}

func ReadMsg(client sqs.Client, name string) (<-chan sqs.Message, <-chan error, error) {
	ch := make(chan sqs.Message)
	errch := make(chan error)

	url, err := client.CreateQueue(name)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		for {
			msg, err := client.GetMessage(url)
			if err != nil {
				errch <- err
				continue
			}

			ch <- *msg
		}
	}()

	return ch, errch, nil
}
