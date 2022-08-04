package structs

import (
	"time"

	"github.com/google/uuid"
)

type RequestMessageType uint

const (
	RequestChallenge RequestMessageType = iota
	RequestQuote
)

type RequestMessage struct {
	Id   string             `json:"id"`
	Type RequestMessageType `json:"type"`
	Body interface{}        `json:"body"`
}

func NewRequestMessage(requestType RequestMessageType, body interface{}) *RequestMessage {
	return &RequestMessage{
		Id:   uuid.New().String(),
		Type: requestType,
		Body: body,
	}
}

type ResponseMessageType uint

const (
	ResponseChallenge ResponseMessageType = iota
	ResponseQuote
)

type ResponseMessage struct {
	Id        string              `json:"id"`
	Type      ResponseMessageType `json:"type"`
	Error     string              `json:"error"`
	Timestamp int64               `json:"timestamp"`
	Body      interface{}         `json:"body"`
}

func NewResponseMessage(Id string, responseType ResponseMessageType, body interface{}) *ResponseMessage {
	responseMsg := &ResponseMessage{
		Id:        Id,
		Type:      responseType,
		Timestamp: time.Now().Unix(),
		Body:      body,
	}

	return responseMsg
}
