// Copyright 2014-2015, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package sqs

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mikedewar/aws4"
)

var (
	// Ref:
	// http://docs.aws.amazon.com/general/latest/gr/rande.html#sqs_region
	EndPoints = map[string]string{
		"ap-northeast-1": "https://sqs.ap-northeast-1.amazonaws.com",
		"ap-southeast-1": "https://sqs.ap-southeast-1.amazonaws.com",
		"ap-southeast-2": "https://sqs.ap-southeast-2.amazonaws.com",
		"eu-west-1":      "https://sqs.eu-west-1.amazonaws.com",
		"sa-east-1":      "https://sqs.sa-east-1.amazonaws.com",
		"us-east-1":      "https://sqs.us-east-1.amazonaws.com",
		"us-west-1":      "https://sqs.us-west-1.amazonaws.com",
		"us-west-2":      "https://sqs.us-west-2.amazonaws.com",
	}

	HTTPTimeout = 25 * time.Second
)

type Client struct {
	Aws4Client  *aws4.Client
	EndPointURL string
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, HTTPTimeout)
}

func NewClient(AccessKey, SecretKey, RegionCode string) (*Client, error) {
	keys := &aws4.Keys{
		AccessKey: AccessKey,
		SecretKey: SecretKey,
	}

	EndPointURL := EndPoints[RegionCode]
	if EndPointURL == "" {
		return nil, errors.New("Unknown region: " + RegionCode)
	}

	transport := &http.Transport{Dial: dialTimeout, ResponseHeaderTimeout: HTTPTimeout}
	client := &http.Client{Transport: transport}

	return &Client{
		Aws4Client:  &aws4.Client{Keys: keys, Client: client},
		EndPointURL: EndPointURL,
	}, nil
}

// Simple wrapper around the aws4 client Post() but less verbose.
func (client *Client) Post(queueURL, data string) (*http.Response, error) {
	return client.Aws4Client.Post(queueURL, SQSContentType,
		strings.NewReader(data))
}

// Simple wrapper around the aws4 Get() to keep it consistent.
func (client *Client) Get(url string) (*http.Response, error) {
	return client.Aws4Client.Get(url)
}

// Return a single message body, with its ReceiptHandle. A lack of message is
// not considered an error but *Message will be nil.
func (client *Client) GetMessagesFromRequest(request *Request) ([]*Message, error) {
	var m ReceiveMessageResult
	var messages []*Message

	// These two settings are required for this function to function.
	request.Set("AttributeName", "SenderId")

	resp, err := client.Get(request.URL())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	err = xml.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(m.Bodies); i++ {
		msg := &Message{
			QueueURL:      request.QueueURL,
			Body:          m.Bodies[i],
			ReceiptHandle: m.ReceiptHandles[i],
			UserID:        m.Values[i],
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

// Return a single message with its ReceiptHandle. A lack of message is not
// considered an error but the return message will be nil.
func (client *Client) GetSingleMessage(url string) (*Message, error) {
	request := NewReceiveMessageRequest(url)
	request.Set("MaxNumberOfMessages", "1")
	messages, err := client.GetMessagesFromRequest(request)
	if err != nil {
		return nil, err
	}

	return messages[0], nil
}

// Return queue messages, with its ReceiptHandle. A lack of message is
// not considered an error but the return message will be nil.
func (client *Client) GetMessages(url string) ([]*Message, error) {
	request := NewReceiveMessageRequest(url)
	request.Set("MaxNumberOfMessages", "10")
	return client.GetMessagesFromRequest(request)
}

// Conduct a DeleteMessage API call on the given queue, using the receipt
// handle from a previously fetched message.
func (client *Client) DeleteMessageFromReceipt(queueURL, receipt string) error {
	url := NewDeleteMessageRequest(queueURL, receipt).URL()

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Conduct a DeleteMessage API call on the given queue, using the receipt
// handle from a previously fetched message.
func (client *Client) DeleteMessage(msg *Message) error {
	url := NewDeleteMessageRequest(msg.QueueURL, msg.ReceiptHandle).URL()

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Conduct a SendMessage API call (POST) on the given queue.
func (client *Client) SendMessage(queueURL, message string) error {
	data := NewSendMessageRequest(queueURL, message).Query()

	resp, err := client.Post(queueURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}

	return nil
}

// Get the queue URL from its name.
func (client *Client) GetQueueURL(name string) (string, error) {
	var parsedResponse GetQueueURLResult
	url := NewGetQueueURLRequest(client.EndPointURL, name).URL()

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}

	err = xml.Unmarshal(body, &parsedResponse)
	if err != nil {
		return "", err
	}

	return parsedResponse.QueueURL, nil
}

// Create a queue using the provided attributes and return its URL. This
// function can be used to obtain the QueueURL for a queue even if it already
// exists.
func (client *Client) CreateQueueWithAttributes(name string, attributes CreateQueueAttributes) (string, error) {
	var parsedResponse CreateQueueResult
	url := buildCreateQueueURL(client.EndPointURL, name, attributes)

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}

	err = xml.Unmarshal(body, &parsedResponse)
	if err != nil {
		return "", err
	}

	return parsedResponse.QueueURL, nil
}

// Create a queue with default parameters and return its URL. This function can
// be used to obtain the QueueURL for a queue even if it already exists.
func (client *Client) CreateQueue(name string) (string, error) {
	url, err := client.CreateQueueWithAttributes(name, CreateQueueAttributes{})
	return url, err
}
