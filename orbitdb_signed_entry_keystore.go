package weshnet

import (
	"context"
	"encoding/hex"
	"sync"

	"berty.tech/go-ipfs-log/keystore"
	"github.com/libp2p/go-libp2p/core/crypto"

	"github.com/core-laboratories/weshnet/v2/pkg/errcode"
)

type WeshSignedKeyStore struct {
	sync.Map
}

func (s *WeshSignedKeyStore) SetKey(pk crypto.PrivKey) error {
	pubKeyBytes, err := pk.GetPublic().Raw()
	if err != nil {
		return errcode.ErrCode_TODO.Wrap(err)
	}

	keyID := hex.EncodeToString(pubKeyBytes)

	s.Store(keyID, pk)

	return nil
}

func (s *WeshSignedKeyStore) HasKey(_ context.Context, id string) (bool, error) {
	_, ok := s.Load(id)

	return ok, nil
}

func (s *WeshSignedKeyStore) CreateKey(ctx context.Context, id string) (crypto.PrivKey, error) {
	return s.GetKey(ctx, id)
}

func (s *WeshSignedKeyStore) GetKey(_ context.Context, id string) (crypto.PrivKey, error) {
	if privKey, ok := s.Load(id); ok {
		if pk, ok := privKey.(crypto.PrivKey); ok {
			return pk, nil
		}
	}

	return nil, errcode.ErrCode_ErrGroupMemberUnknownGroupID
}

func (s *WeshSignedKeyStore) Sign(privKey crypto.PrivKey, bytes []byte) ([]byte, error) {
	return privKey.Sign(bytes)
}

func (s *WeshSignedKeyStore) Verify(signature []byte, publicKey crypto.PubKey, data []byte) error {
	ok, err := publicKey.Verify(data, signature)
	if err != nil {
		return err
	}

	if !ok {
		return errcode.ErrCode_ErrGroupMemberLogEventSignature
	}

	return nil
}

func (s *WeshSignedKeyStore) getIdentityProvider() *weshSignedIdentityProvider {
	return &weshSignedIdentityProvider{
		keyStore: s,
	}
}

var _ keystore.Interface = (*WeshSignedKeyStore)(nil)
