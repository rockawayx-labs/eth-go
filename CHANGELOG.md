## Unreleased

- [Fix] `rpc.Block#BaseFee` is now correctly a `*eth.Uint256` value, you can use `(*uint256.Int)(block.BaseFee).Uint64()` to get back the `uint64` value again (you should check for `nil` value though because it **can** be `nil`).

- [Fix] Encoding of `bytes` in ABI format wasn't properly left padding.

- [Added] `LogEventDef.LogID()` is now exposed publicly.

- [Fix] `LogEventDef.Signature()` is now formatted to be read for `Keccak` processing.
