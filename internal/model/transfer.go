package model

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type TransferDetails struct {
	Source      solana.PublicKey
	Destination solana.PublicKey
	Owner       solana.PublicKey
	Mint        solana.PublicKey
	Decimals    uint8
	IsSystem    bool
	Amount      uint64
	Account     *rpc.Account
}
