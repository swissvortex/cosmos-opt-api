package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Colm3na/cosmos-opt-api/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/blocktime", Blocktime)
	e.GET("/validator/:id", ValidatorUptime)
	e.Logger.Fatal(e.Start(":47567"))
}

func GetValidatorUptime(cosmosvaloper string) (models.Uptime, error) {
	uri := "https://api.cosmostation.io/v1/staking/validator/" + cosmosvaloper
	fmt.Println(uri)
	responseData := HttpGet(uri)

	validator, err := models.UnmarshalValidator(responseData)
	if err != nil {
		log.Fatal(err)
		return models.Uptime{}, errors.New("Cosmosvaloper not found")
	}

	return validator.Uptime, nil
}

func GetBlockHeighAndTime(blockHeigh int) (int, time.Time) {
	uri := "http://cosmos.delega.io:26657/block"
	if blockHeigh != -1 {
		uri = uri + "?height=" + strconv.Itoa(blockHeigh)
	}
	responseData := HttpGet(uri)

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
	responseData := HttpGet(uri)
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

func HttpGet(uri string) []byte {
	response, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData
}

func Blocktime(c echo.Context) error {
	latestBlockHeigh, latestBlockTime := GetBlockHeighAndTime(-1)
	blocktime := models.BlockTime{Average: GetAverageBlockTime(latestBlockHeigh, latestBlockTime)}
	marshalledBlocktime, err := json.Marshal(blocktime)
	if err != nil {
		log.Fatal(err)
	}
	return c.String(http.StatusOK, fmt.Sprintf("%s", string(marshalledBlocktime)))
}

func ValidatorUptime(c echo.Context) error {
	cosmosvaloper := c.Param("id")
	uptime, err := GetValidatorUptime(cosmosvaloper)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	marshalledUptime, err := json.Marshal(uptime)
	if err != nil {
		log.Fatal(err)
	}
	return c.String(http.StatusOK, fmt.Sprintf("%s", string(marshalledUptime)))
}
