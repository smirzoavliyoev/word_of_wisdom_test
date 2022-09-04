package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeservice"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeusagefixer"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeusageservice"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/quoteservice"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/tcp"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/tcp/structs"
	"github.com/smirzoavliyoev/word_of_wisdom_test/pkg/config"
)

func main() {

	challengeService := challengeservice.NewChallengeService()
	challegeUsageService := challengeusageservice.NewChallengeUsageService()
	challengeusagefixer.NewChallengeUsageFixer(challengeService, challegeUsageService)
	quoteService := quoteservice.NewQuotaService()

	cfg, err := config.ReadConfig(config.WithSpecificConfigPathOption)
	if err != nil {
		panic(err)
	}

	server := tcp.NewServer(cfg, challengeService, challegeUsageService, quoteService)

	server.Handle(Handler)
}

func Handler(conn net.Conn,
	s tcp.Server,
	chalangeService *challengeservice.ChallengeService,
	challengeUsageService *challengeusageservice.ChallengeUsageService,
	quotaService *quoteservice.QuoteService) {

	var (
		responseBody        interface{}
		responseError       string
		responseMessageType structs.ResponseMessageType
	)

	defer conn.Close()

	//check ip error

	requestMsg, err := s.ReadMessage(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(requestMsg)

	switch requestMsg.Type {
	case structs.RequestChallenge:

		responseMessageType = structs.ResponseChallenge

		newChallenge, err := chalangeService.GetChallenge()
		if err != nil {
			fmt.Println(err)
			return
		}

		challengeUsageService.SaveChallengeUsage(conn.RemoteAddr().String(), newChallenge)

		responseBody = structs.NewResponseChallengeMessage(newChallenge)

	case structs.RequestQuote:
		responseMessageType = structs.ResponseQuote
		responseBody = structs.NewResponseQuoteMessage("")

		reqBodyJson, err := json.Marshal(requestMsg.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		var body structs.RequestQuoteMessage

		if err = json.Unmarshal(reqBodyJson, &body); err != nil {
			logrus.Error(err)
			return
		}

		quote, err := quotaService.GetQuota()
		if err != nil {
			responseError = err.Error()
			break
		}

		responseBody = structs.NewResponseQuoteMessage(quote)

	default:
		fmt.Println("figushki")
		return
	}

	responseMsg := structs.NewResponseMessage(requestMsg.Id, responseMessageType, responseBody)

	if responseError != "" {
		responseMsg.Error = responseError
	}

	if err = s.WriteMessage(conn, responseMsg); err != nil {
		logrus.Error(err)
	}
}
