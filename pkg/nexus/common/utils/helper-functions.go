package utils

import (
	"fmt"
	"net/http"
	"net/url"

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
