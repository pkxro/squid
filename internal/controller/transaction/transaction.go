package transaction

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkxro/squid/internal/cache"
	"github.com/pkxro/squid/internal/config"
	"github.com/pkxro/squid/internal/controller/transaction/validate"
	"github.com/pkxro/squid/internal/model"
	"go.uber.org/zap"
)

// TransactionManager handles all transaction-related operations
type TransactionManager struct {
	logger *zap.Logger
	rpcc   *rpc.Client
	cache  *cache.Manager
	pKey   solana.PrivateKey
	ocfg   *config.OperationConfig
}

// TransactionConfig holds configuration specific to transactions
type TransactionConfig struct {
	FeePayer             solana.PrivateKey
	MaxSignatures        uint64
	LamportsPerSignature uint64
	SameSourceTimeout    time.Duration
}

// NewTransactionManager creates a new TransactionManager
func NewTransactionManager(
	logger *zap.Logger,
	rpcc *rpc.Client,
	cache *cache.Manager,
	pKey solana.PrivateKey,
	ocfg *config.OperationConfig,
) *TransactionManager {
	return &TransactionManager{
		logger: logger,
		rpcc:   rpcc,
		cache:  cache,
		pKey:   pKey,
		ocfg:   ocfg,
	}
}

// SignWithTokenFee signs a transaction with a token fee
func (tm *TransactionManager) SignWithTokenFee(ctx context.Context, req model.SignWithTokenFeeRequest) (*model.SignWithTokenFeeResponse, error) {
	// Prevent simple duplicate transactions using a hash of the message
	messageBytes, err := req.Transaction.Message.MarshalBinary()
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(messageBytes)
	key := base64.StdEncoding.EncodeToString(hash[:])

	ok, err := tm.cache.Client.Exists(ctx, key, "")
	if err != nil {
		return nil, err
	}

	if ok {
		return nil, errors.New("duplicate tx")
	}

	tm.cache.Client.Set(ctx, key, "", hash, time.Minute*time.Duration(tm.ocfg.CacheLimitMinutes))

	sig, err := validate.ValidateTransaction(
		ctx, tm.rpcc, req.Transaction, tm.pKey, req.LamportsPerSignature, req.MaxSignatures,
	)
	if err != nil {
		return nil, err
	}

	err = validate.ValidateInstructions(req.Transaction)
	if err != nil {
		return nil, err
	}

	err = validate.ValidateTransfer(ctx, tm.rpcc, req.Transaction, tm.ocfg.AllowedTokens)
	if err != nil {
		return nil, err
	}

	// // Implement lockout for the source token account
	// key = "transfer/lastSignature/" + transfer.Keys.Source.PublicKey.String()
	// lastSignatureInterface, err := params.Cache.Get(params.Ctx, key)
	// if err == nil {
	// 	lastSignature, ok := lastSignatureInterface.(int64)
	// 	if ok && time.Now().UnixNano()-lastSignature < params.SameSourceTimeout.Nanoseconds() {
	// 		return nil, errors.New("duplicate transfer")
	// 	}
	// }

	// err = params.Cache.Set(params.Ctx, key, time.Now().UnixNano())
	// if err != nil {
	// 	return nil, err
	// }

	_, err = tm.rpcc.SimulateRawTransactionWithOpts(ctx, messageBytes, &rpc.SimulateTransactionOpts{
		Commitment: rpc.CommitmentConfirmed,
		SigVerify:  false,
	})
	if err != nil {
		return nil, err
	}

	return &model.SignWithTokenFeeResponse{Signature: sig.String()}, nil
}
