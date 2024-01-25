/*
Copyright 2023 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/nuclio/nuclio-sdk-go"
	amqp "github.com/rabbitmq/amqp091-go"

	"context"
	"log"
	"time"
)

const (
	rabbitMQURL = "amqp://jeff:jeff@host.docker.internal:5672/"
	queueName   = "stretch" // queue to publish the result
)

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial(rabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	q, err := declareQueue(ch)
	failOnError(err, "Failed to declare a queue")

	// Take input as string
	input := string(event.GetBody())

	// Convert input to uppercase
	revertedString := reverseString(input)

	// Publish the message to the queue
	err = publishMessage(ch, q, revertedString)
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", revertedString)

	return nuclio.Response{
		StatusCode:  200,
		ContentType: "application/text",
		Body:        []byte(revertedString),
	}, nil
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}

func publishMessage(ch *amqp.Channel, q amqp.Queue, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
