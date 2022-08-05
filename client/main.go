package main

import (
	"fmt"
	"net"

	"github.com/smirzoavliyoev/word_of_wisdom_test/pkg/config"
	"github.com/smirzoavliyoev/word_of_wisdom_test/pkg/tcp"
)

func main() {

	cfg, err := config.ReadConfig(config.WithSpecificConfigPathOption)
	if err != nil {
		panic(err)
	}

	client := tcp.NewClient(cfg)

	for {

		var NoQuotasFlag = false

		client.Request()

		response, err := client.Response()
		if err != nil {
			fmt.Println(err)
		}
		resp := response.Body

		switch v := resp.(type) {
		case string:
			if v == "no quotas" {
				NoQuotasFlag = true
				break
			}
		}

		if NoQuotasFlag {
			break
		}
	}
}

func GetQuota(conn net.Conn) {
	defer conn.Close()

}
