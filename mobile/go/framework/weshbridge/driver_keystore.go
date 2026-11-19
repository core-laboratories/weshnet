package weshbridge

import "github.com/core-laboratories/weshnet/mobile/v2/go/internal/accountutils"

type NativeKeystoreDriver interface {
	accountutils.NativeKeystore
}
