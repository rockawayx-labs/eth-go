// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.10;

import "ds-test/test.sol";
import "src/Codec.sol";

contract CodecTest is DSTest {
    Codec codec;

    function setUp() public {
        codec = new Codec();
    }

    function testFixedArrayAddressArrayUint256ReturnsUint256String() public {
        address[2] memory senders;
        senders[0] = 0xFffDB7377345371817F2b4dD490319755F5899eC;
        senders[1] = 0xFFFdb7377345371817F2B4DD490319755F5899EB;

        address[] memory receivers = new address[](3);
        receivers[0] = 0xaFfdb7377345371817f2b4Dd490319755f5899eC;
        receivers[1] = 0xbfFDB7377345371817F2b4dd490319755f5899eC;
        receivers[2] = 0xCffdb7377345371817F2b4Dd490319755F5899EC;

        bytes memory actual = abi.encodeWithSignature(
            "fixedArrayAddressArrayUint256ReturnsUint256String(address[2],address[])",
            senders,
            receivers
        );

        codec.logBytes(actual);

        require(
            bytesEquals(
                actual,
                hex"74ac01d1000000000000000000000000fffdb7377345371817f2b4dd490319755f5899ec000000000000000000000000fffdb7377345371817f2b4dd490319755f5899eb00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000003000000000000000000000000affdb7377345371817f2b4dd490319755f5899ec000000000000000000000000bffdb7377345371817f2b4dd490319755f5899ec000000000000000000000000cffdb7377345371817f2b4dd490319755f5899ec"
            ),
            "Invalid input encode packed bytes"
        );

        (bool success, bytes memory output) = address(codec).call(actual);
        require(success, "call should have succeeed");

        codec.logBytes(output);

        require(
            bytesEquals(
                output,
                hex"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000"
            ),
            "Invalid output encode packed bytes"
        );
    }

    // FIXME: Shared for all tests ..., copied from test/PersonalSigning
    // Compares the 'len' bytes starting at address 'addr' in memory with the 'len'
    // bytes starting at 'addr2'.
    // Returns 'true' if the bytes are the same, otherwise 'false'.
    function memoryEquals(
        uint256 addr,
        uint256 addr2,
        uint256 len
    ) internal pure returns (bool equal) {
        assembly {
            equal := eq(keccak256(addr, len), keccak256(addr2, len))
        }
    }

    // Checks if two `bytes memory` variables are equal. This is done using hashing,
    // which is much more gas efficient then comparing each byte individually.
    // Equality means that:
    //  - 'self.length == other.length'
    //  - For 'n' in '[0, self.length)', 'self[n] == other[n]'
    function bytesEquals(bytes memory self, bytes memory other)
        internal
        pure
        returns (bool equal)
    {
        if (self.length != other.length) {
            return false;
        }
        uint256 addr;
        uint256 addr2;

        assembly {
            addr := add(
                self,
                /*BYTES_HEADER_SIZE*/
                32
            )
            addr2 := add(
                other,
                /*BYTES_HEADER_SIZE*/
                32
            )
        }

        equal = memoryEquals(addr, addr2, self.length);
    }
}
