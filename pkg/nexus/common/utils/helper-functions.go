package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/nuclio/nuclio/pkg/common/headers"
	"github.com/nuclio/nuclio/pkg/nexus/common/models"
)

// GetEnvironmentHost returns the host of the environment
//
// Currently for linux and mac os it is host.docker.internal
// We set in the docker-compose an external host to it for ensuring that the host.docker.internal will be resolved
// More info: https://docs.docker.com/docker-for-mac/networking/#use-cases-and-workarounds
// Docker compose: profaastinate/deployment/docker/docker-compose.yml
func GetEnvironmentHost() (host string) {
	return "localhost"
}

// TransformRequestToClientRequest transforms the async request send to the nexus from outside the cluster to a request
// that can be sent to the function inside the cluster
func TransformRequestToClientRequest(nexusItemRequest *http.Request) (newRequest *http.Request) {
	var requestUrl url.URL
	requestUrl.Scheme = "http"
	requestUrl.Path = nexusItemRequest.URL.Path
	requestUrl.Host = fmt.Sprintf("%s:%s", GetEnvironmentHost(), models.PORT)

	newRequest, _ = http.NewRequest(nexusItemRequest.Method, requestUrl.String(), nexusItemRequest.Body)

	// Create a new header map and copy the contents
	newRequest.Header = make(http.Header)
	for name, values := range nexusItemRequest.Header {
		for _, value := range values {
			if name == headers.ProcessDeadline {
				continue
			}

			newRequest.Header.Add(name, value)
		}
	}

	// fmt.Println("new Request: ", newRequest)
	return
}

// SetEvaluationHeaders sets the headers for the evaluation request
func SetEvaluationHeaders(req *http.Request, schedulerName string) {
	if schedulerName == "" || req.Header.Get(headers.INCOMING) == "" {
		req.Header.Set(headers.INCOMING, time.Now().Format(time.RFC3339))
	} else {
		req.Header.Set(headers.SYNC_PROCESSING, time.Now().Format(time.RFC3339))
		req.Header.Set(headers.SCHEDULER_NAME, schedulerName)
	}
}
