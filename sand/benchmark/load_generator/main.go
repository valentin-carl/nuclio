package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {

	functionName := os.Args[1]
	fmt.Println(functionName)

	maxValue, err := strconv.Atoi(os.Args[2])
	handle(err)
	fmt.Println(maxValue)

	fileOut := os.Args[3]
	fmt.Println(fileOut)

	measureWithoutDashboard := func(value int) int64 {
		//cmd := exec.Command("kubectl", "exec", "nuclio-entrytask-659cbcd9fb-llhbn", "--", "curl", "localhost:8080", "-d", "12")
		cmd := exec.Command("kubectl", "exec", functionName, "--", "curl", "localhost:8080", "-d", strconv.Itoa(value))

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		tStart := time.Now().UnixMilli()
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
		tEnd := time.Now().UnixMilli()
		return tEnd - tStart
	}

	type measurement struct {
		value   int
		latency int64
	}

	// student in der prüfungsphase sein wie:
	findSleep := func(n int) int {
		if n < 100 {
			return 1
		} else if n < 1000 {
			return 5
		} else {
			return 15
		}
	}

	file, err := os.Create(fileOut)
	handle(err)
	defer file.Close()
	_, err = file.WriteString("value,latency\n")

	for i := 0; i <= maxValue; i += 100 {
		m := measurement{
			value:   i,
			latency: measureWithoutDashboard(i),
		}

		_, err = file.WriteString(fmt.Sprintf("%d,%d\n", m.value, m.latency))
		handle(err)

		sleepFor := findSleep(i)
		time.Sleep(time.Duration(sleepFor) * time.Second)
	}

}

func handle(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
