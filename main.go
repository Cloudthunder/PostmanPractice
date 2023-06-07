package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// variables to change for server ports and IP
const serverPort = 3777
const ipadd = "http://192.168.50.145"

type Machine struct {
	MachineNO    string
	MachineState int
}

func main() {
	//go func() {
	//	mux := http.NewServeMux()
	//	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//		fmt.Printf("server: %s /\n", r.Method)
	//	})
	//	server := http.Server{
	//		Addr:    fmt.Sprintf(":%d", serverPort),
	//		Handler: mux,
	//	}
	//	if err := server.ListenAndServe(); err != nil {
	//		if !errors.Is(err, http.ErrServerClosed) {
	//			fmt.Printf("error running http server: %s\n", err)
	//		}
	//	}
	//
	//	/*mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//		fmt.Printf("server: %s /\n", r.Method)
	//		fmt.Fprintf(w, `{"message": "hello!"}`)
	//	})*/
	//}()

	//delay to allow for server to start up
	time.Sleep(100 * time.Millisecond)

	//retrieving URL
	requestURL := fmt.Sprintf("%s:%d/machineStates/NC-132", ipadd, serverPort)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	//errors if retrieval failed
	if err != nil {
		fmt.Printf("client could not create request: %s\n", err)
		os.Exit(1)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	//printing status if things go well
	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	var givenMachine Machine
	err = json.Unmarshal(resBody, &givenMachine)
	if err != nil {
		panic(err)
	}
	//patch edit
	givenMachine.MachineState = 2

	//patch send
	jsonPayload, err := json.Marshal(givenMachine)
	if err != nil {
		panic(err)
	}

	// Create a Patch Request
	postReq, err := http.NewRequest(http.MethodPatch, requestURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		panic(err)
	}
	postReq.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 100 * time.Millisecond,
	}
	resp, err := client.Do(postReq)

	if err != nil {
		panic(err)
	}

	fmt.Printf("client: patch response successful %v", resp)

	//printing data
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", bytes.Replace(bytes.Replace(resBody, []byte("{"), []byte("{\n\t"), 1), []byte(",\""), []byte(",\n\t\""), -1))
}
