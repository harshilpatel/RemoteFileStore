package utils

import "net/rpc"

type Client struct {
}

func (c *Client) MakeRequest(a ...interface{}) {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	c.
}
