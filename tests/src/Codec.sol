pragma solidity 0.8.10;

contract Codec {
    event EventUFixedArraySubFixed(address[2] param0);

    event EventUFixedArraySubDynamic(bytes[2] param0);

    event EventUBytes8UBytes16UBytes24UBytes32(
        bytes8 param0,
        bytes16 param1,
        bytes24 param2,
        bytes32 param3
    );

    function funFixedArraySubFixed(address[2] calldata) public pure {}

    function funFixedArraySubDynamic(bytes[2] calldata) public pure {}

    function funBytes8Bytes16Bytes24Bytes32(
        bytes8,
        bytes16,
        bytes24,
        bytes32
    ) public pure {}

    function funReturnsString() public pure returns (string memory) {
        return "test";
    }

    function funReturnsStringString()
        public
        pure
        returns (string memory, string memory)
    {
        return ("test1", "test2");
    }

    function funInt8(int8) public pure {}

    function funInt32(int32) public pure {}

    function funInt256(int256) public pure {}

    function funInt8Int32Int64Int256(int8, int32, int64, int256) public pure {}

    function funString(string memory) public pure {}

    function funFixedArrayAddressArrayUint256ReturnsUint256String(
        address[2] memory,
        address[] memory
    ) public pure returns (uint256, string memory) {
        return (0, "");
    }

    function funAll(
        address,
        bytes memory,
        bytes8,
        bytes32,
        int256,
        uint256,
        bool,
        string memory,
        address[2] memory,
        address[] memory
    ) public pure {}

    event EventIArrayAddress(address[] indexed param0);

    function emitEventIArrayAddress() public {
        address[] memory addresses = new address[](2);
        addresses[0] = 0xdB0De9288CF0713De91371969efCC9969dd94117;
        addresses[1] = 0xdB0De9288CF0713De91371969efCC9969dd94117;

        emit EventIArrayAddress(addresses);
    }

    function logBytes(bytes memory data) public pure {
        // Will be recorded automatically by forge, use 'forge test -m <testName> -vvv' to see the results
    }
}
