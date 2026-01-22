pragma solidity ^0.8.33;

import {BearCoin} from "../src/BearCoin.sol";
import {Script} from "forge-std/Script.sol";

contract BearCoinScript is Script {
    BearCoin public bcn;

    function setUp() public {}

    function run() public {
        vm.startBroadcast();
        bcn = new BearCoin();
        vm.stopBroadcast();
    }
}