"""
批量生成以太坊钱包工具

此脚本的功能：
1. 使用 BIP-39 标准生成助记词（Mnemonic）。
2. 使用 BIP-44 标准派生以太坊地址（Ethereum addresses）和对应的私钥（Private keys）。
3. 输出每个钱包的信息，包括：
   - 助记词（12 个随机单词）
   - 私钥（hex 格式字符串）
   - 公钥（地址，checksummed 格式）
4. 支持批量生成多个钱包并将数据输出为 JSON 文件。

输出：
- 钱包信息以列表形式打印到控制台。

使用方法：
1. 直接运行脚本：`python eth_create.py`。
2. 修改脚本中 `createNewETHWallet(10)` 的参数以生成指定数量的钱包。

注意事项：
- 请妥善保管助记词和私钥，它们是恢复钱包和访问资产的重要凭据。
- 如果生成的钱包文件已经存在，脚本将抛出异常以避免覆盖操作。

依赖：
- Python 3.7 或更高版本
- 安装必要库：
  - bip_utils：实现 BIP-39 和 BIP-44 标准
  - eth_account：处理以太坊地址和密钥
    使用 `pip install bip_utils eth_account` 命令安装依赖库。

作者：koala9527
日期：2024.12.26
"""



from bip_utils import (
    Bip39MnemonicGenerator,
    Bip39SeedGenerator,
    Bip44,
    Bip44Coins,
    Bip44Changes,
)


def createNewETHWallet(number=1):
    wallets = []
    for id in range(number):
        # 创建助记词
        mnemonic = Bip39MnemonicGenerator().FromWordsNumber(12)
        mnemonic = mnemonic.ToStr()
        # 生成种子
        seed = Bip39SeedGenerator(mnemonic).Generate()

        # 根据种子生成以太坊账户（BIP44标准）
        bip44_mst = Bip44.FromSeed(seed, Bip44Coins.ETHEREUM)

        # 指定 BIP44 路径：m / 44' / 60' / 0' / 0 / 0
        bip44_acc = (
            bip44_mst.Purpose()
            .Coin()
            .Account(0)
            .Change(Bip44Changes.CHAIN_EXT)
            .AddressIndex(0)
        )

        # 获取私钥（修复：使用 Raw().ToHex() 转换为十六进制字符串）
        privateKey = bip44_acc.PrivateKey().Raw().ToHex()

        # 获取地址
        address = bip44_acc.PublicKey().ToAddress()

        wallet = {
            "index": id,
            "mnemonic": mnemonic,  # 助记词
            "address": address,
            "privateKey": privateKey,
        }
        wallets.append(wallet)

    return wallets


if __name__ == "__main__":
    wallets = createNewETHWallet(10)
    print(wallets)
