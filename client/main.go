package main

import (
	"net"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/tcp"
	"github.com/smirzoavliyoev/word_of_wisdom_test/pkg/config"
)

func main() {

	cfg, err := config.ReadConfig(config.WithSpecificConfigPathOption)
	if err != nil {
		panic(err)
	}

	client := tcp.NewClient(cfg)

	for {

		client.Request()

		client.Response()

	}
}

func GetQuota(conn net.Conn) {
	defer conn.Close()

}
