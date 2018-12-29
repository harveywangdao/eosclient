package test

import (
	"eosclient/eos"
	"eosclient/logger"
	"sync"
)

type EosClientTest struct {
	cli *eos.EosClient
}

func (e *EosClientTest) testApi() error {
	var err error

	err = e.cli.GetInfo()
	if err != nil {
		logger.Error(err)
		return nil
	}

	headBlockNumber, err := e.cli.GetHeadBlockNumber()
	if err != nil {
		logger.Error(err)
		return nil
	}

	err = e.cli.GetBlockByNum(headBlockNumber)
	if err != nil {
		logger.Error(err)
		return nil
	}

	err = e.cli.GetAccount("alice")
	if err != nil {
		logger.Error(err)
		return nil
	}

	contractAccount := "eosio.token"
	err = e.cli.GetABI(contractAccount)
	if err != nil {
		logger.Error(err)
		return nil
	}

	err = e.cli.GetCode(contractAccount)
	if err != nil {
		logger.Error(err)
		return nil
	}

	err = e.cli.GetBalance("bob", "SYS", contractAccount)
	if err != nil {
		logger.Error(err)
		return nil
	}

	err = e.cli.Transfer(contractAccount, "transfer", "alice", "bob", "26.0000 SYS")
	if err != nil {
		logger.Error(err)
		return nil
	}

	/*	err = e.cli.GetProducers()
		if err != nil {
			logger.Error(err)
			return nil
		}*/

	/*	err = e.cli.GetDBSize()
		if err != nil {
			logger.Error(err)
			return nil
		}*/

	return nil
}

func (e *EosClientTest) testing(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error

	err = e.testApi()
	if err != nil {
		logger.Error(err)
		return
	}
}

func NewEosClientTest(ipPort, keosdIpPort string, wg *sync.WaitGroup) (*EosClientTest, error) {
	test := new(EosClientTest)

	cli, err := eos.NewEosClient(ipPort, keosdIpPort)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	test.cli = cli

	go test.testing(wg)

	return test, nil
}
