package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Result struct {
	URL    string
	Status string
	Header http.Header
	Body   []byte
	Err    error
}

func fetch(ctx context.Context, client *http.Client, out chan<- Result, url string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		out <- Result{URL: url, Err: err}
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		out <- Result{URL: url, Err: err}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		out <- Result{URL: url, Err: err}
		return
	}

	out <- Result{
		URL:    url,
		Status: resp.Status,
		Header: resp.Header.Clone(),
		Body:   body,
		Err:    nil,
	}
}

func printResponse(r Result) {
	fmt.Printf("HTTP/1.1 %s\n", r.Status)

	for k, values := range r.Header {
		for _, v := range values {
			fmt.Printf("%s: %s\n", k, v)
		}
	}

	fmt.Println(string(r.Body))
}

func main() {
	var timeoutSeconds int
	var showInf bool

	flag.IntVar(&timeoutSeconds, "t", 15, "timeout in seconds")
	flag.IntVar(&timeoutSeconds, "timeout", 15, "timeout in seconds")

	flag.BoolVar(&showInf, "h", false, "show inf")
	flag.BoolVar(&showInf, "help", false, "show inf")

	flag.Usage = func() {
		fmt.Println(`hedgedcurl — hedged HTTP GET requests
	Usage:
	  hedgedcurl [flags] <url1> <url2> ...

	Flags:
	  -t, --timeout SECONDS   request timeout (default: 15)
	  -h, --help              show this help
	`)
	}

	flag.Parse()
	if showInf {
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "error: must be url")
		flag.Usage()
		os.Exit(1)
	}

	if timeoutSeconds <= 0 {
		fmt.Fprintln(os.Stderr, "error: timeout must be greater than zero")
		os.Exit(1)
	}

	timeout := time.Duration(timeoutSeconds) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := &http.Client{}
	results := make(chan Result, len(args))

	for _, url := range args {
		go fetch(ctx, client, results, url)
	}

	errors := 0
	for errors < len(args) {
		select {
		case r := <-results:
			if r.Err != nil {
				errors++
				continue
			}
			printResponse(r)
			cancel()
			os.Exit(0)
		case <-ctx.Done():
			fmt.Fprintln(os.Stderr, "error: timeout")
			os.Exit(228)
		}
	}
	fmt.Fprintln(os.Stderr, "error: all requests failed")
	os.Exit(1)
}
