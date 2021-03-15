import cn.hyperchain.annotations.StoreField;
import cn.hyperchain.contract.BaseContract;
import cn.hyperchain.core.HyperMap;
import cn.hyperchain.core.Logger;

/**
 * Framework示例合约
 * 便于有一定基础的开发者熟悉一般业务逻辑的合约
 *
 * @author blocface
 */
public class Framework extends BaseContract implements FrameworkAdmin, FrameworkClient {
    /**
     * 提供critical，error，warning，notice，info，debug六种日志等级，且日志级别依次由高到低，日志输出结果会输出到平台的日志文件中。
     */
    Logger logger = Logger.getLogger(Framework.class);

    /**
     * 合约拥有者，即合约部署者
     */
    @StoreField
    private String owner;

    /**
     * 账户地址->用户角色
     */
    @StoreField
    private HyperMap<String, Role> addr2Role = new HyperMap<String, Role>();

    /**
     * 主键->数据表 的映射
     */
    @StoreField
    private HyperMap<Integer, MyData> id2MyData = new HyperMap<Integer, MyData>();

    /**
     * 构造函数
     * 合约部署者就是合约所有者。
     */
    @Override
    public void onInit() {
        String msgSender = getSender();
        this.owner = msgSender;
        addr2Role.put(msgSender, Role.ADMIN);
        logger.info("init owner=" + msgSender);
    }

    /**
     * 注册新用户或者修改已有的角色。仅限管理员操作，且无法修改合约拥有者的角色
     *
     * @param userAddr 待注册或待修改的用户账号地址
     * @param role     待注册为或待修改为的用户角色
     * @return 是否操作成功
     */
    @Override
    public Boolean chmod(String userAddr, Role role) {
        String msgSender = getSender();
        if (!Role.ADMIN.equals(this.addr2Role.get(msgSender))) {
            return false;
        }
        if (this.owner.equals(userAddr)) {
            return false;
        }
        this.addr2Role.put(userAddr, role);
        return true;
    }

    /**
     * 实现保存数据的接口，并且限制只有管理员才能操作
     *
     * @param myData 自定义数据，以JSON字符串形式输入
     * @return 返回是否插入成功
     */
    @Override
    public Boolean save(MyData myData) {
        String msgSender = getSender();
        if (!Role.ADMIN.equals(this.addr2Role.get(msgSender))) {
            return false;
        }
        this.id2MyData.put(myData.getId(), myData);
        return true;
    }

    /**
     * 实现按主键查询数据的接口，并且限制只有管理员和客户才能操作
     *
     * @param id 待查询的数据主键
     * @return 自定义数据
     */
    @Override
    public MyData query(Integer id) {
        String msgSender = getSender();
        if (!Role.ADMIN.equals(this.addr2Role.get(msgSender)) && !Role.CLIENT.equals(this.addr2Role.get(msgSender))) {
            return null;
        }
        return this.id2MyData.get(id);
    }
}
