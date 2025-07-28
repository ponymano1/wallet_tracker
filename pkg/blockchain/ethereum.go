package blockchain

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ERC20 ABI (简化版)
const ERC20ABI = `[
    {
        "constant":true,
        "inputs":[{"name":"_owner","type":"address"}],
        "name":"balanceOf",
        "outputs":[{"name":"balance","type":"uint256"}],
        "type":"function"
    },
    {
        "constant":true,
        "inputs":[],
        "name":"decimals",
        "outputs":[{"name":"","type":"uint8"}],
        "type":"function"
    },
    {
        "constant":true,
        "inputs":[],
        "name":"symbol",
        "outputs":[{"name":"","type":"string"}],
        "type":"function"
    },
    {
        "constant":true,
        "inputs":[],
        "name":"name",
        "outputs":[{"name":"","type":"string"}],
        "type":"function"
    }
]`

type BlockchainClient struct {
	client *ethclient.Client
	abi    abi.ABI
}

func NewBlockchainClient(rpcURL string) (*BlockchainClient, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	contractABI, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return nil, err
	}

	return &BlockchainClient{
		client: client,
		abi:    contractABI,
	}, nil
}

func (bc *BlockchainClient) GetTokenBalance(tokenAddress, walletAddress string) (*big.Int, error) {
	tokenAddr := common.HexToAddress(tokenAddress)
	walletAddr := common.HexToAddress(walletAddress)

	data, err := bc.abi.Pack("balanceOf", walletAddr)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &tokenAddr,
		Data: data,
	}

	result, err := bc.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	balance := new(big.Int)
	if err := bc.abi.UnpackIntoInterface(&balance, "balanceOf", result); err != nil {
		return nil, err
	}

	return balance, nil
}

func (bc *BlockchainClient) GetTokenInfo(tokenAddress string) (string, string, uint8, error) {
	tokenAddr := common.HexToAddress(tokenAddress)

	// 获取 symbol
	symbolData, err := bc.abi.Pack("symbol")
	if err != nil {
		return "", "", 0, err
	}

	symbolMsg := ethereum.CallMsg{
		To:   &tokenAddr,
		Data: symbolData,
	}

	symbolResult, err := bc.client.CallContract(context.Background(), symbolMsg, nil)
	if err != nil {
		return "", "", 0, err
	}

	var symbol string
	if err := bc.abi.UnpackIntoInterface(&symbol, "symbol", symbolResult); err != nil {
		return "", "", 0, err
	}

	// 获取 name
	nameData, err := bc.abi.Pack("name")
	if err != nil {
		return "", "", 0, err
	}

	nameMsg := ethereum.CallMsg{
		To:   &tokenAddr,
		Data: nameData,
	}

	nameResult, err := bc.client.CallContract(context.Background(), nameMsg, nil)
	if err != nil {
		return "", "", 0, err
	}

	var name string
	if err := bc.abi.UnpackIntoInterface(&name, "name", nameResult); err != nil {
		return "", "", 0, err
	}

	// 获取 decimals
	decimalsData, err := bc.abi.Pack("decimals")
	if err != nil {
		return "", "", 0, err
	}

	decimalsMsg := ethereum.CallMsg{
		To:   &tokenAddr,
		Data: decimalsData,
	}

	decimalsResult, err := bc.client.CallContract(context.Background(), decimalsMsg, nil)
	if err != nil {
		return "", "", 0, err
	}

	var decimals uint8
	if err := bc.abi.UnpackIntoInterface(&decimals, "decimals", decimalsResult); err != nil {
		return "", "", 0, err
	}

	return symbol, name, decimals, nil
}

func (bc *BlockchainClient) Close() {
	bc.client.Close()
}
