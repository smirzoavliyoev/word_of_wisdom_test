package tcp

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeservice"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeusageservice"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/quoteservice"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/tcp/structs"
	"github.com/smirzoavliyoev/word_of_wisdom_test/pkg/config"
)

type Server struct {
	challengeService      *challengeservice.ChallengeService
	challengeusageservice *challengeusageservice.ChallengeUsageService
	quoteService          *quoteservice.QuoteService
	config                *config.Config
}

func NewServer(config *config.Config,
	challengeService *challengeservice.ChallengeService,
	challengeUsageService *challengeusageservice.ChallengeUsageService,
	quoteService *quoteservice.QuoteService) *Server {
	return &Server{
		config:                config,
		challengeService:      challengeService,
		challengeusageservice: challengeUsageService,
		quoteService:          quoteService,
	}
}

func (s Server) ReadMessage(conn net.Conn) (*structs.RequestMessage, error) {
	var (
		requestMsgBase64 string
		requestMsgData   []byte
		requestMsg       structs.RequestMessage
	)

	reader := bufio.NewReader(conn)
	requestMsgBase64, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	requestMsgData, err = base64.StdEncoding.DecodeString(requestMsgBase64)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(requestMsgBase64), string(requestMsgData))

	if err = json.Unmarshal(requestMsgData, &requestMsg); err != nil {
		return nil, err
	}

	return &requestMsg, nil
}

func (s Server) WriteMessage(conn net.Conn, responseMsg *structs.ResponseMessage) error {
	responseMsgData, err := json.Marshal(responseMsg)
	if err != nil {
		return err
	}

	responseMsgBase64 := base64.StdEncoding.EncodeToString(responseMsgData)
	_, err = conn.Write([]byte(responseMsgBase64 + "\n"))

	return err
}

func (s Server) Handle(handlefunc func(conn net.Conn,
	s Server,
	challengeService *challengeservice.ChallengeService,
	challengeUsageService *challengeusageservice.ChallengeUsageService,
	quoteService *quoteservice.QuoteService)) error {

	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	fmt.Println(address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	defer ln.Close()

	logrus.WithField("address", address).Info("Wisdom Server started")

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}

		go handlefunc(conn, s, s.challengeService, s.challengeusageservice, s.quoteService)
	}
}
