pragma solidity ^0.8.33;

import {HelloWorld} from "../src/HelloWorld.sol";
import {Script} from "forge-std/Script.sol";

contract HelloWorldScript is Script {
    HelloWorld public hw;

    function setUp() public {}

    function run() public {
        vm.startBroadcast();
        hw = new HelloWorld();
        vm.stopBroadcast();
    }
}
