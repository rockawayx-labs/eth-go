// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.10;

contract PersonalSigning {
    function recoverPersonalSigner(bytes32 messageHash, bytes memory signature)
        public
        pure
        returns (address)
    {
        return recoverSigner(signedMessageHash(messageHash), signature);
    }

    function signedMessageHash(bytes32 messageHash)
        public
        pure
        returns (bytes32)
    {
        return keccak256(signedMessageHashPayload(messageHash));
    }

    function signedMessageHashPayload(bytes32 messageHash)
        public
        pure
        returns (bytes memory)
    {
        return
            abi.encodePacked("\x19Ethereum Signed Message:\n32", messageHash);
    }

    function recoverSigner(
        bytes32 signedMessageHashValue,
        bytes memory signature
    ) public pure returns (address) {
        (bytes32 r, bytes32 s, uint8 v) = splitSignature(signature);

        return ecrecover(signedMessageHashValue, v, r, s);
    }

    function splitSignature(bytes memory sig)
        public
        pure
        returns (
            bytes32 r,
            bytes32 s,
            uint8 v
        )
    {
        require(sig.length == 65, "invalid signature length");

        assembly {
            /*
            First 32 bytes stores the length of the signature

            add(sig, 32) = pointer of sig + 32
            effectively, skips first 32 bytes of signature

            mload(p) loads next 32 bytes starting at the memory address p into memory
            */

            // first 32 bytes, after the length prefix
            r := mload(add(sig, 32))
            // second 32 bytes
            s := mload(add(sig, 64))
            // final byte (first byte of the next 32 bytes)
            v := byte(0, mload(add(sig, 96)))
        }

        // implicitly return (r, s, v)
    }
}
