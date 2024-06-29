package common

// IdempotencyHeader is the header used to store the idempotency key
const IdempotencyHeader = "Idempotency-Key"

// GenesisHash is a type for accessing multi-environment genesis hashes
type GenesisHash string

const (
	GenesisHashTestnet GenesisHash = "4uhcVJyU9pJkvQyS88uRDiswHXSCkY3zQawwpjk2NsNY"

	GenesisHashDevnet GenesisHash = "EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG"

	GenesisHashMainnetBeta GenesisHash = "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d"
)
