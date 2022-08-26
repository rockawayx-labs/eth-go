pragma solidity 0.8.10;

contract Codec {
    function fixedArrayAddressArrayUint256ReturnsUint256String(
        address[2] memory,
        address[] memory
    ) public pure returns (uint256, string memory) {
        return (0, "");
    }

    function logBytes(bytes memory data) public pure {
        // Will be recorded automatically by forge, use 'forge test -m <testName> -vvv' to see the results
    }
}
