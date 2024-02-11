use this to create a custom function that can be used to test the performance of the platform

```go
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
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/nuclio/nuclio-sdk-go"
)

const (
	//INVOCATION_URL = "http://172.17.0.1:8888/evaluation/headers"
	INVOCATION_URL = "http://host.docker.internal:8888/evaluation/headers"
)

func primeNumbers(max uint64) []uint64 {
	var primes []uint64

	for i := uint64(2); i < max; i++ {
		isPrime := true

		for j := uint64(2); j <= uint64(math.Sqrt(float64(i))); j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}

		if isPrime {
			primes = append(primes, i)
		}

        time.Sleep(5 * time.Microsecond)
	}

	return primes
}

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {

	start := time.Now()

	headers := event.GetHeaders()

	fmt.Printf("Headers: %v\n", headers)

	n, err := strconv.ParseUint(string(headers["Max"].(string)), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing input:", err)
		return nil, err
	}

	primes := primeNumbers(n)

	req, err := http.NewRequest(http.MethodPost, INVOCATION_URL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value.(string))
	}

	req.Header.Set("X-Profaastinate-Exec-Start", start.Format(time.RFC3339))
	req.Header.Set("X-Profaastinate-Exec-Stop", time.Now().Format(time.RFC3339))

	// Make the request
	resp, postErr := http.DefaultClient.Do(req)
	if postErr != nil {
		return nil, postErr
	}

	defer resp.Body.Close()

	//print req header
	fmt.Printf("Request headers: %v\n", req.Header)

	return nuclio.Response{
		StatusCode:  200,
		ContentType: "application/text",
		Headers:     headers,
		Body:        []byte(fmt.Sprintf("Found %d prime numbers", len(primes))),
	}, nil
}

```