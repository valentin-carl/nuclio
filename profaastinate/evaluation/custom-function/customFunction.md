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
	}

	return primes
}

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	req, err := http.NewRequest(http.MethodPost, INVOCATION_URL, nil)
	if err != nil {
		context.Logger.Error("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("x-nuclio-function-name", "yeet")

	// Make the request
	resp, postErr := http.DefaultClient.Do(req)
	if postErr != nil {
		context.Logger.Error("Error sending request to evaluation endpoint:", postErr)
		return nil, postErr
	}

	defer resp.Body.Close()

	n, err := strconv.ParseUint(string(event.GetBody()), 10, 64)

	if err != nil {
		context.Logger.Error("Error parsing input:", err)
		return nil, err
	}

	primes := primeNumbers(n)

	return nuclio.Response{
		StatusCode:  200,
		ContentType: "application/text",
		Body:        []byte(fmt.Sprintf("Found %d prime numbers", len(primes))),
	}, nil
}
```