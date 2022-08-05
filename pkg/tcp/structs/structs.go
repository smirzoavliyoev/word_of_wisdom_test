package structs

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1"
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

type ResponseChallengeMessage struct {
	Challenge Challenge `json:"challenge"`
}

func NewResponseChallengeMessage(challenge *Challenge) *ResponseChallengeMessage {
	return &ResponseChallengeMessage{
		Challenge: *challenge,
	}
}

type ResponseQuoteMessage struct {
	Quote string `json:"quote"`
}

func NewResponseQuoteMessage(quote string) *ResponseQuoteMessage {
	return &ResponseQuoteMessage{Quote: quote}
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

type RequestChallengeMessage struct{}

func NewRequestChallengeMessage() *RequestMessage {
	return NewRequestMessage(RequestChallenge, RequestChallengeMessage{})
}

type RequestQuoteMessage struct {
	Challenge Challenge `json:"challenge"`
}

func NewRequestQuoteMessage(challenge *Challenge) *RequestQuoteMessage {
	return &RequestQuoteMessage{
		Challenge: *challenge,
	}
}

type Challenge struct {
	Complexity []byte `json:"complexity"`
	Timestamp  int64  `json:"timestamp"`
	Timeout    int64  `json:"timeout"`
	Entropy    []byte `json:"entropy"`
	Signature  []byte `json:"signature"`
	PublicKey  []byte `json:"public_key"`
	Solution   []byte `json:"solution"`
}

func NewChallenge(privateKeyBytes []byte, timeout int64) (*Challenge, error) {
	_, publicKey := secp256k1.PrivKeyFromBytes(privateKeyBytes)
	complexity, err := GetComplexity()
	if err != nil {
		return nil, err
	}

	entropy, err := GenerateEntropy()
	if err != nil {
		return nil, err
	}

	challenge := &Challenge{
		Complexity: complexity,
		Timestamp:  time.Now().Unix(),
		Timeout:    timeout,
		Entropy:    entropy,
		PublicKey:  publicKey.SerializeCompressed(),
	}

	return challenge, nil
}

func (c Challenge) CheckSolution() bool {
	solution := big.NewInt(0)
	solution.SetBytes(c.Solution)
	complexity := big.NewInt(0)
	complexity.SetBytes(c.Complexity)
	hash := big.NewInt(0)
	hash.SetBytes(c.Hash())
	res := hash.Cmp(complexity)

	if res == -1 {
		return true
	}

	return false
}

func (c Challenge) IsOverdue() bool {
	now := time.Now().Unix()
	deadline := c.Timestamp + c.Timeout

	return now >= deadline
}

func (c *Challenge) Solve() {
	solution := big.NewInt(0)
	solution.SetBytes(c.Solution)
	complexity := big.NewInt(0)
	complexity.SetBytes(c.Complexity)

	for {
		hash := big.NewInt(0)
		hash.SetBytes(c.Hash())
		res := hash.Cmp(complexity)

		if res == -1 {
			break
		}

		solution.Add(solution, big.NewInt(1))
		c.Solution = solution.Bytes()
	}
}

func (c Challenge) Hash() []byte {
	h := sha256.New()
	h.Write(c.Complexity)

	timestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestamp, uint64(c.Timestamp))
	h.Write(timestamp)

	timeout := make([]byte, 8)
	binary.LittleEndian.PutUint64(timeout, uint64(c.Timeout))
	h.Write(timeout)

	h.Write(c.Entropy)
	h.Write(c.PublicKey)
	h.Write(c.Solution)

	return h.Sum(nil)
}

func (c *Challenge) Sign(privateKeyBytes []byte) error {
	privateKey, publicKey := secp256k1.PrivKeyFromBytes(privateKeyBytes)
	sig, err := privateKey.Sign(c.Hash())
	if err != nil {
		return err
	}

	c.Signature = sig.Serialize()
	c.PublicKey = publicKey.SerializeCompressed()

	return nil
}

func (c Challenge) VerifySignature(privateKeyBytes []byte) (bool, error) {
	_, publicKey := secp256k1.PrivKeyFromBytes(privateKeyBytes)
	c.Solution = nil
	sig, err := secp256k1.ParseSignature(c.Signature)
	if err != nil {
		return false, err
	}

	return sig.Verify(c.Hash(), publicKey), nil
}

func GetComplexity() ([]byte, error) {
	return []byte{0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}, nil
}

func GenerateEntropy() ([]byte, error) {
	entropy := make([]byte, 32)
	_, err := rand.Read(entropy)
	if err != nil {
		return nil, err
	}

	return entropy, nil
}
