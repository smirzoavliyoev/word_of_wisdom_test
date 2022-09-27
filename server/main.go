package main

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeservice"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeusagefixer"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/challengeusageservice"
	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/hashverifier"
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

	fmt.Println("---- Request message ---- ", requestMsg)

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
		fmt.Println(1212)
		responseMessageType = structs.ResponseQuote

		var body map[string]interface{}

		fmt.Println(requestMsg.Body)
		body, ok := requestMsg.Body.(map[string]interface{})
		if !ok {
			fmt.Println("a vot xren tebe")
			return
		}

		fmt.Println("----body----- ", body)

		str, ok := body["challenge"].(string)
		if !ok {
			fmt.Println("no chaleng body - ", body["challenge"])
			return
		}

		fmt.Println("str", str)

		fmt.Println(body["challenge"], "here", str)

		hc, err := hashverifier.New(
			&hashverifier.Resource{
				ValidatorFunc: validatorFunc,
				Data:          str,
			},
			nil,
		)

		if err != nil {
			panic(err)
		}

		ok, err = hc.Verify(str)

		if err != nil {
			fmt.Println(err)
			return
		}

		if !ok {
			fmt.Println("solution failed")
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

func validatorFunc(a string) bool {
	return true
}
