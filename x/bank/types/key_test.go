package types_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

func cloneAppend(bz []byte, tail []byte) (res []byte) {
	res = make([]byte, len(bz)+len(tail))
	copy(res, bz)
	copy(res[len(bz):], tail)
	return
}

func TestAddressFromBalancesStore(t *testing.T) {
	addr, err := sdk.AccAddressFromHexUnsafe("0x99cfcc44e6bc906138b370f35c85583a3b54c8d8")
	require.NoError(t, err)
	addrLen := len(addr)
	require.Equal(t, 20, addrLen)
	key := cloneAppend(address.MustLengthPrefix(addr), []byte("stake"))

	tests := []struct {
		name        string
		key         []byte
		wantErr     bool
		expectedKey sdk.AccAddress
	}{
		{"valid", key, false, addr},
		{"#9111", []byte("\xff000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"), false, nil},
		{"empty", []byte(""), true, nil},
		{"invalid", []byte("3AA"), true, nil},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			addr, denom, err := types.AddressAndDenomFromBalancesStore(tc.key)
			if tc.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(types.ErrInvalidKey, err))
			} else {
				assert.NoError(t, err)
			}
			if len(tc.expectedKey) > 0 {
				assert.Equal(t, tc.expectedKey, addr)
				assert.Equal(t, "stake", denom)
			}
		})
	}
}

func TestCreateDenomAddressPrefix(t *testing.T) {
	require := require.New(t)

	key := types.CreateDenomAddressPrefix("")
	require.Len(key, len(types.DenomAddressPrefix)+1)
	require.Equal(append(types.DenomAddressPrefix, 0), key)

	key = types.CreateDenomAddressPrefix("abc")
	require.Len(key, len(types.DenomAddressPrefix)+4)
	require.Equal(append(types.DenomAddressPrefix, 'a', 'b', 'c', 0), key)
}
