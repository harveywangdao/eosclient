package test

import (
	"eosclient/eos"
	"eosclient/logger"
	"sync"
)

type EosClientTest struct {
	EosIpPort string
	coinbase  string
	addrs     []string
}

func (e *EosClientTest) testing(wg *sync.WaitGroup) {
	defer wg.Done()

	cli, err := eos.NewEosClient(e.EosIpPort)
	if err != nil {
		logger.Error(err)
		return
	}

	err = cli.GetInfo()
	if err != nil {
		logger.Error(err)
		return
	}

	err = cli.GetAccount("alice")
	if err != nil {
		logger.Error(err)
		return
	}

	err = cli.GetBlockByNum(1000)
	if err != nil {
		logger.Error(err)
		return
	}

	/*	err = cli.GetProducers()
		if err != nil {
			logger.Error(err)
			return
		}*/

	err = cli.GetBalance("bob", "SYS", "eosio.token")
	if err != nil {
		logger.Error(err)
		return
	}

	/*	err = cli.GetCode("eosio.token")
		if err != nil {
			logger.Error(err)
			return
		}*/

	return

	err = cli.QueryTransaction("0x98593fe3321925c8ef1fb2acdfbd93932a97bbe79654a1bc6d06a8746c7806f0")
	if err != nil {
		logger.Error(err)
		return
	}

	err = cli.GetNewWallet()
	if err != nil {
		logger.Error(err)
		return
	}

	err = cli.Transfer(e.addrs[0], e.addrs[1], "55555")
	if err != nil {
		logger.Error(err)
		return
	}

	err = cli.BlockByNumber(9848)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (e *EosClientTest) setData() {
	e.coinbase = "0xe89d4872b78ab5c5c903583725fe5d485686d6ce"
	e.addrs = append(e.addrs, "0x044b8ab7c603f0938f53e72b7586ec38f3eff044")
	e.addrs = append(e.addrs, "0xced5d036328e9b6da0ed9a7ce1a7770b951fc636")
}

func NewEosClientTest(ipPort string, wg *sync.WaitGroup) (*EosClientTest, error) {
	test := new(EosClientTest)
	test.EosIpPort = ipPort

	test.setData()

	go test.testing(wg)

	return test, nil
}
