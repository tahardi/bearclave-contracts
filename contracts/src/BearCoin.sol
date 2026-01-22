pragma solidity ^0.8.33;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract BearCoin is ERC20 {
    uint8 public constant DECIMALS = 18;
    uint256 public constant TOTAL_SUPPLY = 1_000_000 * 10 ** uint256(DECIMALS);

    address public owner;

    event Mint(address indexed to, uint256 amount);
    event Burn(address indexed from, uint256 amount);

    constructor() ERC20("BearCoin", "BCN") {
        owner = msg.sender;
        _mint(msg.sender, TOTAL_SUPPLY);
    }

    function decimals() public pure override returns (uint8) {
        return DECIMALS;
    }

    function mint(address recipient, uint256 amount) public onlyOwner {
        require(totalSupply() + amount <= TOTAL_SUPPLY, "Minting exceeds total supply");
        _mint(recipient, amount);
        emit Mint(recipient, amount);
    }

    function burn(uint256 amount) public {
        _burn(msg.sender, amount);
        emit Burn(msg.sender, amount);
    }

    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0), "New owner cannot be null");
        owner = newOwner;
    }

    // https://getfoundry.sh/forge/linting/#unwrapped-modifier-logic
    modifier onlyOwner() {
        _checkOwner(msg.sender);
        _;
    }

    function _checkOwner(address who) internal view {
        require(who == owner, "Not owner");
    }
}
