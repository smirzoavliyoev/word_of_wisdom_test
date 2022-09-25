package tcp

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/tcp/structs"
	"github.com/smirzoavliyoev/word_of_wisdom_test/pkg/config"
)

type Client struct {
	config *config.Config
	conn   net.Conn
}

func NewClient(config *config.Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) connect() error {
	address := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *Client) getIp() string {
	return c.conn.LocalAddr().String()
}

func (c *Client) Request(req interface{}) (string, error) {
	if err := c.connect(); err != nil {
		return "", err
	}

	// requestBytes, err := json.Marshal(req)
	// if err != nil {
	// 	return err
	// }
	requestBytes := req.([]byte)

	reqBase64 := base64.StdEncoding.EncodeToString(requestBytes)

	fmt.Println(reqBase64)

	_, err := c.conn.Write([]byte(reqBase64 + "\n"))
	if err != nil {
		return "", nil
	}
	return c.getIp(), nil
}

func (c *Client) Response() (*structs.ResponseMessage, error) {
	defer c.conn.Close()

	resp := &structs.ResponseMessage{}
	reader := bufio.NewReader(c.conn)
	respBase64, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		return nil, err
	}

	if err == io.EOF {
		fmt.Println("this is eof")
	}

	respData, err := base64.StdEncoding.DecodeString(respBase64)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(respData))

	if err = json.Unmarshal(respData, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
