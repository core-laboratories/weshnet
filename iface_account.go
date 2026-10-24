package weshnet

import (
	"github.com/libp2p/go-libp2p/core/crypto"

	"github.com/core-laboratories/weshnet/v2/pkg/protocoltypes"
	"github.com/core-laboratories/weshnet/v2/pkg/secretstore"
)

type AccountKeys interface {
	AccountPrivKey() (crypto.PrivKey, error)
	AccountProofPrivKey() (crypto.PrivKey, error)
	DevicePrivKey() (crypto.PrivKey, error)
	ContactGroupPrivKey(pk crypto.PubKey) (crypto.PrivKey, error)
	MemberDeviceForGroup(g *protocoltypes.Group) (secretstore.OwnMemberDevice, error)
}
