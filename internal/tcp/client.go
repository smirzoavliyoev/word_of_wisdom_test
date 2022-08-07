package tcp

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func (c *Client) Request(req interface{}) error {
	if err := c.connect(); err != nil {
		return err
	}

	requestBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	reqBase64 := base64.StdEncoding.EncodeToString(requestBytes)

	_, err = c.conn.Write([]byte(reqBase64 + "\n"))
	return err
}

func (c *Client) Response() (*structs.ResponseMessage, error) {
	defer c.conn.Close()

	resp := &structs.ResponseMessage{}
	reader := bufio.NewReader(c.conn)
	respBase64, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	respData, err := base64.StdEncoding.DecodeString(respBase64)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(respData, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// func (c Client) RequestChallenge() (*structs.Challenge, error) {
// 	req := structs.NewRequestChallengeMessage()

// 	if err := c.request(req); err != nil {
// 		return nil, err
// 	}

// 	response, err := c.response()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if response.Type != structs.ResponseChallenge {
// 		return nil, errors.New("response type mismatch") // TODO: make const and move out of here
// 	}

// 	bodyJson, err := json.Marshal(response.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var body structs.ResponseChallengeMessage

// 	if err = json.Unmarshal(bodyJson, &body); err != nil {
// 		return nil, err
// 	}

// 	return &body.Challenge, nil
// }

// func (c Client) RequestQuote(solvedChallenge *structs.Challenge) (string, error) {
// 	reqBody := structs.NewRequestQuoteMessage(solvedChallenge)
// 	req := structs.NewRequestMessage(structs.RequestQuote, reqBody)

// 	if err := c.request(req); err != nil {
// 		return "", err
// 	}

// 	response, err := c.response()
// 	if err != nil {
// 		return "", err
// 	}

// 	if response.Error != "" {
// 		return "", errors.New(response.Error)
// 	}

// 	bodyJson, err := json.Marshal(response.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	var body structs.ResponseQuoteMessage

// 	if err = json.Unmarshal(bodyJson, &body); err != nil {
// 		return "", err
// 	}

// 	return body.Quote, nil
// }
