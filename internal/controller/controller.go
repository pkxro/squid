package controller

import (
	"github.com/gagliardetto/solana-go"
	"github.com/pkxro/squid/internal/cache"
	"github.com/pkxro/squid/internal/common"
	"github.com/pkxro/squid/internal/config"
	"github.com/pkxro/squid/internal/controller/transaction"
	"go.uber.org/zap"
)

// Manager is the struct that the constructor implements
type Manager struct {
	Scfg   *config.SecretConfig
	Ocfg   *config.OperationConfig
	Logger *zap.Logger
	Tx     *transaction.TransactionManager
}

// NewController initializes a Manager
func NewControllerManager(
	logger *zap.Logger,
	scfg *config.SecretConfig,
	ocfg *config.OperationConfig,
	cache *cache.Manager,

) Manager {
	rpcc := common.NewRPCClient(scfg.RpcUrl)

	// todo, don't panic here, panic earlier
	keyFromSecret := solana.MustPrivateKeyFromBase58(scfg.FeePayer)

	tx := transaction.NewTransactionManager(logger, rpcc, cache, keyFromSecret, *ocfg)

	return Manager{
		Scfg:   scfg,
		Ocfg:   ocfg,
		Logger: logger,
		Tx:     tx,
	}
}
