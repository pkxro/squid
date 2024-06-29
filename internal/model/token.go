package model

import "github.com/gagliardetto/solana-go"

type Token struct {
	Mint     solana.PublicKey
	Account  solana.PublicKey
	Fee      uint64
	Decimals uint8
}
