pragma solidity ^0.8.33;

import {Test} from "forge-std/Test.sol";
import {BearCoin} from "../src/BearCoin.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Errors} from "@openzeppelin/contracts/interfaces/draft-IERC6093.sol";

contract BearCoinTest is Test {
    BearCoin public bcn;
    address public owner;
    address public alice;
    address public bob;

    function setUp() public {
        owner = address(1);
        alice = address(2);
        bob = address(3);

        vm.prank(owner);
        bcn = new BearCoin();
    }

    function test_BearCoin() public view {
        assertEq(bcn.name(), "BearCoin");
        assertEq(bcn.symbol(), "BCN");
        assertEq(bcn.decimals(), 18);
        assertEq(bcn.totalSupply(), 1000000 * 10 ** 18);

        assertEq(bcn.owner(), owner);
        assertEq(bcn.balanceOf(owner), 1000000 * 10 ** 18);

        assertEq(bcn.balanceOf(alice), 0);
        assertEq(bcn.balanceOf(bob), 0);
    }

    function test_approve() public {
        // given
        uint256 approveAmount = 500 * 10 ** bcn.decimals();
        assertEq(bcn.allowance(owner, alice), 0);

        vm.prank(owner);
        vm.expectEmit(true, true, false, true);
        emit IERC20.Approval(owner, alice, approveAmount);

        // when
        bcn.approve(alice, approveAmount);

        // then
        assertEq(bcn.allowance(owner, alice), approveAmount);
    }

    function test_burn() public {
        // given
        uint256 burnAmount = 100 * 10 ** bcn.decimals();
        uint256 initialBalance = bcn.balanceOf(owner);
        uint256 initialSupply = bcn.totalSupply();

        vm.prank(owner);
        vm.expectEmit(true, false, false, true);
        emit BearCoin.Burn(owner, burnAmount);

        // when
        bcn.burn(burnAmount);

        // then
        assertEq(bcn.balanceOf(owner), initialBalance - burnAmount);
        assertEq(bcn.totalSupply(), initialSupply - burnAmount);
    }

    function test_burn_non_owner() public {
        // given
        uint256 burnAmount = 50 * 10 ** bcn.decimals();
        uint256 transferAmount = 100 * 10 ** bcn.decimals();
        uint256 initialSupply = bcn.totalSupply();

        vm.prank(owner);
        vm.expectEmit(true, true, false, true);
        emit IERC20.Transfer(owner, alice, transferAmount);

        require(bcn.transfer(alice, transferAmount), "Transfer failed");
        assertEq(bcn.balanceOf(alice), transferAmount);

        // when
        vm.prank(alice);
        vm.expectEmit(true, false, false, true);
        emit BearCoin.Burn(alice, burnAmount);

        bcn.burn(burnAmount);

        // then
        assertEq(bcn.balanceOf(alice), transferAmount - burnAmount);
        assertEq(bcn.totalSupply(), initialSupply - burnAmount);
    }

    function test_burn_revert_insufficient_balance() public {
        // given
        uint256 burnAmount = 100 * 10 ** bcn.decimals();
        uint256 balance = bcn.balanceOf(alice);
        assertEq(balance, 0);

        vm.prank(alice);
        vm.expectRevert(
            abi.encodeWithSelector(IERC20Errors.ERC20InsufficientBalance.selector, alice, balance, burnAmount)
        );

        // when/then
        bcn.burn(burnAmount);
    }

    function test_mint() public {
        // given
        uint256 burnAmount = 100 * 10 ** bcn.decimals();
        uint256 initialBalance = bcn.balanceOf(owner);
        uint256 initialSupply = bcn.totalSupply();

        vm.prank(owner);
        vm.expectEmit(true, false, false, true);
        emit BearCoin.Burn(owner, burnAmount);

        bcn.burn(burnAmount);

        assertEq(bcn.balanceOf(owner), initialBalance - burnAmount);
        assertEq(bcn.totalSupply(), initialSupply - burnAmount);

        // when
        vm.prank(owner);
        vm.expectEmit(true, false, false, true);
        emit BearCoin.Mint(alice, burnAmount);

        bcn.mint(alice, burnAmount);

        // then
        assertEq(bcn.balanceOf(alice), burnAmount);
        assertEq(bcn.totalSupply(), initialSupply);
    }

    function test_mint_revert_not_owner() public {
        // given
        uint256 mintAmount = 100 * 10 ** bcn.decimals();
        vm.expectRevert("Not owner");

        // when/then
        vm.prank(alice);
        bcn.mint(alice, mintAmount);
    }

    function test_mint_revert_exceeds_total_supply() public {
        // given
        uint256 exceedAmount = bcn.TOTAL_SUPPLY() + 1;
        vm.prank(owner);
        vm.expectRevert("Minting exceeds total supply");

        // when/then
        bcn.mint(alice, exceedAmount);
    }

    function test_transfer() public {
        // given
        uint256 transferAmount = 100 * 10 ** 18;

        vm.prank(owner);
        vm.expectEmit(true, true, false, true);
        emit IERC20.Transfer(owner, alice, transferAmount);

        // when
        // https://getfoundry.sh/forge/linting/#erc20-unchecked-transfer
        require(bcn.transfer(alice, transferAmount), "Transfer failed");

        // then
        assertEq(bcn.balanceOf(alice), transferAmount);
        assertEq(bcn.balanceOf(owner), bcn.totalSupply() - transferAmount);
    }

    function test_transfer_to_self() public {
        // given
        uint256 transferAmount = 100 * 10 ** 18;
        uint256 initialBalance = bcn.balanceOf(owner);

        vm.prank(owner);
        vm.expectEmit(true, true, false, true);
        emit IERC20.Transfer(owner, owner, transferAmount);

        // when
        require(bcn.transfer(owner, transferAmount), "Transfer failed");

        // then
        assertEq(bcn.balanceOf(owner), initialBalance);
    }

    function test_transfer_revert_insufficient_balance() public {
        // given
        uint256 balance = bcn.balanceOf(alice);
        uint256 transferAmount = 100 * 10 ** 18;

        vm.prank(alice);
        vm.expectRevert(
            abi.encodeWithSelector(IERC20Errors.ERC20InsufficientBalance.selector, alice, balance, transferAmount)
        );

        // when/then
        require(!bcn.transfer(bob, transferAmount), "Transfer succeeded");
    }

    function test_transferFrom() public {
        // given
        uint256 approveAmount = 200 * 10 ** bcn.decimals();
        uint256 transferAmount = 100 * 10 ** bcn.decimals();

        vm.prank(owner);
        vm.expectEmit(true, true, false, true);
        emit IERC20.Approval(owner, alice, approveAmount);

        bcn.approve(alice, approveAmount);
        assertEq(bcn.allowance(owner, alice), approveAmount);

        // when
        vm.prank(alice);
        vm.expectEmit(true, true, false, true);
        emit IERC20.Transfer(owner, bob, transferAmount);

        require(bcn.transferFrom(owner, bob, transferAmount), "Transfer failed");

        // then
        assertEq(bcn.balanceOf(bob), transferAmount);
        assertEq(bcn.balanceOf(owner), bcn.totalSupply() - transferAmount);
        assertEq(bcn.allowance(owner, alice), approveAmount - transferAmount);
    }

    function test_transferFrom_unlimited_allowance() public {
        // given
        uint256 approveAmount = type(uint256).max;
        uint256 transferAmount = 100 * 10 ** bcn.decimals();

        vm.prank(owner);
        vm.expectEmit(true, true, false, true);
        emit IERC20.Approval(owner, alice, approveAmount);

        bcn.approve(alice, approveAmount);
        assertEq(bcn.allowance(owner, alice), approveAmount);

        // when
        vm.prank(alice);
        vm.expectEmit(true, true, false, true);
        emit IERC20.Transfer(owner, bob, transferAmount);

        require(bcn.transferFrom(owner, bob, transferAmount), "Transfer failed");

        // then
        assertEq(bcn.allowance(owner, alice), approveAmount);
    }

    function test_transferFrom_revert_insufficient_allowance() public {
        // given
        uint256 allowance = bcn.allowance(owner, alice);
        assertEq(allowance, 0);

        uint256 transferAmount = 100 * 10 ** bcn.decimals();

        vm.prank(alice);
        vm.expectRevert(
            abi.encodeWithSelector(IERC20Errors.ERC20InsufficientAllowance.selector, alice, allowance, transferAmount)
        );

        // when/then
        require(!bcn.transferFrom(owner, bob, transferAmount), "Transfer succeded");
    }

    function test_transferOwnership() public {
        // given
        assertEq(bcn.owner(), owner);

        // when
        vm.prank(owner);
        bcn.transferOwnership(alice);

        // then
        assertEq(bcn.owner(), alice);
    }

    function test_transferOwnership_revert_not_owner() public {
        vm.prank(alice);
        vm.expectRevert("Not owner");
        bcn.transferOwnership(bob);
    }
}
