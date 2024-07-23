package service

import (
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/chainpusher/blockchain/model"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

type EthereumServiceAssembler struct {
	parsedABI      abi.ABI
	transferMethod *abi.Method
}

func NewEthereumServiceAssembler() (*EthereumServiceAssembler, error) {

	theAbi, err := abi.JSON(strings.NewReader(EthereumUsdtAbi))
	if err != nil {
		return nil, err
	}

	method, err := theAbi.MethodById(EthereumUsdtMethodTransfer)
	if err != nil {
		return nil, err
	}

	return &EthereumServiceAssembler{
		parsedABI:      theAbi,
		transferMethod: method,
	}, nil
}

func (a *EthereumServiceAssembler) ToBlock(block *types.Block) *model.Block {

	createdAt := time.Unix(int64(block.Time()), 0)
	height := block.Number()
	id := block.Hash().String()

	return &model.Block{
		Height:       height,
		Id:           id,
		CreatedAt:    createdAt,
		Transactions: a.BlockToTransactions(block),
	}
}

func (a *EthereumServiceAssembler) ToUSDTTransferArguments(data *[]byte) (*EthereumContractUsdtTransfer, error) {

	var to common.Address
	var amount *big.Int

	b := *data
	if len(b) == 0 {
		return nil, errors.New("empty data")
	}
	args, err := a.transferMethod.Inputs.Unpack(b[4:])
	if err != nil {
		return nil, err
	}

	to = args[0].(common.Address)
	amount = args[1].(*big.Int)

	return &EthereumContractUsdtTransfer{
		To:    &to,
		Value: amount,
	}, nil
}

func (a *EthereumServiceAssembler) ToTransaction(t *types.Transaction) *model.Transaction {

	var crypto model.CryptoCurrency
	var from = ParseEthereumTransactionFromAddress(t)
	var to string
	var amount *big.Int

	// this is an USDT transfer
	txTo := t.To()
	if txTo == nil {
		return nil
	}
	if t.To().String() == EthereumUsdtAddress {
		crypto = model.EthereumUSDT
		data := t.Data()
		transfer, err := a.ToUSDTTransferArguments(&data)
		if err != nil {
			return nil
		}

		amount = transfer.Value
		to = transfer.To.String()
	} else if t.Value().Cmp(big.NewInt(0)) == 1 { // this is a normal transfer
		crypto = model.Ether
		to = t.To().String()
		amount = t.Value()
	} else { // this is a contract call
		return nil
	}

	var createdAt time.Time = t.Time()
	// Time.MarshalJSON: year outside of range [0,9999]
	if createdAt.Year() < 0 || createdAt.Year() > 9999 {
		createdAt = time.Unix(0, 0)
	}

	txId := t.Hash().String()

	return &model.Transaction{
		Id:             txId,
		Platform:       model.PlatformEthereum,
		CryptoCurrency: crypto,
		Payee:          from,
		Payer:          to,
		Amount:         *amount,
		CreatedAt:      t.Time(),
	}
}

func (a *EthereumServiceAssembler) BlockToTransactions(block *types.Block) []*model.Transaction {
	transactions := make([]*model.Transaction, 0)

	for _, tx := range block.Transactions() {
		t := a.ToTransaction(tx)
		if t == nil {
			continue
		}
		transactions = append(transactions, t)
	}

	return transactions

}

func ParseEthereumTransactionFromAddress(t *types.Transaction) string {
	var signer types.Signer
	if t.Type() == types.AccessListTxType {
		signer = types.NewEIP2930Signer(t.ChainId())
	} else if t.Type() == types.DynamicFeeTxType {
		signer = types.NewLondonSigner(t.ChainId())
	} else if t.Type() == types.BlobTxType {
		logrus.Tracef("Blob transaction: %v", t.Hash().String())
		return EthereumEmptyAddress
	} else {
		signer = types.NewEIP155Signer(t.ChainId())
	}

	sender, err := types.Sender(signer, t)

	if err != nil {
		logrus.Error("Failed to get sender: ", err, t.Type())
		return EthereumEmptyAddress
	}
	return sender.String()
}
