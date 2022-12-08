## Unreleased

- **Breaking** The `ABI` has changed so that multiple events/functions of the name or same id are parsed correctly, in order defined.

- [Changed] Added `out of gas` and `Out of gas` as deterministic error (with the constraint that all provider of `eth_call` used have a `gasCap` configured >= than `gasLimit` used for a call, which should be fixed).

- [Fix] `rpc.Block#BaseFee` is now correctly a `*eth.Uint256` value, you can use `(*uint256.Int)(block.BaseFee).Uint64()` to get back the `uint64` value again (you should check for `nil` value though because it **can** be `nil`).

- [Fix] Encoding of `bytes` in ABI format wasn't properly left padding.

- [Added] `LogEventDef.LogID()` is now exposed publicly.

- [Fix] `LogEventDef.Signature()` is now formatted to be read for `Keccak` processing.

- [Fix] `rpc.Block.Nonce` is now a FixedUint64 to enforce `0x0000000000000000` encoding.

- [Fix] `rpc.Block.Timestamp` is now encoded as a `uint64` instead of a time.RFC3339 string

- [Fix] `rpc.BlockRef` decoding fixed to support either a `BlockNumber` or a `BlockHash`