package model

import (
	"time"

	"github.com/gagliardetto/solana-go"
)

type SignWithTokenFeeRequest struct {
	Transaction          *solana.Transaction
	MaxSignatures        uint64
	LamportsPerSignature uint64
	SameSourceTimeout    time.Duration
}

type SignWithTokenFeeResponse struct {
	Signature string `json:"signature"`
}
