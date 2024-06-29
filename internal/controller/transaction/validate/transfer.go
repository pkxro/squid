package validate

import (
	"context"
	"errors"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkxro/squid/internal/model"
)

func ValidateTransfer(
	ctx context.Context,
	rpcc *rpc.Client,
	transaction *solana.Transaction,
	allowedTokens []model.Token,
) error {

	if len(transaction.Message.Instructions) == 0 {
		return errors.New("missing instructions")
	}

	// Get the first instruction of the transaction
	firstInstruction := transaction.Message.Instructions[0]

	// Accounts
	accounts, err := firstInstruction.ResolveInstructionAccounts(&transaction.Message)
	if err != nil {
		return err
	}

	ix, err := system.DecodeInstruction(accounts, firstInstruction.Data)
	if err != nil {
		return err
	}

	if sysXfer, ok := ix.Impl.(*system.Transfer); ok {
		account, err := rpcc.GetAccountInfo(ctx, sysXfer.GetFundingAccount().PublicKey)
		if err != nil {
			return err
		}

		return validateSystemTransfer(ctx, rpcc, transaction, sysXfer, account.Value, allowedTokens)
	}

	decTok, err := token.DecodeInstruction(accounts, firstInstruction.Data)
	if err != nil {
		return err
	}

	tokXfer := decTok.Impl.(*token.Transfer)

	account, err := rpcc.GetAccountInfo(ctx, tokXfer.GetSourceAccount().PublicKey)
	if err != nil {
		return err
	}

	err = validateTokenTransfer(ctx, rpcc, transaction, tokXfer, account.Value, allowedTokens)
	if err != nil {
		return nil
	}
	err = validateTokenTransferChecked(ctx, rpcc, decTok, allowedTokens)
	if err != nil {
		return nil
	}

	return nil
}

func validateSystemTransfer(ctx context.Context, rpcc *rpc.Client, transaction *solana.Transaction, xfer *system.Transfer, account *rpc.Account, allowedTokens []model.Token) error {
	// source := xfer.GetFundingAccount().PublicKey
	// TODO

	return nil
}

func validateTokenTransferChecked(ctx context.Context, rpcc *rpc.Client, decTok *token.Instruction, allowedTokens []model.Token) error {
	xferChecked := decTok.Impl.(*token.TransferChecked)
	if xferChecked == nil {
		return nil
	}

	// Check if the mint of the source account is in the allowed tokens list
	var vToken *model.Token
	for _, token := range allowedTokens {
		if token.Mint.Equals(xferChecked.GetMintAccount().PublicKey) {
			vToken = &token
			break
		}
	}
	if vToken == nil {
		return errors.New("mint not in allowed tokens list")
	}

	if xferChecked.Decimals != &vToken.Decimals {
		return errors.New("invalid decimals")
	}
	if xferChecked.GetMintAccount().PublicKey != vToken.Mint {
		return errors.New("invalid mint")
	}
	if xferChecked.GetMintAccount().IsWritable {
		return errors.New("mint is writable")
	}
	if xferChecked.GetMintAccount().IsSigner {
		return errors.New("mint is signer")
	}

	return nil

}

func validateTokenTransfer(ctx context.Context, rpcc *rpc.Client, transaction *solana.Transaction, xfer *token.Transfer, account *rpc.Account, allowedTokens []model.Token) error {
	amt := *xfer.Amount
	sourceAcct := xfer.GetSourceAccount()
	destAcct := xfer.GetDestinationAccount()
	ownerAcct := xfer.GetOwnerAccount()

	// Check if the mint of the source account is in the allowed tokens list
	// transfer doesn't have a min getter, so try on the 2nd pubkey
	var vToken *model.Token
	for _, token := range allowedTokens {
		if token.Mint.Equals(xfer.Accounts[1].PublicKey) {
			vToken = &token
			break
		}
	}
	if vToken == nil {
		return errors.New("mint not in allowed tokens list")
	}

	acct, err := rpcc.GetAccountInfo(ctx, sourceAcct.PublicKey)
	if err != nil {
		return err
	}

	if acct.Value.Executable {
		return errors.New("cannot be executable")
	}

	// Check that the instruction is going to pay the fee
	if amt < vToken.Fee {
		return errors.New("invalid amount")
	}

	if acct.Value.Owner != ownerAcct.PublicKey {
		return errors.New("invalid owner")
	}

	if acct.Value.Lamports < amt {
		return errors.New("insufficient balance of source account")
	}

	if *xfer.Amount < vToken.Fee {
		return errors.New("invalid amount")
	}

	if sourceAcct.IsSigner {
		return errors.New("source is signer")
	}

	if sourceAcct.IsWritable {
		return errors.New("source is not writable")
	}

	if destAcct.PublicKey == vToken.Account {
		return errors.New("invalid destination")
	}

	if !destAcct.IsWritable {
		return errors.New("destination not writable")
	}

	if destAcct.IsSigner {
		return errors.New("destination is signer")
	}

	return nil
}
