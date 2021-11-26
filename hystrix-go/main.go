package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

type Handle struct{}

func main() {
	go serve()
		
	url := os.Getenv("URL")
	for range time.Tick(time.Second) {
		concurrent := getEnvInt("CONCURRENT")

		wg := &sync.WaitGroup{}
		for i := 0; i < concurrent; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Printf("return error: %+v\n", sendRequest(url))
			}()
		}
		wg.Wait()
	}
}

func sendRequest(url string) error {
	hystrix.ConfigureCommand("my-command", hystrix.CommandConfig{
		Timeout:                int(time.Duration(getEnvInt("TIMEOUT")) * time.Millisecond),
		MaxConcurrentRequests:  getEnvInt("MAX_CONCURRENT_REQUESTS"),
		SleepWindow:            getEnvInt("SLEEP_WINDOW"),
		RequestVolumeThreshold: getEnvInt("REQUEST_VOLUME_THRESHOD"),
		ErrorPercentThreshold:  getEnvInt("ERROR_PERCENT_THRESHOLD"),
	})

	return hystrix.Do("my-command", func() error {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("get error: %v\n", err)
			return err
		}

		defer resp.Body.Close()
		content, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("get response: %s\n", content)

		return nil
	}, func(err error) error {
		fmt.Printf("handle error: %v\n", err)
		return nil
	})
}

func getEnvInt(env string) int {
	res, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		return 0
	}
	return res
}

func serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request)  {
		time.Sleep(time.Duration(getEnvInt("SLEEP_TIME")) * time.Second)
		w.Write([]byte("Hello world"))
	})
	http.ListenAndServe(":8888", nil)
}
