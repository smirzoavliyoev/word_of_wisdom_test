package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/tcp"
	"github.com/smirzoavliyoev/word_of_wisdom_test/pkg/config"
)

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func main() {

	// ip, err := externalIP()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(ip)

	cfg, err := config.ReadConfig(config.WithSpecificConfigPathOption)
	if err != nil {
		panic(err)
	}

	client := tcp.NewClient(cfg)

	b, err := ioutil.ReadFile("./example.json")
	if err != nil {
		panic(err)
	}

	ip, err := client.Request(b)
	if err != nil {
		panic(err)
	}

	response, err := client.Response()
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	if err == io.EOF {
		fmt.Println("this is eof", err)
	}

	fmt.Println("response", response, ip)

	hc, err := New(
		&Resource{
			Data:          ip,
			ValidatorFunc: validateResource,
		},
		nil, // use default config.
	)
	if err != nil {
		fmt.Println("error solution finder", err)
		return
	}

	solution, err := hc.Compute()
	if err != nil {
		if err != ErrSolutionFail {
			// did not find a solution, can call compute again.
			fmt.Println("no any solutions")
		}
	}
	fmt.Println("sloution", solution)

	fmt.Println("response", response)

}

func validateResource(data string) bool {
	return true
}

func GetQuota(conn net.Conn) {
	defer conn.Close()

}
