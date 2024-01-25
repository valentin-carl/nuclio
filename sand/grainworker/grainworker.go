package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

// The grainworker is a separate program from Nuclio. With each function (in the paper: "grain"), a grain
// worker is deployed. The grain worker subscribes to the respective queue and invokes the function.
//
// - For later: look into how hard it would be to also run the grain worker in the function's container
//   => look at what these containers are => if it's alpine e.g., adding the grain worker is simple enough,
//      if its something else, this might be more difficult
//   => also, having the grain worker outside makes it available regardless of which runtime we're using
//
// TODO how to deploy the grain worker automatically with the function? => do it manually for now
//
// Note: so far, this only supports async calls; in the paper it appears to only work that way,
// but maybe we could use rabbitmq's request-response pattern
//
// TODO should this be included in "make build"? maybe create "make grainworker"?

type Grainworker struct {
	functionName string // get this as CLI argument
	functionIp   string // tells us where to reach the function
	amqpUrl      string // also get this as CLI argument, tells us where to reach the queue
	conn         *amqp.Connection
	ch           *amqp.Channel
	queue        amqp.Queue
	client       *http.Client // used to invoke the functions directly
}

func NewGrainworker(functionName, functionIp, amqpUrl string) *Grainworker {

	// create connection to broker
	conn, err := amqp.Dial(amqpUrl)
	handle(err)

	// the connection handles low level stuff, the channel gives access to most of the rabbitmq api
	ch, err := conn.Channel()
	handle(err)

	// create the queue if it doesn't exist yet
	queue, err := ch.QueueDeclare(
		functionName,
		false,
		false,
		false,
		false,
		nil,
	)

	// assign variables
	return &Grainworker{
		functionName: functionName,
		functionIp:   functionIp,
		amqpUrl:      amqpUrl,
		conn:         conn,
		ch:           ch,
		queue:        queue,
		client:       &http.Client{},
	}
}

func (gw *Grainworker) start() {

	defer gw.shutdown()

	log.Println("grain worker started")

	// register at queue as consumer
	events, err := gw.ch.Consume(
		gw.functionName,
		fmt.Sprintf("grain-worker-%s", gw.functionName),
		false,
		false, // exclusive=true could be a good idea
		false,
		false,
		nil,
	)
	handle(err)

	log.Println("grain worker registered as consumer at broker")

	// event loop: listen for requests and invoke the function
	for event := range events {

		log.Println("received event, invoking function ...")

		// invoke the function
		go func() {
			err := gw.invoke(&event)
			handle(err)
		}()

		// send acknowledgement to broker
		err := event.Ack(false)

		log.Println("sent acknowledgement to broker")

		handle(err)
	}

	// this shouldn't be reached
	log.Panic("events channel closed unexpectedly")
}

// Invoke calls the function this grainworker handles
func (gw *Grainworker) invoke(event *amqp.Delivery) error {

	log.Println("invoking function")

	// create the request
	requestBody := bytes.NewReader(event.Body)
	req, err := http.NewRequest(
		// TODO does this have to be/should it be a post request?
		//  => maybe add some logic with the rabbitmq headers here later on
		http.MethodPost,
		gw.functionIp,
		requestBody,
	)
	if err != nil {
		return err
	}

	log.Println("read request body", string(event.Body))

	// also send headers from rabbitmq event
	err = event.Headers.Validate()
	handle(err)
	for key, value := range event.Headers {
		req.Header.Set(key, value.(string))
	}

	log.Println("copied headers", event.Headers)

	// send it!
	response, err := gw.client.Do(req)
	if err != nil {
		return err
	}

	log.Println("sent request to function container")

	// TODO what to do with the response?
	//  in workflows, we don't really expect a direct response from the function
	//  => for now, just log it and check if it's a 200

	// read + log the response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	log.Println("received response")
	log.Println(response.Status, string(responseBody))

	if response.StatusCode != http.StatusOK {
		return errors.New("didn't get 200 back, what happened? " + response.Status)
	}

	// returns nil if invocation was successful
	return nil
}

// Shutdown closes the amqp connection + channel
func (gw *Grainworker) shutdown() {

	log.Println("shutting down ...")

	// TODO do the go stuff where you listed to ctrl c to cancel and then close conn, channel and queue etc
	//  => do that to determine when to call this function

	// close the channel
	err := gw.ch.Close()
	handle(err)

	// close the connection
	err = gw.conn.Close()
	handle(err)
}

// //////////////// //
// helper functions //
// //////////////// //

// String show gw config as json
func (gw *Grainworker) String() string {
	return fmt.Sprintf(`{
	"functionName": "%s",
	"functionIp": "%s",
	"amqpUrl": "%s"
}`, gw.functionName, gw.functionIp, gw.amqpUrl)
}

func handle(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

// /////////////////// //
// program entry point //
// /////////////////// //

type gwConfig struct {
	FunctionName string `json:"functionName"`
	FunctionUrl  string `json:"functionUrl"`
	BrokerUrl    string `json:"brokerUrl"`
}

func main() {

	var config gwConfig

	fname := os.Getenv("FUNCTION_NAME")
	if fname == "" {
		panic("my name is jeff")
	}
	fport := os.Getenv("FUNCTION_PORT")
	if fport == "" {
		panic("my name is jeff ... NOT")
	}

	config = gwConfig{
		fname,
		fmt.Sprintf("http://localhost:%s", fport),
		fmt.Sprintf("amqp://jeff:jeff@localhost:5672"),
	}

	// create grain worker
	gw := NewGrainworker(config.FunctionName, config.FunctionUrl, config.BrokerUrl)
	log.Println("created new grain worker:", gw.String())

	// start it
	go gw.start()
	log.Println("started grain worker")

	// don't let the main goroutine end immediately
	<-make(<-chan any)
}
