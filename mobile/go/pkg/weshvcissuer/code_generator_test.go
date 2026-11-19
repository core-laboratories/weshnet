package weshvcissuer_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/weshvcissuer"
	"github.com/core-laboratories/weshnet/v2/pkg/verifiablecredstypes"
)

func Test_CodeGeneratorEightDigits(t *testing.T) {
	f := weshvcissuer.CodeGeneratorEightDigits([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	code, err := f(&verifiablecredstypes.StateCode{})
	require.Error(t, err)
	require.Empty(t, code)

	code1Data, err := f(&verifiablecredstypes.StateCode{
		AccountLink: "wesh://1",
		Identifier:  "+15555555555",
		Timestamp:   []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	})
	require.NoError(t, err)
	require.NotEmpty(t, code1Data)
	require.Len(t, code1Data, 8)

	code1Data2, err := f(&verifiablecredstypes.StateCode{
		AccountLink: "wesh://2",
		Identifier:  "+15555555555",
		Timestamp:   []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 0},
	})
	require.NoError(t, err)
	require.NotEmpty(t, code1Data2)
	require.Len(t, code1Data2, 8)

	f2 := weshvcissuer.CodeGeneratorEightDigits([]byte{0, 1, 2, 3, 4, 5, 6, 7})
	code2, err := f2(&verifiablecredstypes.StateCode{})
	require.Error(t, err)
	require.Empty(t, code2)

	code2Data, err := f2(&verifiablecredstypes.StateCode{
		AccountLink: "wesh://1",
		Identifier:  "+15555555555",
		Timestamp:   []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	})
	require.NoError(t, err)
	require.NotEmpty(t, code2Data)
	require.Len(t, code2Data, 8)

	require.NotEqual(t, code1Data, code2Data)
	require.NotEqual(t, code1Data, code1Data2)
}
