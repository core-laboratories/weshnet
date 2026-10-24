package weshnet

import (
	"context"

	"berty.tech/go-ipfs-log/identityprovider"
	"github.com/libp2p/go-libp2p/core/crypto"

	"github.com/core-laboratories/weshnet/v2/pkg/errcode"
)

const (
	identityGroupIDKey = "group_id"
	storeTypeKey       = "store_type"
	identityType       = "betry_group_entry"
)

type weshSignedIdentityProvider struct {
	keyStore *WeshSignedKeyStore
}

func (b *weshSignedIdentityProvider) UnmarshalPublicKey(data []byte) (crypto.PubKey, error) {
	return crypto.UnmarshalPublicKey(data)
}

func (b *weshSignedIdentityProvider) GetID(_ context.Context, opts *identityprovider.CreateIdentityOptions) (string, error) {
	return opts.ID, nil
}

//nolint:revive
func (b *weshSignedIdentityProvider) SignIdentity(ctx context.Context, data []byte, id string) ([]byte, error) {
	return nil, nil
}

func (b *weshSignedIdentityProvider) GetType() string {
	return identityType
}

func (b *weshSignedIdentityProvider) VerifyIdentity(*identityprovider.Identity) error {
	return nil
}

func (b *weshSignedIdentityProvider) Sign(ctx context.Context, identity *identityprovider.Identity, bytes []byte) ([]byte, error) {
	key, err := b.keyStore.GetKey(ctx, identity.ID)
	if err != nil {
		return nil, errcode.ErrCode_TODO.Wrap(err)
	}

	sig, err := key.Sign(bytes)
	if err != nil {
		return nil, errcode.ErrCode_TODO.Wrap(err)
	}

	return sig, nil
}

func (b *weshSignedIdentityProvider) signID(ctx context.Context, id string) (crypto.PubKey, []byte, error) {
	privKey, err := b.keyStore.GetKey(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	idSignature, err := b.keyStore.Sign(privKey, []byte(id))
	if err != nil {
		return nil, nil, err
	}

	return privKey.GetPublic(), idSignature, nil
}

func (b *weshSignedIdentityProvider) createIdentity(ctx context.Context, options *identityprovider.CreateIdentityOptions) (*identityprovider.Identity, error) {
	id, err := b.GetID(ctx, options)
	if err != nil {
		return nil, err
	}

	publicKey, idSignature, err := b.signID(ctx, id)
	if err != nil {
		return nil, err
	}

	publicKeyRaw, err := publicKey.Raw()
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := crypto.MarshalPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	pubKeyIDSignature, err := b.SignIdentity(ctx, append(publicKeyRaw, idSignature...), options.ID)
	if err != nil {
		return nil, err
	}

	return &identityprovider.Identity{
		ID:        id,
		PublicKey: publicKeyBytes,
		Signatures: &identityprovider.IdentitySignature{
			ID:        idSignature,
			PublicKey: pubKeyIDSignature,
		},
		Type:     b.GetType(),
		Provider: b,
	}, nil
}

var _ identityprovider.Interface = (*weshSignedIdentityProvider)(nil)
