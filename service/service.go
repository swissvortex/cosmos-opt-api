package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/swissvortex/cosmos-opt-api/constants"
	"github.com/swissvortex/cosmos-opt-api/logger"
	"github.com/swissvortex/cosmos-opt-api/models"
	"github.com/swissvortex/cosmos-opt-api/repository"
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

	uri := constants.CosmostationApi + cosmosvaloper
	responseData, err := s.repository.HttpGetBody(uri)
	if err != nil {
		s.log.InternalErrorWithContext(s.log.FileContext(), err)
		return models.Uptime{}, err
	}

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

	uri := constants.CosmosApiUrl + constants.BlockPath
	if blockHeigh != constants.LatestBlock {
		uri = uri + constants.BlockHeightParam + strconv.Itoa(blockHeigh)
	}
	responseData, err := s.repository.HttpGetBody(uri)
	if err != nil {
		s.log.InternalErrorWithContext(s.log.FileContext(), err)
		return 0, time.Time{}, err
	}

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

	var blockTimeArray [constants.AverageBlockWindow]time.Time
	var blockPeriodArray [constants.AverageBlockWindow]int64

	uri := constants.CosmosApiUrl + constants.BlockchainPath + constants.MinHeightParam + strconv.Itoa(latestBlockHeigh-constants.AverageBlockWindow-1) + constants.MaxHeightParam + strconv.Itoa(latestBlockHeigh-1)
	responseData, err := s.repository.HttpGetBody(uri)
	if err != nil {
		s.log.InternalErrorWithContext(s.log.FileContext(), err)
		return 0, err
	}
	blockchain, err := models.UnmarshalBlockchain(responseData)
	if err != nil {
		s.log.InternalErrorWithContext(s.log.FileContext(), err)
		return 0, err
	}
	for i := 0; i < constants.AverageBlockWindow; i++ {
		blockTimeArray[i], err = time.Parse(time.RFC3339, blockchain.BlockchainResult.BlockMetas[i].Header.Time)
		if err != nil {
			s.log.InternalErrorWithContext(s.log.FileContext(), err)
			return 0, err
		}
	}
	blockPeriodArray[0] = int64(latestBlockTime.Sub(blockTimeArray[0]) / time.Millisecond)
	avg := blockPeriodArray[0]
	for i := 1; i < constants.AverageBlockWindow; i++ {
		blockPeriodArray[i] = int64(blockTimeArray[i-1].Sub(blockTimeArray[i]) / time.Millisecond)
		avg = avg + blockPeriodArray[i]
	}
	avgBlockTime := float64(avg) / float64(constants.AverageBlockWindow*1000)
	s.log.DebugWithContext(s.log.FileContext(), fmt.Sprintf("Block array: %v\n", blockPeriodArray))
	s.log.DebugWithContext(s.log.FileContext(), fmt.Sprintf("Average block time = %fs\n", avgBlockTime))

	s.log.ExitWithContext(s.log.FileContext(), avgBlockTime)
	return avgBlockTime, nil
}
