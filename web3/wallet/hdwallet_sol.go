package wallet

import (
	"crypto/ed25519"
	"fmt"
	"log"

	"github.com/blocto/solana-go-sdk/pkg/hdwallet"
	"github.com/btcsuite/btcutil/base58"
)

type HDWalletSOL struct {
	AddressIndex int
	PrivateKey   string
	Address      string
	Key          *HDMasterKey
	SolKey       hdwallet.Key
	priKey       ed25519.PrivateKey
	pubKey       ed25519.PrivateKey
}

func (wallet *HDWalletSOL) NewWalletFromMnemonic(mnemonic string) (err error) {
	if wallet.Key == nil {
		wallet.Key = new(HDMasterKey)
	}
	err = wallet.Key.CreateByMnemonic(mnemonic, "")
	if err != nil {
		return err
	}
	// _ = wallet.DerivePrivateKey()
	return nil
}

func (wallet *HDWalletSOL) NewAccount() (*HDWalletSOL, error) {
	var account = new(HDWalletSOL)
	var err error
	account.Key = wallet.Key
	account.SolKey, err = wallet.DeriveSubKey(wallet.Key.Index)
	if err != nil {
		return account, err
	}
	wallet.Key.Index += 1
	_ = account.DerivePrivateKey()
	return account, err
}

func (wallet *HDWalletSOL) DeriveSubKey(index int) (hdwallet.Key, error) {
	// path := GetHDPath("SOL", index, 44)
	path := fmt.Sprintf("m/44'/501'/%d'/%d'", 0, index)
	k, err := hdwallet.Derived(path, wallet.Key.Seed)
	if err != nil {
		log.Fatalf(err.Error())
		return hdwallet.Key{}, err
	}
	return k, nil
}

func (wallet *HDWalletSOL) DerivePrivateKey() error {
	wallet.priKey = ed25519.NewKeyFromSeed(wallet.SolKey.PrivateKey)
	wallet.pubKey = make([]byte, ed25519.PublicKeySize)
	copy(wallet.pubKey, wallet.priKey[32:])
	wallet.PrivateKey = base58.Encode(wallet.priKey)
	wallet.Address = base58.Encode(wallet.pubKey)
	return nil
}
