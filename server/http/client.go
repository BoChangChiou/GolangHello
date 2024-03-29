package http

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	interval = 4 * time.Second
)

var cancelContext context.Context

// Go當client打Http Get
func httpApiClientGet(ctx context.Context) {
	cancelContext = ctx

	time.Sleep(time.Second * 2)
	go fetchByTimer()
	time.Sleep(time.Second * 2)
	go fetchByTicker()
}

func newClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func fetchByTimer() {
	client := newClient()
	u, err := url.ParseRequestURI("http://localhost:9090/sqlGet")
	if err != nil {
		fmt.Println("ParseRequestURI err:", err)
		return
	}
	data := url.Values{}
	data.Set("Name", "Steven")
	u.RawQuery = data.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		fmt.Printf("NewReq err: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Do err: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Println("status code:", resp.StatusCode)
	log.Println(string(body))

	timer := time.NewTimer(interval)
	select {
	case <-timer.C:
		go fetchByTimer()
	case <- cancelContext.Done():
		fmt.Println("stop fetchByTimer")
		return
	}
}

func fetchByTicker() {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			go tickerClient()
		case <- cancelContext.Done():
			fmt.Println("stop fetchByTicker")
			return
		}
	}
}

func tickerClient() {
	client := newClient()
	req, err := http.NewRequest(http.MethodGet, "http://localhost:9090/sqlGet?Name=Joe", nil)
	if err != nil {
		fmt.Printf("NewReq err: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Do2 err: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Println("status code:", resp.StatusCode)
	log.Println(string(body))
}
