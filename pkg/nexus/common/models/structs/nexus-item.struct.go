package structs

import (
	"net/http"
	"time"
)

type NexusItem struct {
	// The request to be processed. Later it will be used to create a new request.
	Request *http.Request `json:"request"`
	// The index of the NexusEntry in the queue.
	Index int `json:"index"`
	// The priority of the NexusEntry in the queue.
	Deadline time.Time `json:"x-nuclio-async-deadline"`
	// The name of the NexusEntry in the queue.
	Name string `json:"x-nuclio-function-name"`
}

// NewNexusItem allows to create a new NexusItem.
func NewNexusItem(request *http.Request, deadline time.Time, name string) *NexusItem {
	return &NexusItem{
		Request:  request,
		Deadline: deadline,
		Name:     name,
	}
}
