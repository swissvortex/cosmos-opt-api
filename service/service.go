package service

import (
	"fmt"
	"strconv"
	"time"

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
	logger     *logger.Logger
}

func NewService(repository repository.Repository, logger *logger.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s *service) GetValidatorUptime(cosmosvaloper string) (models.Uptime, error) {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), cosmosvaloper)

	uri := "https://api.cosmostation.io/v1/staking/validator/" + cosmosvaloper
	responseData := s.repository.HttpGetBody(uri)

	validator, err := models.UnmarshalValidator(responseData)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return models.Uptime{}, err
	}

	s.logger.LogOnExitWithContext(s.logger.GetContext(), validator.Uptime)
	return validator.Uptime, nil
}

func (s *service) GetBlockHeighAndTime(blockHeigh int) (int, time.Time, error) {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), blockHeigh)

	uri := "http://cosmos.delega.io:26657/block"
	if blockHeigh != -1 {
		uri = uri + "?height=" + strconv.Itoa(blockHeigh)
	}
	responseData := s.repository.HttpGetBody(uri)

	block, err := models.UnmarshalBlock(responseData)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return 0, time.Time{}, err
	}

	blockTime, err := time.Parse(time.RFC3339, block.Result.Block.Header.Time)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return 0, time.Time{}, err
	}

	blockHeigh, err = strconv.Atoi(block.Result.Block.Header.Height)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return 0, time.Time{}, err
	}

	s.logger.LogOnExitWithContext(s.logger.GetContext(), blockHeigh, blockTime)
	return blockHeigh, blockTime, nil
}

func (s *service) GetAverageBlockTime(latestBlockHeigh int, latestBlockTime time.Time) (float64, error) {
	s.logger.LogOnEntryWithContext(s.logger.GetContext(), latestBlockHeigh, latestBlockTime)

	var blockTimeArray [20]time.Time
	var blockPeriodArray [20]int64

	uri := "http://cosmos.delega.io:26657/blockchain?minHeight=" + strconv.Itoa(latestBlockHeigh-21) + "&maxHeight=" + strconv.Itoa(latestBlockHeigh-1)
	responseData := s.repository.HttpGetBody(uri)
	blockchain, err := models.UnmarshalBlockchain(responseData)
	if err != nil {
		s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
		return 0, err
	}
	for i := 0; i < 20; i++ {
		blockTimeArray[i], err = time.Parse(time.RFC3339, blockchain.BlockchainResult.BlockMetas[i].Header.Time)
		if err != nil {
			s.logger.LogOnInternalErrorWithContext(s.logger.GetContext(), err)
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
	fmt.Printf("%v\n", blockPeriodArray)
	fmt.Printf("Average block time = %fs\n", avgBlockTime)

	s.logger.LogOnExitWithContext(s.logger.GetContext(), avgBlockTime)
	return avgBlockTime, nil
}
