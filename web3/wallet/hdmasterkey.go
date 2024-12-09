package wallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	// "github.com/btcsuite/btcd/btcec"

	// "github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const (
	LEGACY        string = "Legacy"
	TAPROOT       string = "Taproot"
	SEGWIT_NESTED string = "SegWit_Nested"
	SEGWIT_NATIVE string = "SegWit_Native"
)

type HDMasterKey struct {
	Mnemonic string
	Password string
	Seed     []byte
	Root     *bip32.Key
	Parent   *bip32.Key
	CoinType string
	Index    int
	HdPath   string
	Key      *bip32.Key
	// ExtendedKey  *hdkeychain.ExtendedKey
	ECPivate     *btcec.PrivateKey
	ECPublic     *btcec.PublicKey
	ECDSAPrivate *ecdsa.PrivateKey
	ECDSAPublic  *ecdsa.PublicKey
}

func (hdmkey *HDMasterKey) CreateByMnemonic(mnemonic string, password string) (err error) {
	err = hdmkey.NewMasterKey(mnemonic, password)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return
}

func (hdmkey *HDMasterKey) Create() (err error) {
	mnemonic := NewMnemonic()
	// mnemonic := "kingdom track about unveil coconut bid hard circle solve outdoor swim flash"
	err = hdmkey.NewMasterKey(mnemonic, "")
	if err != nil {
		log.Fatalf(err.Error())
	}

	// fmt.Printf("MasterKey EC Private:%v\n", hdmkey.ECPivate)
	// fmt.Printf("MasterKey EC ECPublic:%v\n", hdmkey.ECPublic)
	// fmt.Printf("MasterKey EC Private to ECDSA:%v\n", hdmkey.ECPivate.ToECDSA())
	return
}

func (hdmkey *HDMasterKey) DeriveSubKey(symbol string, index int, propose int) (retKey *HDMasterKey, errs error) {
	path := GetHDPath(symbol, index, propose)
	// fmt.Println(path)
	retKey, errs = hdmkey.NewChildKeyByPathString(path)
	retKey.CoinType = symbol
	retKey.Index = index
	return
}

func (hdmkey *HDMasterKey) NewChildKeyByPathString(childPath string) (retWallet *HDMasterKey, errs error) {
	var uid uint32
	var subKey *bip32.Key
	retWallet = new(HDMasterKey)
	subKey = hdmkey.Key
	arr := strings.Split(childPath, "/")
	// currentKey := key.MasterKey
	for _, part := range arr {
		if part == "m" {
			continue
		}

		var harden = false
		if strings.HasSuffix(part, "'") {
			harden = true
			part = strings.TrimSuffix(part, "'")
		}

		id, err := strconv.ParseUint(part, 10, 31)
		if err != nil {
			return nil, err
		}

		uid = uint32(id)
		if harden {
			uid |= bip32.FirstHardenedChild
		}

		subKey, err = subKey.NewChildKey(uid)
		if err != nil {
			log.Fatalf(err.Error())
			return nil, err
		}
	}
	retWallet.Key = subKey
	retWallet.Root = hdmkey.Root
	retWallet.Parent = hdmkey.Key
	retWallet.HdPath = childPath
	retWallet.ECPivate = DeriveECPrivate(retWallet.Key)
	retWallet.ECPublic = DeriveECPublic(retWallet.Key) //btcec.ParsePubKey(retWallet.Key.PublicKey().Key)
	retWallet.ECDSAPrivate = DeriveECDSAPrivate(retWallet.Key)
	retWallet.ECDSAPublic = DeriveECDSAPublic(retWallet.Key)
	return retWallet, nil
}

func (hdmkey *HDMasterKey) NewMasterKey(mnemonic string, password string) (err error) {
	// hdmkey = new(HDMasterKey)
	hdmkey.Mnemonic = mnemonic
	hdmkey.Password = password
	hdmkey.Seed = NewSeed(hdmkey.Mnemonic, hdmkey.Password)
	hdmkey.Root, err = bip32.NewMasterKey(hdmkey.Seed)
	if err != nil {
		log.Fatalf(err.Error())
	}
	hdmkey.Key = hdmkey.Root
	// hdmkey.ExtendedKey = tempKey.ExtendedKey
	hdmkey.ECPivate = DeriveECPrivate(hdmkey.Key)
	hdmkey.ECPublic = DeriveECPublic(hdmkey.Key)
	hdmkey.ECDSAPrivate = DeriveECDSAPrivate(hdmkey.Key) // retWallet.ECPivate.ToECDSA()
	hdmkey.ECDSAPublic = DeriveECDSAPublic(hdmkey.Key)
	return
}

func NewMnemonic() string {
	// 由熵源生成助记词
	// @参数 128 => 12个单词
	// @参数 256 => 24个单词
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return string(mnemonic)
}

func NewPrivateKey() (*secp256k1.PrivateKey, error) {
	return btcec.NewPrivateKey()
}

func NewSeed(mnemonic string, password string) []byte {
	// 由助记词生成种子(Seed)
	seed := bip39.NewSeed(mnemonic, password)
	return seed
}

func GetHDPath(symbol string, index int, propose int) string {
	// 44-通用多地址，BTC p2pkh 常规地址
	// 49-BTC P2SH 隔离见证地址（ Nested SegWit） P2WPKH-nested-in-P2SH
	// 84-BTC P2WPKH NativeSegWit
	// 86-BTC P2R Taproot 地址
	coinIndex := GetBIP44Index(symbol)
	var pathPropose int
	pathPropose = propose
	if pathPropose == 0 {
		pathPropose = 44
	}
	return fmt.Sprintf(`m/%d'/%d'/0'/0/%d`, pathPropose, coinIndex, index)

	// return fmt.Sprintf(`m/49'/%d'/0'/0/%d`, coinIndex, index) //btc隔离见证地址 P2WPKH-nested-in-P2SH based accounts Nested SegWit
	// return fmt.Sprintf(`m/44'/%d'/0'/0/%d`, coinIndex, index)
	// return fmt.Sprintf(`m/84'/%d'/0'/0/%d`, coinIndex, index)
	// return fmt.Sprintf(`m/86'/%d'/0'/0/%d`, coinIndex, index)
}

func GetBIP44Index(symbol string) (index int) {
	BIP44Dict := map[string]int{
		"BTC": 0,
		"ETH": 60,
		"TRX": 195,
		"SUI": 784,
		"SOL": 501,
	}
	value, exists := BIP44Dict[symbol]
	if exists {
		return value
	} else {
		return 0
	}
}

func DeriveECPrivate(masterKey *bip32.Key) *btcec.PrivateKey {
	return secp256k1.PrivKeyFromBytes(masterKey.Key)
}
func DeriveECPublic(masterKey *bip32.Key) *btcec.PublicKey {
	// return btcec.ParsePubKey(masterKey.PublicKey().Key)
	return secp256k1.PrivKeyFromBytes(masterKey.Key).PubKey()
}

func DeriveECDSAPrivate(masterKey *bip32.Key) *ecdsa.PrivateKey {
	return secp256k1.PrivKeyFromBytes(masterKey.Key).ToECDSA()
}
func DeriveECDSAPublic(masterKey *bip32.Key) *ecdsa.PublicKey {
	publicKey := secp256k1.PrivKeyFromBytes(masterKey.Key).ToECDSA().Public()
	ECDSAPublic, _ := publicKey.(*ecdsa.PublicKey)
	return ECDSAPublic
}

func BTCPublicPubKeyToAddress(publicKey []byte, addrType string, network *chaincfg.Params) (string, error) {
	// publicKey := ECpublic.SerializeCompressed()
	if network == nil {
		network = &chaincfg.MainNetParams
	}
	if addrType == LEGACY {
		p2pkh, err := btcutil.NewAddressPubKey(publicKey, network)
		if err != nil {
			return "", err
		}

		return p2pkh.EncodeAddress(), nil
	} else if addrType == SEGWIT_NATIVE {
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), network)
		if err != nil {
			return "", err
		}

		return p2wpkh.EncodeAddress(), nil
	} else if addrType == SEGWIT_NESTED {
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), network)
		if err != nil {
			return "", err
		}
		redeemScript, err := txscript.PayToAddrScript(p2wpkh)
		if err != nil {
			return "", err
		}
		p2sh, err := btcutil.NewAddressScriptHash(redeemScript, network)
		if err != nil {
			return "", err
		}

		return p2sh.EncodeAddress(), nil
	} else if addrType == TAPROOT {
		internalKey, err := btcec.ParsePubKey(publicKey)
		if err != nil {
			return "", err
		}
		p2tr, err := btcutil.NewAddressTaproot(txscript.ComputeTaprootKeyNoScript(internalKey).SerializeCompressed()[1:], network)
		if err != nil {
			return "", err
		}

		return p2tr.EncodeAddress(), nil
	} else {
		return "", errors.New("address type not supported")
	}
}
