package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

// ----------------------------------------------------------------------------
// Request and Response
// ----------------------------------------------------------------------------

// clientRequest represents a JSON-RPC request sent by a client.
type ClientRequest struct {
	// A String containing the name of the method to be invoked.
	Method string `json:"method"`
	// Object to pass as request parameter to the method.
	Params [1]interface{} `json:"params"`
	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	Id uint64 `json:"id"`
}

// clientResponse represents a JSON-RPC response returned to a client.
type ClientResponse struct {
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
	Id     uint64           `json:"id"`
}

// EncodeClientRequest encodes parameters for a JSON-RPC client request.
func EncodeClientRequest(method string, args interface{}) ([]byte, error) {
	// fmt.Println("args : ",args)
	c := &ClientRequest{
		Method: method,
		Params: [1]interface{}{args},
		Id:     uint64(rand.Int63()),
	}
	return json.Marshal(c)
}

// DecodeClientResponse decodes the response body of a client request into
// the interface reply.
func DecodeClientResponse(r io.Reader, reply interface{}) error {
	var c ClientResponse

	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return err
	}

	if c.Error != nil {
		return fmt.Errorf("%v", c.Error)
	}

	if c.Result == nil {
		return errors.New("result is null")
	}

	return json.Unmarshal(*c.Result, reply)
}

type StockArgs struct {
	StockSymbolandPercentage string
	Budget                   float32
}

type StockReply struct {
	Tradeid        int
	Stocks         string
	UnvestedAmount float32
}

type CheckArgs struct {
	Tradeid int
}
type CheckReply struct {
	Stocks           string
	UnvestedAmount   float32
	CurrentMarketVal float32
}

func Execute(method string, req, res interface{}) error {

	buf, _ := EncodeClientRequest(method, req)
	body := bytes.NewBuffer(buf)

	r, _ := http.NewRequest("POST", "http://localhost:8080/rpc/", body)
	r.Header.Set("Content-Type", "application/json")
	fmt.Println("calling http client")
	client := &http.Client{}
	fmt.Println("called http client")
	resp, _ := client.Do(r)
	fmt.Println("initiated resp")
	fmt.Println("response ", resp)
	fmt.Println("response status ", resp.Status)
	return DecodeClientResponse(resp.Body, res)
}

func main() {
	var reply StockReply
	var reply2 CheckReply
	inputArgs := os.Args[1:]

	if len(inputArgs) == 2 {
		var args StockArgs
		args.StockSymbolandPercentage = inputArgs[0]
		budget64, err := strconv.ParseFloat(inputArgs[1], 64)
		if err != nil {
			// error condition
		}
		args.Budget = float32(budget64)
		fmt.Println(args.Budget)

		if err := Execute("StockService.Say", &args, &reply); err != nil {
			fmt.Println("Error!")
		} else {
			fmt.Println(reply)
		}
	} else {
		var args CheckArgs
		i, err := strconv.Atoi(inputArgs[0])
		if err != nil {
			// error condition
		}
		args.Tradeid = i
		fmt.Println(args.Tradeid)
		if err := Execute("CheckPortfolioService.Che", &args, &reply2); err != nil {
			fmt.Println("Error!")
		} else {
			fmt.Println(reply2)
		}
	}
}
