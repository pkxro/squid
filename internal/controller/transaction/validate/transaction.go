package validate

import (
	"context"
	"errors"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func ValidateTransaction(ctx context.Context, rpcc *rpc.Client, transaction *solana.Transaction, feePayer solana.PrivateKey, lamportsPerSignature uint64, maxSignatures uint64) (*solana.Signature, error) {
	if !transaction.Message.AccountKeys.First().Equals(feePayer.PublicKey()) {
		return nil, errors.New("bad fee payer")
	}

	if transaction.Message.RecentBlockhash.IsZero() {
		return nil, errors.New("missing blockhash")
	}

	f, err := rpcc.GetFeeForMessage(ctx, transaction.Message.RecentBlockhash.String(), rpc.CommitmentConfirmed)
	if err != nil {
		return nil, errors.New("error fetching fee")
	}

	if *f.Value > lamportsPerSignature {
		return nil, errors.New("fee too high")
	}

	if len(transaction.Signatures) == 0 {
		return nil, errors.New("no sigs")
	}

	if len(transaction.Signatures) > int(maxSignatures) {
		return nil, errors.New("too many signatures")
	}

	signatures, _ := transaction.PartialSign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(feePayer.PublicKey()) {
			return &feePayer
		}
		return nil
	})

	if err := transaction.VerifySignatures(); err != nil {
		return nil, errors.New("error verifying transaction signatures")
	}

	return &signatures[0], nil
}
