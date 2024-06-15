package model

import (
	"math/big"
	"time"
)

type Block struct {
	Height *big.Int `json:"height"`

	Id string `json:"id"`

	Transactions []*Transaction `json:"transactions"`

	CreatedAt time.Time `json:"created_at"`
}

func (b *Block) GenerateTimeToNextBlock() time.Time {
	return time.Now()
}

// Mabye this block is across multiple blocks
func (b *Block) IsAcrossMultipleBlocks() bool {
	return false
}

func (b *Block) CloneWithTransactions(transactions []*Transaction) *Block {
	return &Block{
		Height:       b.Height,
		Id:           b.Id,
		Transactions: transactions,
		CreatedAt:    b.CreatedAt,
	}
}
