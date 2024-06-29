package common

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// ValidateInstructions prevents draining by ensuring that the fee payer isn't provided as writable or a signer to any instruction.
// Returns an error if the transaction contains instructions that could potentially drain the fee payer.
func ValidateInstructions(transaction *solana.Transaction, feePayer solana.PublicKey) error {
	for _, instruction := range transaction.Message.Instructions {
		for _, accountIdx := range instruction.Accounts {
			if accountIdx >= uint16(len(transaction.Message.AccountKeys)) {
				return fmt.Errorf("invalid account index in instruction")
			}
			account := transaction.Message.AccountKeys[accountIdx]
			isWritable, err := transaction.IsWritable(account)
			if err != nil {
				return err
			}
			isSigner := transaction.IsSigner(account)

			if (isWritable || isSigner) && account.Equals(feePayer) {
				return fmt.Errorf("invalid account: fee payer (%s) is writable or signer in instruction", feePayer.String())
			}
		}
	}
	return nil
}
