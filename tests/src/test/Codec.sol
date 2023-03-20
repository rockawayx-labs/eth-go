// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.10;

import "forge-std/Test.sol";
import "src/Codec.sol";

contract CodecTest is Test {
    Codec codec;

    function setUp() public {
        codec = new Codec();
    }

    function testFunFixedArraySubFixed() public {
        address[2] memory input;
        input[0] = 0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa;
        input[1] = 0xFFfFfFffFFfffFFfFFfFFFFFffFFFffffFfFFFfF;

        bytes memory actual = abi.encodeWithSignature(
            "funFixedArraySubFixed(address[2])",
            input
        );

        codec.logBytes(actual);

        require(
            bytesEquals(
                actual,
                hex"49508494000000000000000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa000000000000000000000000ffffffffffffffffffffffffffffffffffffffff"
            ),
            "Invalid input encode packed bytes"
        );

        (bool success, ) = address(codec).call(actual);
        require(success, "call should have succeeed");
    }

    function testFunFixedArraySubDynamic() public {
        bytes memory bytes_0 = hex"aaaaaaaaaa";
        bytes memory bytes_1 = hex"ffffffffff";

        bytes[2] memory input;
        input[0] = bytes_0;
        input[1] = bytes_1;

        bytes memory actual = abi.encodeWithSignature(
            "funFixedArraySubDynamic(bytes[2])",
            input
        );

        codec.logBytes(actual);

        require(
            bytesEquals(
                actual,
                hex"cd9f57d00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000005aaaaaaaaaa0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005ffffffffff000000000000000000000000000000000000000000000000000000"
            ),
            "Invalid input encode packed bytes"
        );

        (bool success, ) = address(codec).call(actual);
        require(success, "call should have succeeed");
    }

    function testFunBytes8Bytes16Bytes24Bytes32() public {
        bytes8 fixed_bytes_8_type = 0xc5abac1e99944b1d;
        bytes16 fixed_bytes_16_type = 0x57dbc30b9acfebfb86bcc5f9e2fe3fa0;
        bytes24 fixed_bytes_24_type = 0x04a81d8d5c3958b07e558ff8e58e1edf1871c14b34ecdc1c;
        bytes32 fixed_bytes_32_type = 0xf154bf9817019c089414b85e6c5a19fd5d1ea04c103fcd039314132b354ca184;

        bytes memory actual = abi.encodeWithSignature(
            "funBytes8Bytes16Bytes24Bytes32(bytes8,bytes16,bytes24,bytes32)",
            fixed_bytes_8_type,
            fixed_bytes_16_type,
            fixed_bytes_24_type,
            fixed_bytes_32_type
        );

        codec.logBytes(actual);

        require(
            bytesEquals(
                actual,
                hex"680657bdc5abac1e99944b1d00000000000000000000000000000000000000000000000057dbc30b9acfebfb86bcc5f9e2fe3fa00000000000000000000000000000000004a81d8d5c3958b07e558ff8e58e1edf1871c14b34ecdc1c0000000000000000f154bf9817019c089414b85e6c5a19fd5d1ea04c103fcd039314132b354ca184"
            ),
            "Invalid input encode packed bytes"
        );

        (bool success, ) = address(codec).call(actual);
        require(success, "call should have succeeed");
    }

    function testFunInt8() public {
        int8 value = -127;

        bytes memory actual = abi.encodeWithSignature("funInt8(int8)", value);

        codec.logBytes(actual);

        require(
            bytesEquals(
                actual,
                hex"3036e687ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81"
            ),
            "Invalid input encode packed bytes"
        );

        (bool success, ) = address(codec).call(actual);
        require(success, "call should have succeeed");
    }

    function testFunInt32() public {
        int32 value = -898877731;

        bytes memory actual = abi.encodeWithSignature("funInt32(int32)", value);

        codec.logBytes(actual);

        require(
            bytesEquals(
                actual,
                hex"d78caab3ffffffffffffffffffffffffffffffffffffffffffffffffffffffffca6c36dd"
            ),
            "Invalid input encode packed bytes"
        );

        (bool success, ) = address(codec).call(actual);
        require(success, "call should have succeeed");
    }

    function testFunInt256() public {
        int256 value = -9809887317731;

        bytes memory actual = abi.encodeWithSignature(
            "funInt256(int256)",
            value
        );

        codec.logBytes(actual);

        require(
            bytesEquals(
                actual,
                hex"f70af73bfffffffffffffffffffffffffffffffffffffffffffffffffffff713f526b11d"
            ),
            "Invalid input encode packed bytes"
        );

        (bool success, ) = address(codec).call(actual);
        require(success, "call should have succeeed");
    }

    function testFunInt8Int32Int64Int256() public {
        int8 value8 = -127;
        int32 value32 = -898877731;
        int64 value64 = -9809887317731;
        int256 value256 = -223372036854775808;

        bytes memory actual = abi.encodeWithSignature(
            "funInt8Int32Int64Int256(int8,int32,int64,int256)",
            value8,
            value32,
            value64,
            value256
        );

        codec.logBytes(actual);

        require(
            bytesEquals(
                actual,
                hex"db617e8fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81ffffffffffffffffffffffffffffffffffffffffffffffffffffffffca6c36ddfffffffffffffffffffffffffffffffffffffffffffffffffffff713f526b11dfffffffffffffffffffffffffffffffffffffffffffffffffce66c50e2840000"
            ),
            "Invalid input encode packed bytes"
        );

        (bool success, ) = address(codec).call(actual);
        require(success, "call should have succeeed");
    }

    function testFunString() public {
        bytes memory actual = abi.encodeWithSignature(
            "funString(string)",
            "test"
        );

        codec.logBytes(actual);

        require(
            bytesEquals(
                actual,
                hex"b0d94419000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000047465737400000000000000000000000000000000000000000000000000000000"
            ),
            "Invalid input encode packed bytes"
        );

        (bool success, ) = address(codec).call(actual);
        require(success, "call should have succeeed");
    }

    function testFunReturnsString() public {
        bytes memory actual = abi.encodeWithSignature("funReturnsString()");

        codec.logBytes(actual);

        require(
            bytesEquals(actual, hex"7a3719f0"),
            "Invalid input encode packed bytes"
        );

        (bool success, bytes memory output) = address(codec).call(actual);
        require(success, "call should have succeeed");

        codec.logBytes(output);

        require(
            bytesEquals(
                output,
                hex"000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000047465737400000000000000000000000000000000000000000000000000000000"
            ),
            "Invalid output encode packed bytes"
        );
    }

    function testFunReturnsStringString() public {
        bytes memory actual = abi.encodeWithSignature(
            "funReturnsStringString()"
        );

        codec.logBytes(actual);

        require(
            bytesEquals(actual, hex"85032f7c"),
            "Invalid input encode packed bytes"
        );

        (bool success, bytes memory output) = address(codec).call(actual);
        require(success, "call should have succeeed");

        codec.logBytes(output);

        require(
            bytesEquals(
                output,
                hex"000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000005746573743100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000057465737432000000000000000000000000000000000000000000000000000000"
            ),
            "Invalid output encode packed bytes"
        );
    }

    function testFunFixedArrayAddressArrayUint256ReturnsUint256String() public {
        address[2] memory senders;
        senders[0] = 0xFffDB7377345371817F2b4dD490319755F5899eC;
        senders[1] = 0xFFFdb7377345371817F2B4DD490319755F5899EB;

        address[] memory receivers = new address[](3);
        receivers[0] = 0xaFfdb7377345371817f2b4Dd490319755f5899eC;
        receivers[1] = 0xbfFDB7377345371817F2b4dd490319755f5899eC;
        receivers[2] = 0xCffdb7377345371817F2b4Dd490319755F5899EC;

        bytes memory actual = abi.encodeWithSignature(
            "funFixedArrayAddressArrayUint256ReturnsUint256String(address[2],address[])",
            senders,
            receivers
        );

        codec.logBytes(actual);

        require(
            bytesEquals(
                actual,
                hex"b2e8fed2000000000000000000000000fffdb7377345371817f2b4dd490319755f5899ec000000000000000000000000fffdb7377345371817f2b4dd490319755f5899eb00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000003000000000000000000000000affdb7377345371817f2b4dd490319755f5899ec000000000000000000000000bffdb7377345371817f2b4dd490319755f5899ec000000000000000000000000cffdb7377345371817f2b4dd490319755f5899ec"
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

    function testFunAll() public {
        address address_type = 0xFffDB7377345371817F2b4dD490319755F5899eC;
        bytes memory bytes_type = hex"b2";
        bytes8 fixed_bytes_8_type = 0xcf36ac4f97dc10d9;
        bytes32 fixed_bytes_32_type = 0xcf36ac4f97dc10d91fc2cbb20d718e94a8cbfe0f82eaedc6a4aa38946fb797cd;
        int256 fixed_int_256_type = -9809887317731;
        uint256 fixed_uint_256_type = 1827641804;
        bool bool_type = true;
        string memory string_type = "test";
        address[2] memory fixed_array_type;
        address[] memory array_type = new address[](0);

        bytes memory actual = abi.encodeWithSignature(
            "funAll(address,bytes,bytes8,bytes32,int256,uint256,bool,string,address[2],address[])",
            address_type,
            bytes_type,
            fixed_bytes_8_type,
            fixed_bytes_32_type,
            fixed_int_256_type,
            fixed_uint_256_type,
            bool_type,
            string_type,
            fixed_array_type,
            array_type
        );

        codec.logBytes(actual);

        require(
            bytesEquals(
                actual,
                hex"1af93c31000000000000000000000000fffdb7377345371817f2b4dd490319755f5899ec0000000000000000000000000000000000000000000000000000000000000160cf36ac4f97dc10d9000000000000000000000000000000000000000000000000cf36ac4f97dc10d91fc2cbb20d718e94a8cbfe0f82eaedc6a4aa38946fb797cdfffffffffffffffffffffffffffffffffffffffffffffffffffff713f526b11d000000000000000000000000000000000000000000000000000000006cef99cc000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000001a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001e00000000000000000000000000000000000000000000000000000000000000001b200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000474657374000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
            ),
            "Invalid input encode packed bytes"
        );

        (bool success, ) = address(codec).call(actual);
        require(success, "call should have succeeed");
    }

    function testEmitEventIArrayAddress() public {
        vm.recordLogs();

        codec.emitEventIArrayAddress();

        Vm.Log[] memory logs = vm.getRecordedLogs();
        require(logs.length == 1, "Logs length invalid");

        assertEq(logs[0].topics.length, 2);
        assertEq(logs[0].topics[0], keccak256("EventIArrayAddress(address[])"));

        codec.logBytes(logs[0].data);

        require(
            bytesEquals(logs[0].data, hex""),
            "Invalid input encode packed bytes"
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
    function bytesEquals(
        bytes memory self,
        bytes memory other
    ) internal pure returns (bool equal) {
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
