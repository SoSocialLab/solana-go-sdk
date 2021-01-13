package tokenprog

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
)

const accountSize = 165

type AccountState uint8

const (
	AccountStateUninitialized AccountState = iota
	AccountStateInitialized
	AccountFrozen
)

var (
	Some = []byte{1, 0, 0, 0}
	None = []byte{0, 0, 0, 0}
)

// Account is token program account
// https://github.com/solana-labs/solana-program-library/blob/master/token/program-v3/src/state.rs#L86
type Account struct {
	Mint            common.PublicKey
	Owner           common.PublicKey
	Amount          uint64
	Delegate        *common.PublicKey
	State           AccountState
	IsNative        *uint64
	DelegatedAmount uint64
	CloseAuthority  *common.PublicKey
}

func AccountFromData(data []byte) (*Account, error) {
	if len(data) != accountSize {
		return nil, fmt.Errorf("data length not match")
	}

	mint := common.PublicKeyFromBytes(data[:32])

	owner := common.PublicKeyFromBytes(data[32:64])

	amount := binary.LittleEndian.Uint64(data[64:72])

	var delegate *common.PublicKey
	if bytes.Equal(data[72:76], Some) {
		key := common.PublicKeyFromBytes(data[76:108])
		delegate = &key
	}

	state := AccountState(data[108])

	var isNative *uint64
	if bytes.Equal(data[109:113], Some) {
		num := binary.LittleEndian.Uint64(data[113:121])
		isNative = &num
	}

	delegatedAmount := binary.LittleEndian.Uint64(data[121:129])

	var closeAuthority *common.PublicKey
	if bytes.Equal(data[129:133], Some) {
		key := common.PublicKeyFromBytes(data[133:165])
		closeAuthority = &key
	}

	return &Account{
		Mint:            mint,
		Owner:           owner,
		Amount:          amount,
		Delegate:        delegate,
		State:           state,
		IsNative:        isNative,
		DelegatedAmount: delegatedAmount,
		CloseAuthority:  closeAuthority,
	}, nil
}