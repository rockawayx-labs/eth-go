pragma solidity 0.8.10;

contract Codec {
    uint256 _nextId;

    event EventUFixedArraySubFixed(address[2] param0);

    event EventUFixedArraySubDynamic(bytes[2] param0);

    event EventUBytes8UBytes16UBytes24UBytes32(
        bytes8 param0,
        bytes16 param1,
        bytes24 param2,
        bytes32 param3
    );

    function funDynamicBoolArray(bool[] calldata param0) public pure {
        assert(param0[0]);
        assert(!param0[1]);
    }

    function funFixedArraySubFixed(address[2] calldata) public pure {}

    function funFixedArraySubDynamic(bytes[2] calldata) public pure {}

    function funBytes8Bytes16Bytes24Bytes32(
        bytes8,
        bytes16,
        bytes24,
        bytes32
    ) public pure {}

    function mint() public returns (uint256 tokenId, uint256 nextId) {
        logMint((tokenId = _nextId++));

        return (tokenId, _nextId);
    }

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

    function funInt128(int128) public pure {}

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

    struct Tuple1 {
        address param0;
    }

    event EventUTuple1(Tuple1 param0);

    function emitEventUTuple1() public {
        Tuple1 memory tuple = Tuple1(
            0xdB0De9288CF0713De91371969efCC9969dd94117
        );

        emit EventUTuple1(tuple);
    }

    struct TupleBool {
        bool param0;
    }

    event EventUTupleBool(TupleBool param0);

    function emitEventUTupleBool() public {
        TupleBool memory tuple = TupleBool(true);

        emit EventUTupleBool(tuple);
    }

    event EventUArrayBool(bool[] param0);

    function emitEventUArrayBool() public {
        bool[] memory bools = new bool[](2);
        bools[0] = true;
        bools[1] = false;

        emit EventUArrayBool(bools);
    }

    event EventUFixedArrayString(string[2] param0);

    function emitEventUFixedArrayString() public {
        string[2] memory strings;
        strings[0] = "first";
        strings[1] = "second";

        emit EventUFixedArrayString(strings);
    }

    event EventUFixedArrayBool(bool[2] param0);

    function emitEventUFixedArrayBool() public {
        bool[2] memory bools;
        bools[0] = true;
        bools[1] = false;

        emit EventUFixedArrayBool(bools);
    }

    event EventMint(uint256 indexed tokenID);

    function logMint(uint256 id) public {
        emit EventMint(id);
    }

    function logBytes(bytes memory data) public pure {
        // Will be recorded automatically by forge, use 'forge test -m <testName> -vvv' to see the results
    }
}
