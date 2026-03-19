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
	Proto  string
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
		Proto:  resp.Proto,
		Status: resp.Status,
		Header: resp.Header.Clone(),
		Body:   body,
		Err:    nil,
	}
}

func printResponse(r Result) {
	fmt.Printf("%s %s\n", r.Proto, r.Status)

	for k, values := range r.Header {
		for _, v := range values {
			fmt.Printf("%s: %s\n", k, v)
		}
	}

	fmt.Println(string(r.Body))
}

func main() {
	var timeoutShort int
	var timeoutLong int
	var showInf bool

	flag.IntVar(&timeoutShort, "t", -1, "timeout in seconds")
	flag.IntVar(&timeoutLong, "timeout", -1, "timeout in seconds")
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

	var timeoutSeconds int
	switch {
	case timeoutShort != -1 && timeoutLong != -1:
		fmt.Fprintln(os.Stderr, "error: use either -t or -timeout, not both")
		os.Exit(1)
	case timeoutShort != -1:
		timeoutSeconds = timeoutShort
	case timeoutLong != -1:
		timeoutSeconds = timeoutLong
	default:
		timeoutSeconds = 15
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

	received := 0

	for received < len(args) {
		r := <-results
		received++

		if r.Err != nil {
			continue
		}

		printResponse(r)
		cancel()
		os.Exit(0)
	}

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Fprintln(os.Stderr, "error: timeout")
		os.Exit(228)
	}
	fmt.Fprintln(os.Stderr, "error: all requests failed")
	os.Exit(1)
}
