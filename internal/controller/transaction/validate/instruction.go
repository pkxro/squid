package validate

import (
	"github.com/gagliardetto/solana-go"
)

// ValidateInstructions prevents draining by ensuring that the fee payer isn't provided as writable or a signer to any instruction.
// Returns an error if the transaction contains instructions that could potentially drain the fee payer.
func ValidateInstructions(transaction *solana.Transaction) error {
	for _, instruction := range transaction.Message.Instructions {
		for _, accountMeta := range instruction.Accounts {
			account := transaction.Message.AccountKeys[accountMeta]
			isWritable, err := transaction.Message.IsWritable(account)
			if err != nil {
				return err
			}
			if isWritable || transaction.Message.IsSigner(account) {
				return err
			}
		}
	}
	return nil
}
