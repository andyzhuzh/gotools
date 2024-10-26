package wallet

import (
	"encoding/hex"

	// "github.com/btcsuite/btcd/btcec"

	// "github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

type HDWalletBTC struct {
	AddressIndex           int
	PrivateKey             string
	WIF                    string
	WIFSegWitNative        string
	PrivateKeySegWitNative string
	PrivateKeySegWitNested string
	WIFSegWitNested        string
	PrivateKeyTAPROOT      string
	WIFTAPROOT             string

	AddressLegency      string
	AddressSegWitNative string
	AddressSegWitNested string
	AddressTAPROOT      string
	Key                 *HDMasterKey
	Key49               *HDMasterKey
	Key84               *HDMasterKey
	Key86               *HDMasterKey
	Counter             int
}

func (wallet *HDWalletBTC) NewWallet() (err error) {
	if wallet.Key == nil {
		wallet.Key = new(HDMasterKey)
	}
	err = wallet.Key.Create()
	if err != nil {
		return err
	}
	_ = wallet.DerivePrivateKey()

	return nil
}

func (wallet *HDWalletBTC) NewWalletFromMnemonic(mnemonic string) (err error) {
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

func (wallet *HDWalletBTC) NewWalletPrivate() (err error) {
	if wallet.Key == nil {
		wallet.Key = new(HDMasterKey)
	}
	wallet.Counter += 1
	privKey, err := NewPrivateKey()
	if err != nil {
		return err
	}

	wallet.Key.ECPivate = privKey
	wallet.Key.ECPublic = wallet.Key.ECPivate.PubKey()
	err = wallet.DerivePrivateKey()
	if err != nil {
		return err
	}
	if wallet.AddressLegency == "" && wallet.Counter <= 100 {
		wallet.NewWalletPrivate()
	}
	return nil
}

func (wallet *HDWalletBTC) NewAccount() (*HDWalletBTC, error) {
	var account = new(HDWalletBTC)
	var err error
	account.Key, err = wallet.Key.DeriveSubKey("BTC", wallet.Key.Index, 44) //LEGACY
	if err != nil {
		return account, err
	}

	account.Key49, err = wallet.Key.DeriveSubKey("BTC", wallet.Key.Index, 49) //SEGWIT_NESTED
	if err != nil {
		return account, err
	}

	account.Key84, err = wallet.Key.DeriveSubKey("BTC", wallet.Key.Index, 84) //SEGWIT_NATIVE
	if err != nil {
		return account, err
	}
	account.Key86, err = wallet.Key.DeriveSubKey("BTC", wallet.Key.Index, 86) //TAPROOT
	if err != nil {
		return account, err
	}

	network := &chaincfg.MainNetParams

	account.PrivateKey = hex.EncodeToString(account.Key.ECPivate.Serialize())
	WIF, _ := btcutil.NewWIF(account.Key.ECPivate, network, true)
	account.WIF = WIF.String()
	account.AddressLegency, _ = BTCPublicPubKeyToAddress(account.Key.ECPublic.SerializeCompressed(), LEGACY, network)

	account.PrivateKeySegWitNative = hex.EncodeToString(account.Key84.ECPivate.Serialize())
	WIFSegWitNative, _ := btcutil.NewWIF(account.Key84.ECPivate, network, true)
	account.WIFSegWitNative = WIFSegWitNative.String()
	account.AddressSegWitNative, _ = BTCPublicPubKeyToAddress(account.Key84.ECPublic.SerializeCompressed(), SEGWIT_NATIVE, network)

	account.PrivateKeySegWitNested = hex.EncodeToString(account.Key49.ECPivate.Serialize())
	WIFSegWitNested, _ := btcutil.NewWIF(account.Key49.ECPivate, network, true)
	account.WIFSegWitNested = WIFSegWitNested.String()
	account.AddressSegWitNested, _ = BTCPublicPubKeyToAddress(account.Key49.ECPublic.SerializeCompressed(), SEGWIT_NESTED, network)

	account.PrivateKeyTAPROOT = hex.EncodeToString(account.Key86.ECPivate.Serialize())
	WIFTAPROOT, _ := btcutil.NewWIF(account.Key86.ECPivate, network, true)
	account.WIFTAPROOT = WIFTAPROOT.String()
	account.AddressTAPROOT, _ = BTCPublicPubKeyToAddress(account.Key86.ECPublic.SerializeCompressed(), TAPROOT, network)

	wallet.Key.Index += 1
	return account, err
}

func (wallet *HDWalletBTC) NewFromPrivateKey(wif string) error {
	// 解码WIF格式的私钥
	privKey, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return err
	}
	if wallet.Key == nil {
		wallet.Key = new(HDMasterKey)
	}
	wallet.Key.ECPivate = privKey.PrivKey
	wallet.Key.ECPublic = wallet.Key.ECPivate.PubKey()
	err = wallet.DerivePrivateKey()
	if err != nil {
		return err
	}
	return nil
}

func (wallet *HDWalletBTC) DerivePrivateKey() error {
	// var errs error
	network := &chaincfg.MainNetParams

	// wallet.Private.Key
	// wallet.PrivateKey = hex.EncodeToString(wallet.Private.Serialize())
	wallet.PrivateKey = hex.EncodeToString(wallet.Key.ECPivate.Serialize())
	// wallet.PrivateKey = string(wallet.Private.Serialize())
	// PrivateKeyWIF, _ := btcutil.NewWIF(wallet.Private, network, true)
	PrivateKeyWIF, _ := btcutil.NewWIF(wallet.Key.ECPivate, network, true)
	wallet.WIF = PrivateKeyWIF.String()
	public := wallet.Key.ECPivate.PubKey()
	publicKey := public.SerializeCompressed()
	// if errs != nil {
	// 	return errs
	// }

	wallet.AddressLegency, _ = BTCPublicPubKeyToAddress(publicKey, LEGACY, network)
	wallet.AddressSegWitNative, _ = BTCPublicPubKeyToAddress(publicKey, SEGWIT_NATIVE, network)
	wallet.AddressSegWitNested, _ = BTCPublicPubKeyToAddress(publicKey, SEGWIT_NESTED, network)
	wallet.AddressTAPROOT, _ = BTCPublicPubKeyToAddress(publicKey, TAPROOT, network)
	return nil
}
