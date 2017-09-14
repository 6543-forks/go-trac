package trac

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client handles communications with Trac JSONRPC.
type Client struct {
	server     string // https://user:passwd@trac.example.com/login/jsonrpc
	httpClient *http.Client

	// RPC functions
	Search *Search
	System *System
	Ticket *Ticket
	Wiki   *Wiki
}

// Request is send to Trac JSONRPC via a HTTP POST request.
type Request struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

// Response represents a response returned by Trac JSONRPC.
type Response struct {
	Error  RPCError        `json:"error,omitempty"`
	ID     string          `json:"id,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

// RPCError is the RPC error returned within the response.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

func (r *RPCError) Error() string {
	return fmt.Sprintf("%v(%d): %v", r.Name, r.Code, r.Message)
}

// NewClient returns a new Trac JSONRPC client.
func NewClient(server string) *Client {
	c := &Client{
		server: server,
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	// RPC exported functions
	c.Search = &Search{client: c}
	c.System = &System{client: c}
	c.Ticket = &Ticket{client: c}
	c.Wiki = &Wiki{client: c}
	return c
}

// Query sends a Request and returns a Response.
// Response.Result is unmarshaled by Client.Do
func (c *Client) Query(function string, params ...interface{}) (Response, error) {
	var response = Response{}
	query := Request{function, params}
	body, err := json.Marshal(query)
	if err != nil {
		return response, err
	}

	fmt.Printf("%v\n", string(body))

	res, err := http.Post(c.server, "application/json", bytes.NewReader(body))
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}

	if err := json.Unmarshal(resp, &response); err != nil {
		return response, err
	}
	if response.Error.Code != 0 {
		return response, &response.Error
	}
	return response, nil
}

// Do wraps Client.Query to unmarshal Response.Result in the value pointed to
// by v
func (c *Client) Do(function string, v interface{}, params ...interface{}) (interface{}, error) {
	r, err := c.Query(function, params...)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(r.Result, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// All returns a slice of names. To be used for endpoints which returns lists
// of names. E.g. components, milestones, priorities.
func (c *Client) All(function string) ([]string, error) {
	var r []string
	_, err := c.Do(function, &r)
	return r, err
}
