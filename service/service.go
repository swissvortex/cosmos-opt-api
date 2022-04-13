package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Colm3na/cosmos-opt-api/constants"
	"github.com/Colm3na/cosmos-opt-api/logger"
	"github.com/Colm3na/cosmos-opt-api/models"
	"github.com/Colm3na/cosmos-opt-api/repository"
)

type Service interface {
	GetValidatorUptime(cosmosvaloper string) (models.Uptime, error)
	GetBlockHeighAndTime(blockHeigh int) (int, time.Time, error)
	GetAverageBlockTime(latestBlockHeigh int, latestBlockTime time.Time) (float64, error)
}

type service struct {
	repository repository.Repository
	log        logger.Logger
}

func NewService(repository repository.Repository, log logger.Logger) Service {
	return &service{
		repository: repository,
		log:        log,
	}
}

func (s *service) GetValidatorUptime(cosmosvaloper string) (models.Uptime, error) {
	s.log.EntryWithContext(s.log.FileContext(), cosmosvaloper)

	uri := "https://api.cosmostation.io/v1/staking/validator/" + cosmosvaloper
	responseData := s.repository.HttpGetBody(uri)

	validator, err := models.UnmarshalValidator(responseData)
	if err != nil {
		s.log.InternalErrorWithContext(s.log.FileContext(), err)
		return models.Uptime{}, err
	}

	s.log.ExitWithContext(s.log.FileContext(), validator.Uptime)
	return validator.Uptime, nil
}

func (s *service) GetBlockHeighAndTime(blockHeigh int) (int, time.Time, error) {
	s.log.EntryWithContext(s.log.FileContext(), blockHeigh)

	uri := "http://localhost:26657/block"
	if blockHeigh != constants.LatestBlock {
		uri = uri + "?height=" + strconv.Itoa(blockHeigh)
	}
	responseData := s.repository.HttpGetBody(uri)

	block, err := models.UnmarshalBlock(responseData)
	if err != nil {
		s.log.InternalErrorWithContext(s.log.FileContext(), err)
		return 0, time.Time{}, err
	}

	blockTime, err := time.Parse(time.RFC3339, block.Result.Block.Header.Time)
	if err != nil {
		s.log.InternalErrorWithContext(s.log.FileContext(), err)
		return 0, time.Time{}, err
	}

	blockHeigh, err = strconv.Atoi(block.Result.Block.Header.Height)
	if err != nil {
		s.log.InternalErrorWithContext(s.log.FileContext(), err)
		return 0, time.Time{}, err
	}

	s.log.ExitWithContext(s.log.FileContext(), blockHeigh, blockTime)
	return blockHeigh, blockTime, nil
}

func (s *service) GetAverageBlockTime(latestBlockHeigh int, latestBlockTime time.Time) (float64, error) {
	s.log.EntryWithContext(s.log.FileContext(), latestBlockHeigh, latestBlockTime)

	var blockTimeArray [20]time.Time
	var blockPeriodArray [20]int64

	uri := "http://localhost:26657/blockchain?minHeight=" + strconv.Itoa(latestBlockHeigh-21) + "&maxHeight=" + strconv.Itoa(latestBlockHeigh-1)
	responseData := s.repository.HttpGetBody(uri)
	blockchain, err := models.UnmarshalBlockchain(responseData)
	if err != nil {
		s.log.InternalErrorWithContext(s.log.FileContext(), err)
		return 0, err
	}
	for i := 0; i < 20; i++ {
		blockTimeArray[i], err = time.Parse(time.RFC3339, blockchain.BlockchainResult.BlockMetas[i].Header.Time)
		if err != nil {
			s.log.InternalErrorWithContext(s.log.FileContext(), err)
			return 0, err
		}
	}
	blockPeriodArray[0] = int64(latestBlockTime.Sub(blockTimeArray[0]) / time.Millisecond)
	avg := blockPeriodArray[0]
	for i := 1; i < 20; i++ {
		blockPeriodArray[i] = int64(blockTimeArray[i-1].Sub(blockTimeArray[i]) / time.Millisecond)
		avg = avg + blockPeriodArray[i]
	}
	avgBlockTime := float64(avg) / (20 * 1000)
	s.log.InfoMessageWithContext(s.log.FileContext(), fmt.Sprintf("Block array: %v\n", blockPeriodArray))
	s.log.InfoMessageWithContext(s.log.FileContext(), fmt.Sprintf("Average block time = %fs\n", avgBlockTime))

	s.log.ExitWithContext(s.log.FileContext(), avgBlockTime)
	return avgBlockTime, nil
}
