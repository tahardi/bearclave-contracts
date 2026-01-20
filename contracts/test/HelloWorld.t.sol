pragma solidity ^0.8.26;

import {Test} from "forge-std/Test.sol";
import {HelloWorld} from "../src/HelloWorld.sol";

contract HelloWorldTest is Test {
    HelloWorld public hw;

    function setUp() public {
        hw = new HelloWorld();
    }

    function test_greet() public view {
        assertEq(hw.greet(), "Hello, World!");
    }
}
