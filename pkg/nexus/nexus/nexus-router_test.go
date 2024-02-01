package nexus

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type HelperSuite struct {
	suite.Suite
	*http.Client
	testServer  *httptest.Server
	nexusRouter *NexusRouter
}

func (helperSuite *HelperSuite) SetupTest() {
	nexus := Initialize()

	helperSuite.nexusRouter = NewNexusRouter(nexus)
	helperSuite.nexusRouter.Initialize()
	helperSuite.Client = &http.Client{}

	helperSuite.testServer = httptest.NewServer(helperSuite.nexusRouter.Router)
}

func (helperSuite *HelperSuite) TearDownTest() {
	helperSuite.testServer.Close()
}

func (helperSuite *HelperSuite) TestInitialize() {
	testCases := []struct {
		name       string
		method     string
		path       string
		statusCode int
	}{
		{"StartScheduler", http.MethodPost, "/scheduler/deadline/start", http.StatusOK},
		{"StopScheduler", http.MethodPost, "/scheduler/deadline/stop", http.StatusOK},
		{"GetAllSchedulersWithStatus", http.MethodGet, "/scheduler", http.StatusOK},
		{"modifyLoadBalancer", http.MethodPut, "/load-balancer", http.StatusOK},
		{"startLoadBalancer", http.MethodPost, "/load-balancer/start", http.StatusOK},
		{"stopLoadBalancer", http.MethodPost, "/load-balancer/stop", http.StatusOK},
	}

	for _, tc := range testCases {
		helperSuite.Run(tc.name, func() {
			req, err := http.NewRequest(tc.method, helperSuite.testServer.URL+tc.path, nil)
			assert.NoError(helperSuite.T(), err)

			resp, respErr := helperSuite.Client.Do(req)
			assert.NoError(helperSuite.T(), respErr)

			assert.Equal(helperSuite.T(), tc.statusCode, resp.StatusCode)
		})
	}

}

func (helperSuite *HelperSuite) TestModifyLoadBalancer() {
	queryParams := url.Values{}
	queryParams.Add("targetLoadCPU", "10")
	queryParams.Add("targetLoadMemory", "10")
	queryParams.Add("maxParallelRequests", "14")
	queryParams.Add("limitMaxParallelRequests", "12")

	pathWithQuery := fmt.Sprintf("/load-balancer?%s", queryParams.Encode())
	req, err := http.NewRequest(http.MethodPut, helperSuite.testServer.URL+pathWithQuery, nil)
	assert.NoError(helperSuite.T(), err)

	resp, respErr := helperSuite.Client.Do(req)
	assert.NoError(helperSuite.T(), respErr)

	assert.Equal(helperSuite.T(), http.StatusAccepted, resp.StatusCode)

	body, readErr := io.ReadAll(resp.Body)
	assert.NoError(helperSuite.T(), readErr)
	assert.Equal(helperSuite.T(), "Target CPU load set to 10.0\nTarget memory load set to 10.0\nMax parallel requests set to 14\nLimit max parallel requests set to 12\n", string(body))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(HelperSuite))
}
