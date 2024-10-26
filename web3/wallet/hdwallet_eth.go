package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	// "github.com/btcsuite/btcd/btcec"
	// "github.com/ethereum/go-ethereum/accounts/keystore"
)

type HDWalletETH struct {
	AddressIndex int
	PrivateKey   string
	Address      string
	Key          *HDMasterKey
}

func (wallet *HDWalletETH) NewAccount() (*HDWalletETH, error) {
	var account = new(HDWalletETH)
	var err error
	account.Key, err = wallet.Key.DeriveSubKey("ETH", wallet.Key.Index, 44)
	if err != nil {
		return account, err
	}
	wallet.Key.Index += 1
	_ = account.DerivePrivateKey()
	return account, err
}

func (wallet *HDWalletETH) NewWallet() (err error) {
	err = wallet.Key.Create()
	if err != nil {
		return err
	}
	_ = wallet.DerivePrivateKey()
	return nil
}
func (wallet *HDWalletETH) NewWalletFromMnemonic(mnemonic string) (err error) {
	if wallet.Key == nil {
		wallet.Key = new(HDMasterKey)
	}
	err = wallet.Key.CreateByMnemonic(mnemonic, "")
	if err != nil {
		return err
	}
	_ = wallet.DerivePrivateKey()
	return nil
}

func (wallet *HDWalletETH) DerivePrivateKey() (err error) {
	wallet.PrivateKey = hex.EncodeToString(crypto.FromECDSA(wallet.Key.ECDSAPrivate))
	pubKey := wallet.Key.ECDSAPublic
	wallet.Address = crypto.PubkeyToAddress(*pubKey).Hex()
	return nil
}

func (wallet *HDWalletETH) NewFromPrivateKey(HexPrivateKey string) error {
	wallet.PrivateKey = HexPrivateKey
	privateKey, err := crypto.HexToECDSA(HexPrivateKey)
	if err != nil {
		return err
	}

	// hex.DecodeString('aaa')
	derivedPublicKey := privateKey.Public().(*ecdsa.PublicKey)
	wallet.Address = crypto.PubkeyToAddress(*derivedPublicKey).Hex()
	return nil
}

func newWalletKeystore(passPhrase string) {
	// 第一个参数：./keystore 指定账户所在的文件夹
	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	// acc, err := ks.NewAccount("your password")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("账户地址：", acc.Address.String())
	for _, acc := range ks.Accounts() {
		fmt.Println("账户地址：", acc.Address.String())
	}

	// 如果希望签名交易的话，可以使用钱包
	// for _, wallet := range ks.Wallets() {
	// 	// 这里每个钱包获得的第一个account和上面第4步的账户是一一对应的
	// 	_, _ = wallet.SignTx(wallet.Accounts()[0], nil, nil)
	// 	_, _ = wallet.SignData(wallet.Accounts()[0], "", nil)
	// }

	// // 先把账户的 json 文件导出为[]byte
	kjson, err := ks.Export(ks.Accounts()[0], passPhrase, "newpassword")
	if err != nil {
		panic(err.Error())
	}

	// 通过kjson+密码可以解密私钥
	key, err := keystore.DecryptKey(kjson, "newpassword")
	if err != nil {
		panic(err.Error())
	}

	privateKey := key.PrivateKey
	fmt.Println(common.BytesToHash(crypto.FromECDSA(privateKey)).String())
}
