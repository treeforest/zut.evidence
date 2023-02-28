package main

import (
	"encoding/json"
	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
	log "github.com/treeforest/logger"
	evidence "github.com/treeforest/zut.evidence/blockchain/contracts/evidence"
	"io/ioutil"
)

type ContractAddress struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func NewContractAddress(name, address string) ContractAddress {
	log.Infof("contract name: %s", name)
	log.Infof("contract address: %s", address)
	return ContractAddress{Name: name, Address: address}
}

type ContractInfo struct {
	Addresses []ContractAddress `json:"addresses"`
}

func main() {
	configs, err := conf.ParseConfigFile("config.toml")
	if err != nil {
		log.Fatal("parse config error: ", err)
	}
	config := &configs[0]

	client, err := client.Dial(config)
	if err != nil {
		log.Fatal("dial fisco bcos node error: ", err)
	}

	// deploy evidence
	evidenceAddress, _, _, err := evidence.DeployEvidence(client.GetTransactOpts(), client) // deploy contract
	if err != nil {
		log.Fatal("deploy evidence error: ", err)
	}

	info := &ContractInfo{Addresses: make([]ContractAddress, 0)}
	info.Addresses = []ContractAddress{
		NewContractAddress("evidence", evidenceAddress.Hex()),
	}

	data, _ := json.MarshalIndent(info, "", "\t")
	err = ioutil.WriteFile("address.json", data, 777)
	if err != nil {
		log.Fatal(err)
	}
}
