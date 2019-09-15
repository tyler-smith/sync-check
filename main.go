package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gcash/bchd/rpcclient"
	"github.com/tyler-smith/sync-check/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	httpClient = &http.Client{Timeout: time.Second * 60}
)

func main() {
	heights, maxHeight, err := getHeights(getHostsFromArgs())
	if err != nil {
		log.Fatal(err)
	}

	failedCount := 0
	for host, height := range heights {
		if height < (maxHeight - 3) {
			failedCount++
			fmt.Println(fmt.Sprintf(
				"%s is behind by %d blocks (%d)", host, maxHeight-height, height))
		}
	}

	os.Exit(failedCount)
}

func getHostsFromArgs() []string {
	if len(os.Args) < 2 {
		return []string{"grpc://bchd.greyh.at:8335"}
	}
	return os.Args[1:]
}

func getHeights(hosts []string) (map[string]int64, int64, error) {
	heights := make(map[string]int64, len(hosts)+1)
	var maxHeight int64

	// Query each host for their height
	for _, host := range hosts {
		u, err := url.Parse(host)
		if err != nil {
			return nil, 0, err
		}

		host = u.Hostname() + ":" + u.Port()

		var height int64
		switch u.Scheme {
		case "grpc":
			height, err = getBestHeightGRPC(host)
		case "rpc":
			height, err = getBestHeightRPC(host)
		default:
			return nil, 0, errors.New("Unknown protocol scheme: " + u.Scheme)
		}

		if err != nil {
			return nil, 0, err
		}

		heights[host] = height
		if height > maxHeight {
			maxHeight = height
		}
	}

	// Query bitcoin.com's API
	bitcoinDotComheight, err := getBestHeightBitcoinDotCom()
	if err != nil {
		return nil, 0, err
	}

	if bitcoinDotComheight > maxHeight {
		maxHeight = bitcoinDotComheight
	}

	heights["bitcoin.com"] = bitcoinDotComheight

	return heights, maxHeight, nil
}

func getBestHeightGRPC(host string) (int64, error) {
	conn, err := grpc.Dial(
		host,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return 0, err
	}

	client := pb.NewBchrpcClient(conn)
	blockchainInfo, err := client.GetBlockchainInfo(
		context.Background(), &pb.GetBlockchainInfoRequest{})
	if err != nil {
		return 0, err
	}

	return int64(blockchainInfo.GetBestHeight()), nil
}

func getBestHeightRPC(host string) (int64, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         "",
		Pass:         "",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return 0, err
	}
	defer client.Shutdown()

	blockCount, err := client.GetBlockCount()
	if err != nil {
		return 0, err
	}

	return blockCount, nil
}

func getBestHeightBitcoinDotCom() (int64, error) {
	data := &struct {
		Blocks int64 `json:"blocks"`
	}{}

	resp, err := httpClient.Get(
		"https://rest.bitcoin.com/v2/blockchain/getBlockchainInfo")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(respBody, data)
	if err != nil {
		return 0, err
	}

	return data.Blocks, nil
}
