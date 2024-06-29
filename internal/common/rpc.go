package common

import (
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/time/rate"
)

func NewRPCClient(rpcUrl string) *rpc.Client {

	r := rpc.NewWithCustomRPCClient(rpc.NewWithLimiter(
		rpcUrl,
		rate.Every(time.Second), // time frame
		5,                       // limit of requests per time frame
	))

	return r
}
