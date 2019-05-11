package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

var (
	ErrFailedToConvertTransactionValue = fmt.Errorf("Failed to Convert Transaction Value")
)

/*
{
  "status": "1",
  "message": "OK",
  "result": [
    {
      "blockNumber": "4338491",
      "timeStamp": "1507190164",
      "hash": "0x40f1f53edfb3a43c57065aea524448ebf626cc0e9856311172007b54753623ed",
      "nonce": "201",
      "blockHash": "0xf9f953dd153b4be809d89930cf467b0caad09cc942b2c66928d2c4c6dc16568b",
      "from": "0xcbcc5828b0e28d874c1a3c2e5a77d733072bfd67",
      "to": "0x96844a81703a21ffd27a4096baec3f4df1543a51",
      "contractAddress": "0xb8c77482e45f1f44de1745f52c74426c631bdd52",
      "value": "100000000000000000",
      "tokenName": "BNB",
      "tokenSymbol": "BNB",
      "tokenDecimal": "18",
      "transactionIndex": "50",
      "gas": "52514",
      "gasPrice": "1000000000",
      "gasUsed": "37514",
      "cumulativeGasUsed": "1258706",
      "input": "0xa9059cbb00000000000000000000000096844a81703a21ffd27a4096baec3f4df1543a51000000000000000000000000000000000000000000000000016345785d8a0000",
      "confirmations": "3397022"
    }
  ]
}
*/

type Result struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Result  []*Transaction `json:"result"`
}

type Transaction struct {
	BlockNumber       uint64    `json:"blockNumber,omitempty"`
	TimeStamp         time.Time `json:"timeStamp,omitempty"`
	Hash              string    `json:"hash,omitempty"`
	Nonce             uint64    `json:"nonce,omitempty"`
	BlockHash         string    `json:"blockHash,omitempty"`
	From              string    `json:"from,omitempty"`
	To                string    `json:"to,omitempty"`
	ContractAddress   string    `json:"contractAddress,omitempty"`
	Value             *big.Int  `json:"value,omitempty"`
	TokenName         string    `json:"tokenName,omitempty"`
	TokenSymbol       string    `json:"tokenSymbol,omitempty"`
	TokenDecimal      uint64    `json:"tokenDecimal,omitempty"`
	TransactionIndex  uint64    `json:"transactionIndex,omitempty"`
	Gas               uint64    `json:"gas,omitempty"`
	GasPrice          uint64    `json:"gasPrice,omitempty"`
	GasUsed           uint64    `json:"gasUsed,omitempty"`
	CumulativeGasUsed uint64    `json:"cumulativeGasUsed,omitempty"`
	Input             string    `json:"input,omitempty"`
	Confirmations     uint64    `json:"confirmations,omitempty"`
}

func (self *Transaction) UnmarshalJSON(
	bs []byte,
) error {

	tx := &transaction{}
	if err := json.Unmarshal(bs, tx); err != nil {
		return err
	}
	var err error
	if tx.BlockNumber != "" {
		self.BlockNumber, err = strconv.ParseUint(tx.BlockNumber, 10, 64)
		if err != nil {
			return err
		}
	}
	if tx.TimeStamp != "" {
		ts, err := strconv.ParseUint(tx.TimeStamp, 10, 64)
		if err != nil {
			return err
		}
		self.TimeStamp = time.Unix(int64(ts), 0)
	}
	self.Hash = tx.Hash
	if tx.Nonce != "" {
		self.Nonce, err = strconv.ParseUint(tx.Nonce, 10, 64)
		if err != nil {
			return err
		}
	}
	self.BlockHash = tx.BlockHash
	self.From = tx.From
	self.To = tx.To
	self.ContractAddress = tx.ContractAddress
	if tx.Value != "" {
		self.Value = big.NewInt(0)
		_, ok := self.Value.SetString(tx.Value, 10)
		if !ok {
			return ErrFailedToConvertTransactionValue
		}
	}
	self.TokenName = tx.TokenName
	self.TokenSymbol = tx.TokenSymbol
	if tx.TokenDecimal != "" {
		self.TokenDecimal, err = strconv.ParseUint(tx.TokenDecimal, 10, 64)
		if err != nil {
			return err
		}
	}
	if tx.TransactionIndex != "" {
		self.TransactionIndex, err = strconv.ParseUint(tx.TransactionIndex, 10, 64)
		if err != nil {
			return err
		}
	}
	if tx.Gas != "" {
		self.Gas, err = strconv.ParseUint(tx.Gas, 10, 64)
		if err != nil {
			return err
		}
	}
	if tx.GasPrice != "" {
		self.GasPrice, err = strconv.ParseUint(tx.GasPrice, 10, 64)
		if err != nil {
			return err
		}
	}
	if tx.GasUsed != "" {
		self.GasUsed, err = strconv.ParseUint(tx.GasUsed, 10, 64)
		if err != nil {
			return err
		}
	}
	if tx.CumulativeGasUsed != "" {
		self.CumulativeGasUsed, err = strconv.ParseUint(tx.CumulativeGasUsed, 10, 64)
		if err != nil {
			return err
		}
	}
	self.Input = tx.Input
	if tx.Confirmations != "" {
		self.Confirmations, err = strconv.ParseUint(tx.Confirmations, 10, 64)
		if err != nil {
			return err
		}
	}
	return nil
}

type transaction struct {
	BlockNumber       string `json:"blockNumber,omitempty"`
	TimeStamp         string `json:"timeStamp,omitempty"`
	Hash              string `json:"hash,omitempty"`
	Nonce             string `json:"nonce,omitempty"`
	BlockHash         string `json:"blockHash,omitempty"`
	From              string `json:"from,omitempty"`
	To                string `json:"to,omitempty"`
	ContractAddress   string `json:"contractAddress,omitempty"`
	Value             string `json:"value,omitempty"`
	TokenName         string `json:"tokenName,omitempty"`
	TokenSymbol       string `json:"tokenSymbol,omitempty"`
	TokenDecimal      string `json:"tokenDecimal,omitempty"`
	TransactionIndex  string `json:"transactionIndex,omitempty"`
	Gas               string `json:"gas,omitempty"`
	GasPrice          string `json:"gasPrice,omitempty"`
	GasUsed           string `json:"gasUsed,omitempty"`
	CumulativeGasUsed string `json:"cumulativeGasUsed,omitempty"`
	Input             string `json:"input,omitempty"`
	Confirmations     string `json:"confirmations,omitempty"`
}
