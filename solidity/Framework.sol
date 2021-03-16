// 推荐使用Hyperchain目前支持的0.4.15版本及以上的Solidity编译器
pragma solidity ^0.4.15;

/// @title Framework示例合约
/// @notice 便于有一定基础的开发者熟悉一般业务逻辑的合约
/// @author blocface check
contract Framework {
    /// @notice 部署的时候被调用，保存合约部署者，并将合约部署者设置为合约管理员角色
    /// @dev 合约构造函数，可用于初始化合约状态变量。msg.sender表示当前交易发起方的账户地址，在构造函数中即为合约部署者
    function Framework() {
        owner = msg.sender;
        addr2Role[msg.sender] = Role.ADMIN;
    }

    /// @dev 合约拥有者，即合约部署者
    address owner;

    /// @dev 用户角色，用于权限校验
    /// @dev 注意，Role.ANY=0，用于占住零值
    enum Role {ANY, ADMIN, CLIENT}

    /// @dev 通用数据表
    struct MyData {
        uint256 id;
        bytes32 field1;
        bytes field2;
        string field3;
    }

    /// @dev 账户地址->用户角色
    mapping(address => Role) addr2Role;

    /// @dev 主键->数据表 的映射
    mapping(uint256 => MyData) id2MyData;


    /// @notice 注册新用户或者修改已有的角色。仅限管理员操作，且无法修改合约拥有者的角色
    /// @param userAddr 待注册或待修改的用户账号地址
    /// @param role 待注册为或待修改为的用户角色
    /// @return bool 是否操作成功
    function chmod(address userAddr, Role role) public returns (bool) {
        // 仅限合约管理员操作
        if (addr2Role[msg.sender] != Role.ADMIN) {
            return false;
        }
        // 无法修改合约拥有者的角色
        if (userAddr == owner) {
            return false;
        }
        // 修改用户的角色
        addr2Role[userAddr] = role;
        return true;
    }

    /// @notice 实现保存数据的接口，并且限制只有管理员才能操作
    /// @dev 这里体现了三种字符类型的入参
    /// @param field1 bytes32类型的数据，输入示例为 0x6162630000000000000000000000000000000000000000000000000000000000
    /// @param field2 bytes类型的数据，输入示例为 0x616263
    /// @param field3 string类型的数据，输入示例为 "abc"
    /// @return bool 返回是否插入成功
    function save(uint256 id, bytes32 field1, bytes field2, string field3) public returns (bool) {
        // 仅限合约管理员操作
        if (addr2Role[msg.sender] != Role.ADMIN) {
            return false;
        }
        MyData storage myData = id2MyData[id];
        myData.id = id;
        myData.field1 = field1;
        myData.field2 = field2;
        myData.field3 = field3;
        return true;
    }

    /// @notice 实现按主键查询数据的接口，并且限制只有管理员和客户才能操作
    /// @param id 待查询的数据主键
    /// @return field1 bytes32类型的数据
    /// @return field2 bytes类型的数据
    /// @return field3 string类型的数据
    function query(uint256 id) public returns (bytes32 field1, bytes field2, string field3) {
        // 仅限合约管理员或客户操作
        if (addr2Role[msg.sender] != Role.ADMIN && addr2Role[msg.sender] != Role.CLIENT) {
            return;
        }
        MyData memory myData = id2MyData[id];
        field1 = myData.field1;
        field2 = myData.field2;
        field3 = myData.field3;
        return;
    }
}
