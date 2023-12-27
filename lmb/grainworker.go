package lmb

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nuclio/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// The grainworker is a separate program from Nuclio
// With each function (in the paper: "grain"), a grain worker is deployed.
// The grain worker subscribes to the respective queue and invokes the function.
// TODO open question: should the grain worker be deployed in a different container than the function but the same pod or in the same container?
// => currently leaning towards giving it a different container and to call the function from outside that container,
// but it could also be called from inside
// => look at what these containers are => if its alpine e.g., adding the grain worker is simple enough
// if its something else, this might be more difficult -- also, having the grain worker outside makes it available
// regardless of which runtime we're using
// maybe we still need one because of the deployment stuff?
// TODO how to deploy the grain worker automatically with the function?
// => do it manually for now
// Note: so far, this only supports async calls; in the paper it appears to
// only work that way, but maybe we could use rabbitmq's request-response pattern
// TODO add lots of logging => to file!
// TODO should this be included in "make build"? maybe create "make grainworker"?

type Grainworker struct {
	functionName string // get this as CLI argument
	functionIp   string // tells us where to reach the function TODO correct format? TODO what if both are in the same container?
	amqpUrl      string // also get this as CLI argument, tells us where to reach the queue
	conn         *amqp.Connection
	ch           *amqp.Channel
	queue        amqp.Queue   // TODO should this be a pointer? => the rabbitmq amqp implementation doesn't treat it that way
	client       *http.Client // used to invoke the functions directly
}

// TODO add different queue options here later on
// maybe in a config file?
func NewGrainworker(functionName, functionIp, amqpUrl string) *Grainworker {

	// create connection to broker
	conn, err := amqp.Dial(amqpUrl)
	handle(err)

	// the connection handles low level stuff, the channel gives access to most of the rabbitmq api
	ch, err := conn.Channel()
	handle(err)

	// create the queue if it doesn't exist yet
	// TODO turn some of the options to true later on
	// For now, they don't seem very relevant for the evaluation
	// But stuff like exclusive (at least for the consumer) etc. could be a good idea in general
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
		amqpUrl:      amqpUrl,
		conn:         conn,
		ch:           ch,
		queue:        queue,
		client:       &http.Client{},
	}
}

func (gw *Grainworker) start() {

	// TODO maybe change this to a different goroutine that listens to Ctrl + C
	// see for example https://stackoverflow.com/questions/18106749/golang-catch-signals
	defer gw.shutdown()

	// register at queue as consumer
	events, err := gw.ch.Consume(
		gw.functionName,
		fmt.Sprintf("grain-worker-%s", gw.functionName),
		false,
		false, // TODO this might actually be a good idea
		false,
		false,
		nil,
	)
	handle(err)

	// event loop: listen for requests and invoke the function
	for event := range events {

		// invoke the function
		go func() {
			err := gw.invoke(&event)
			handle(err)
		}()

		// send acknowledgement to broker
		err := event.Ack(false)
		handle(err)
	}

	// this shouldn't be reached
	log.Panic("events channel closed unexpectedly")
}

// Invoke calls the function this grainworker handles
func (gw *Grainworker) invoke(event *amqp.Delivery) error {

	// TODO create different go program to test invoking functions with args programmatically

	// create the request
	requestBody := bytes.NewReader(event.Body) // TODO can we just pass the body on?? experiment with this to figure it out
	req, err := http.NewRequest(
		http.MethodPost, // TODO does this have to be/should it be a post request?
		gw.functionIp,
		requestBody,
	)
	if err != nil {
		return err
	}

	// set the relevant nuclio headers
	headers := map[string]string{
		// TODO are these event necessary when calling the function directly?
		"x-nuclio-function-namespace": "default", // TODO is it "default" or "nuclio"?
		"x-nuclio-function-name":      gw.functionName,
		// TODO what else???
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// send it!
	response, err := gw.client.Do(req)
	if err != nil {
		return err
	}

	// TODO what to do with the response?
	// in workflows, we don't really expect a direct response from the function
	// => for now, just log it and check if it's a 200

	// read + log the response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	log.Printf("Response status: \"%s\"\nResponse body: \"%s\"", response.Status, string(responseBody))
	if response.StatusCode != http.StatusOK {
		return errors.New("didn't get 200 back, what happened? " + response.Status)
	}
	return nil
}

// Shutdown closes the amqp connection + channel
func (gw *Grainworker) shutdown() {

	// TODO do the go stuff where you listed to ctrl c to cancel and then close conn, channel and queue etc
	// => do that to determine when to call this function

	// close the channel
	err := gw.ch.Close()
	handle(err)

	// close the connection
	err = gw.conn.Close()
	handle(err)

	// TODO what else?
}

// //////////////// //
// helper functions //
// //////////////// //

// TODO handle errors better
func handle(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

// /////////////////// //
// program entry point //
// /////////////////// //

// blub
func main() {
	// TODO	get cli args, create and start the grain worker
}
