// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.10;

import "ds-test/test.sol";
import "src/PersonalSigning.sol";

contract PersonalSigningTest is DSTest {
    PersonalSigning signing;

    function setUp() public {
        signing = new PersonalSigning();
    }

    function testRecoverPersonalSigner() public {
        bytes
            memory signature = hex"cfc0c9160c1fbe884f02298b76e194611904a4f18b6814dfd22f95b659b2f90b2564e455543316f125c3c51ec72b22917401d4872ddabb9a7a2d0bea1a827db31b";
        address signer = signing.recoverPersonalSigner(
            0xcf36ac4f97dc10d91fc2cbb20d718e94a8cbfe0f82eaedc6a4aa38946fb797cd,
            signature
        );

        assertEq(signer, 0xFffDB7377345371817F2b4dD490319755F5899eC);
    }

    function testSignedMessageHash() public {
        bytes32 actual = signing.signedMessageHash(
            0xcf36ac4f97dc10d91fc2cbb20d718e94a8cbfe0f82eaedc6a4aa38946fb797cd
        );

        assertEq(
            actual,
            0x58749f0b9677f513b6cf2a4e163dc7a09d61d6e4168e05b25fd11a4ffd62944c
        );
    }

    function testSignedMessageHashPayload() public view {
        bytes memory actual = signing.signedMessageHashPayload(
            0xcf36ac4f97dc10d91fc2cbb20d718e94a8cbfe0f82eaedc6a4aa38946fb797cd
        );

        require(
            bytesEquals(
                actual,
                hex"19457468657265756d205369676e6564204d6573736167653a0a3332cf36ac4f97dc10d91fc2cbb20d718e94a8cbfe0f82eaedc6a4aa38946fb797cd"
            ),
            "Invalid encode packaed bytes"
        );
    }

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
