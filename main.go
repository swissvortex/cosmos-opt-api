package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Colm3na/cosmos-opt-api/models"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		latestBlockHeigh, latestBlockTime := GetBlockHeighAndTime(-1)
		return c.String(http.StatusOK, fmt.Sprintf("%f", GetAverageBlockTime(latestBlockHeigh, latestBlockTime)))
	})
	e.Logger.Fatal(e.Start(":47567"))
}

func GetBlockHeighAndTime(blockHeigh int) (int, time.Time) {
	uri := "http://cosmos.delega.io:26657/block"
	if blockHeigh != -1 {
		uri = uri + "?height=" + strconv.Itoa(blockHeigh)
	}
	response, err := http.Get(uri)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	block, err := models.UnmarshalBlock(responseData)
	if err != nil {
		log.Fatal(err)
	}

	blockTime, err := time.Parse(time.RFC3339, block.Result.Block.Header.Time)
	if err != nil {
		log.Fatal(err)
	}

	blockHeigh, err = strconv.Atoi(block.Result.Block.Header.Height)
	if err != nil {
		log.Fatal(err)
	}

	return blockHeigh, blockTime
}

func GetAverageBlockTime(latestBlockHeigh int, latestBlockTime time.Time) float64 {
	var blockTimeArray [20]time.Time
	var blockPeriodArray [20]int64

	uri := "http://cosmos.delega.io:26657/blockchain?minHeight=" + strconv.Itoa(latestBlockHeigh-21) + "&maxHeight=" + strconv.Itoa(latestBlockHeigh-1)
	response, err := http.Get(uri)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	blockchain, err := models.UnmarshalBlockchain(responseData)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		blockTimeArray[i], err = time.Parse(time.RFC3339, blockchain.BlockchainResult.BlockMetas[i].Header.Time)
		if err != nil {
			log.Fatal(err)
		}
	}
	blockPeriodArray[0] = int64(latestBlockTime.Sub(blockTimeArray[0]) / time.Millisecond)
	avg := blockPeriodArray[0]
	for i := 1; i < 20; i++ {
		blockPeriodArray[i] = int64(blockTimeArray[i-1].Sub(blockTimeArray[i]) / time.Millisecond)
		avg = avg + blockPeriodArray[i]
	}
	fmt.Printf("%v\n", blockPeriodArray)
	fmt.Printf("Average block time = %fs\n", float64(avg)/(20*1000))

	return float64(avg) / (20 * 1000)
}
