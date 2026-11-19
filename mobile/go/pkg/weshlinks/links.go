package weshlinks

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"net/url"
	"strings"

	"github.com/eknkc/basex"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/proto"

	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/errcode"
	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/messengertypes"
	"github.com/core-laboratories/weshnet/v2/pkg/cryptoutil"
	"github.com/core-laboratories/weshnet/v2/pkg/protocoltypes"
)

// MarshalLink returns shareable web and internal URLs.
//
// The web URL is meant to:
// - be short,
// - have some parts that are human-readable,
// - to point to a sub-page of the wesh.tech website where some JS code will parse the human-readable part.
//
// The internal URL is meant to generate the most tiny QR codes. These QR codes can only be opened by a Wesh app.
//
// Marshal will return an error if the provided link does not contain all the mandatory fields;
// it may also filter-out some sensitive data.
func MarshalLink(link *messengertypes.WeshLink) (internal string, web string, err error) {
	if link == nil || link.Kind == messengertypes.WeshLink_UnknownKind {
		return "", "", errcode.ErrCode_ErrMissingInput
	}

	if err := link.IsValid(); err != nil {
		return "", "", err
	}

	var (
		// web
		kind    string
		machine = &messengertypes.WeshLink{}
		human   = url.Values{}

		// internal
		qrOptimized *messengertypes.WeshLink
	)

	switch link.Kind {
	case messengertypes.WeshLink_ContactInviteV1Kind:
		kind = "contact"
		machine.WeshId = &messengertypes.WeshID{
			PublicRendezvousSeed: link.WeshId.PublicRendezvousSeed,
			AccountPk:            link.WeshId.AccountPk,
		}
		if link.WeshId.DisplayName != "" {
			human.Add("name", link.WeshId.DisplayName)
		}

		// for contact sharing, there are no fields to hide, so just copy the input link
		qrOptimized = proto.Clone(link).(*messengertypes.WeshLink)
	case messengertypes.WeshLink_GroupV1Kind:
		kind = "group"
		machine.WeshGroup = &messengertypes.WeshGroup{
			Group: &protocoltypes.Group{
				PublicKey:  link.WeshGroup.Group.PublicKey,
				Secret:     link.WeshGroup.Group.Secret,
				SecretSig:  link.WeshGroup.Group.SecretSig,
				GroupType:  link.WeshGroup.Group.GroupType,
				SignPub:    link.WeshGroup.Group.SignPub,
				LinkKeySig: link.WeshGroup.Group.LinkKeySig,
			},
		}
		if link.WeshGroup.DisplayName != "" {
			human.Add("name", link.WeshGroup.DisplayName)
		}
		qrOptimized = proto.Clone(link).(*messengertypes.WeshLink)
	case messengertypes.WeshLink_EncryptedV1Kind:
		kind = "enc"
		machine.Encrypted = &messengertypes.WeshLink_Encrypted{
			Kind:     link.Encrypted.Kind,
			Nonce:    link.Encrypted.Nonce,
			Checksum: link.Encrypted.Checksum,
		}
		if link.Encrypted.DisplayName != "" {
			human.Add("name", link.Encrypted.DisplayName)
		}
		switch link.Encrypted.Kind {
		case messengertypes.WeshLink_ContactInviteV1Kind:
			machine.Encrypted.ContactAccountPk = link.Encrypted.ContactAccountPk
			machine.Encrypted.ContactPublicRendezvousSeed = link.Encrypted.ContactPublicRendezvousSeed
		case messengertypes.WeshLink_GroupV1Kind:
			machine.Encrypted.GroupPublicKey = link.Encrypted.GroupPublicKey
			machine.Encrypted.GroupSecret = link.Encrypted.GroupSecret
			machine.Encrypted.GroupSecretSig = link.Encrypted.GroupSecretSig
			machine.Encrypted.GroupSignPub = link.Encrypted.GroupSignPub
			machine.Encrypted.GroupType = link.Encrypted.GroupType
			machine.Encrypted.GroupLinkKeySig = link.Encrypted.GroupLinkKeySig
		}
		qrOptimized = proto.Clone(link).(*messengertypes.WeshLink)
	case messengertypes.WeshLink_MessageV1Kind:
		kind = "message"
		machine.WeshMessageRef = &messengertypes.WeshLink_WeshMessageRef{
			AccountId: link.WeshMessageRef.AccountId,
			GroupPk:   link.WeshMessageRef.GroupPk,
			MessageId: link.WeshMessageRef.MessageId,
		}
		qrOptimized = proto.Clone(link).(*messengertypes.WeshLink)
	default:
		return "", "", errcode.ErrCode_ErrInvalidInput
	}

	// compute the web shareable link.
	// in this mode, we have:
	// - a human-readable link kind
	// - a base58-encoded binary (proto) representation of the link (without the kind and metadata)
	// - human-readable metadata, encoded as query string (including display name)
	{
		machineBin, err := proto.Marshal(machine)
		if err != nil {
			return "", "", errcode.ErrCode_ErrInvalidInput.Wrap(err)
		}
		// here we use base58 which is compressed enough whilst being easy to read by a human.
		// another candidate could be base58.RawURLEncoding which is a little bit more compressed and also only containing unescaped URL chars.
		machineEncoded := base58.Encode(machineBin)
		path := kind + "/" + machineEncoded
		if len(human) > 0 {
			path += "/" + human.Encode()
		}
		// we use a '#' to improve privacy by preventing the webservers to get aware of the right part of this URL
		web = LinkWebPrefix + path
	}

	// compute the internal shareable link.
	// in this mode, the url is as short as possible, in the format: wesh://{base45(proto.marshal(link))}.
	{
		qrBin, err := proto.Marshal(qrOptimized)
		if err != nil {
			return "", "", errcode.ErrCode_ErrInvalidInput.Wrap(err)
		}
		// using uppercase to stay in the QR AlphaNum's 45chars alphabet
		internal = LinkInternalPrefix + "PB/" + qrBaseEncoder.Encode(qrBin)
	}

	return internal, web, nil
}

// UnmarshalLink takes an URL generated by WeshLink.Marshal (or manually crafted), and returns a WeshLink object.
func UnmarshalLink(uri string, key []byte) (*messengertypes.WeshLink, error) {
	if uri == "" {
		return nil, errcode.ErrCode_ErrMissingInput
	}

	// internal format
	if strings.HasPrefix(strings.ToLower(uri), strings.ToLower(LinkInternalPrefix)) {
		right := uri[len(LinkInternalPrefix):]
		parts := strings.Split(right, "/")
		if len(parts) < 2 {
			return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("URI should have at least 2 parts"))
		}
		switch strings.ToLower(parts[0]) {
		case "pb":
			blob := strings.Join(parts[1:], "/")
			qrBin, err := qrBaseEncoder.Decode(blob)
			if err != nil {
				return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
			}
			var link messengertypes.WeshLink
			err = proto.Unmarshal(qrBin, &link)
			if err != nil {
				return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
			}

			if link.Kind == messengertypes.WeshLink_EncryptedV1Kind {
				return decryptLink(&link, key)
			}

			return &link, nil
		default:
			return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("unsupported link type"))
		}
	}

	// web format
	if strings.HasPrefix(strings.ToLower(uri), strings.ToLower(LinkWebPrefix)) {
		parsed, err := url.Parse(uri)
		if err != nil {
			return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
		}
		if parsed.Fragment == "" {
			return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
		}

		rawFragment := strings.Join(strings.Split(uri, "#")[1:], "#") // required by go1.14
		// when minimal version of wesh will be go1.15, we can just use `parsed.EscapedFragment()`

		link := messengertypes.WeshLink{}
		parts := strings.Split(rawFragment, "/")
		if len(parts) < 2 {
			return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("URI should have at least 2 parts"))
		}

		// decode blob
		machineBin, err := base58.Decode(parts[1])
		if err != nil {
			return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
		}
		if err := proto.Unmarshal(machineBin, &link); err != nil {
			return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
		}

		// decode url.Values
		var human url.Values
		if len(parts) > 2 {
			encodedValues := strings.Join(parts[2:], "/")
			human, err = url.ParseQuery(encodedValues)
			if err != nil {
				return nil, errcode.ErrCode_ErrInvalidInput.Wrap(err)
			}
		}

		// per-kind merging strategies and checks
		switch kind := parts[0]; kind {
		case "contact":
			link.Kind = messengertypes.WeshLink_ContactInviteV1Kind
			if link.WeshId == nil {
				link.WeshId = &messengertypes.WeshID{}
			}
			if name := human.Get("name"); name != "" && link.WeshId.DisplayName == "" {
				link.WeshId.DisplayName = name
			}
		case "group":
			link.Kind = messengertypes.WeshLink_GroupV1Kind
			if link.WeshGroup == nil {
				link.WeshGroup = &messengertypes.WeshGroup{}
			}
			if name := human.Get("name"); name != "" && link.WeshGroup.DisplayName == "" {
				link.WeshGroup.DisplayName = name
			}
		case "enc":
			link.Kind = messengertypes.WeshLink_EncryptedV1Kind
			if link.Encrypted == nil {
				return nil, errcode.ErrCode_ErrInvalidInput
			}
			if name := human.Get("name"); name != "" && link.Encrypted.DisplayName == "" {
				link.Encrypted.DisplayName = name
			}
		case "message":
			link.Kind = messengertypes.WeshLink_MessageV1Kind
			if link.Encrypted == nil {
				link.WeshMessageRef = &messengertypes.WeshLink_WeshMessageRef{}
			}
		default:
			return nil, errcode.ErrCode_ErrInvalidInput
		}

		if link.Kind == messengertypes.WeshLink_EncryptedV1Kind {
			return decryptLink(&link, key)
		}

		return &link, nil
	}

	return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("unsupported link format"))
}

func decryptLink(link *messengertypes.WeshLink, passphrase []byte) (*messengertypes.WeshLink, error) {
	if link == nil || link.Kind != messengertypes.WeshLink_EncryptedV1Kind {
		return nil, errcode.ErrCode_ErrInvalidInput
	}
	if err := link.IsValid(); err != nil {
		return nil, err
	}
	if passphrase == nil || string(passphrase) == "" {
		return link, nil
	}

	decrypted := messengertypes.WeshLink{
		Kind: link.Encrypted.Kind,
	}

	// derive key for AES
	key, _, err := cryptoutil.DeriveKey(passphrase, link.Encrypted.Nonce)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	// create stream cipher
	stream, err := cryptoutil.AESCTRStream(key, link.Encrypted.Nonce)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	switch decrypted.Kind {
	case messengertypes.WeshLink_ContactInviteV1Kind:
		decrypted.WeshId = &messengertypes.WeshID{
			AccountPk:            make([]byte, len(link.Encrypted.ContactAccountPk)),
			PublicRendezvousSeed: make([]byte, len(link.Encrypted.ContactPublicRendezvousSeed)),
		}
		stream.XORKeyStream(decrypted.WeshId.PublicRendezvousSeed, link.Encrypted.ContactPublicRendezvousSeed)
		stream.XORKeyStream(decrypted.WeshId.AccountPk, link.Encrypted.ContactAccountPk)
		decrypted.WeshId.DisplayName = link.Encrypted.DisplayName

	case messengertypes.WeshLink_GroupV1Kind:
		decrypted.WeshGroup = &messengertypes.WeshGroup{
			Group: &protocoltypes.Group{
				PublicKey:  make([]byte, len(link.Encrypted.GroupPublicKey)),
				Secret:     make([]byte, len(link.Encrypted.GroupSecret)),
				SecretSig:  make([]byte, len(link.Encrypted.GroupSecretSig)),
				SignPub:    make([]byte, len(link.Encrypted.GroupSignPub)),
				LinkKeySig: make([]byte, len(link.Encrypted.GroupLinkKeySig)),
				GroupType:  link.Encrypted.GroupType,
			},
		}
		stream.XORKeyStream(decrypted.WeshGroup.Group.PublicKey, link.Encrypted.GroupPublicKey)
		stream.XORKeyStream(decrypted.WeshGroup.Group.Secret, link.Encrypted.GroupSecret)
		stream.XORKeyStream(decrypted.WeshGroup.Group.SecretSig, link.Encrypted.GroupSecretSig)
		stream.XORKeyStream(decrypted.WeshGroup.Group.SignPub, link.Encrypted.GroupSignPub)
		stream.XORKeyStream(decrypted.WeshGroup.Group.LinkKeySig, link.Encrypted.GroupLinkKeySig)
		decrypted.WeshGroup.DisplayName = link.Encrypted.DisplayName
	}

	if len(link.Encrypted.Checksum) > 0 {
		checksum := make([]byte, len(link.Encrypted.Checksum))
		err := clearLinkChecksum(&decrypted, checksum)
		if err != nil {
			return nil, errcode.ErrCode_ErrInternal.Wrap(err)
		}
		if !bytes.Equal(link.Encrypted.Checksum, checksum) {
			return nil, errcode.ErrCode_ErrMessengerDeepLinkInvalidPassphrase
		}
	}

	return &decrypted, nil
}

// EncryptLink converts a clear WeshLink into an encrypted one.
//
// The encrypted link can be marshaled by UnmarshalLink if a passphrase is provided in the URL.
func EncryptLink(link *messengertypes.WeshLink, passphrase []byte) (*messengertypes.WeshLink, error) {
	if link == nil {
		return nil, errcode.ErrCode_ErrInvalidInput
	}
	encrypted := messengertypes.WeshLink{
		Kind: messengertypes.WeshLink_EncryptedV1Kind,
		Encrypted: &messengertypes.WeshLink_Encrypted{
			Kind: link.Kind, // inherit kind from the clear link.
		},
	}
	if link.Encrypted != nil {
		// if display name is set in the encrypted part of the input link,
		// then we want the display name to be available in the generated URL.
		encrypted.Encrypted.DisplayName = link.Encrypted.DisplayName
		// if checksum is set, it will be the size of the hash (SHAKE256).
		encrypted.Encrypted.Checksum = link.Encrypted.Checksum
	}
	if encrypted.Encrypted.Checksum == nil {
		encrypted.Encrypted.Checksum = make([]byte, DefaultChecksumSize)
	}

	// generate nonce with AES' blocksize
	nonce, err := cryptoutil.GenerateNonceSize(aes.BlockSize)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}
	encrypted.Encrypted.Nonce = nonce

	// derive key for AES
	key, _, err := cryptoutil.DeriveKey(passphrase, nonce)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	// encrypt
	stream, err := cryptoutil.AESCTRStream(key, nonce)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	switch link.Kind {
	case messengertypes.WeshLink_ContactInviteV1Kind:
		if link.WeshId == nil {
			return nil, errcode.ErrCode_ErrInvalidInput
		}

		// encrypt fields (order is important)
		encrypted.Encrypted.ContactPublicRendezvousSeed = make([]byte, len(link.WeshId.PublicRendezvousSeed))
		encrypted.Encrypted.ContactAccountPk = make([]byte, len(link.WeshId.AccountPk))
		stream.XORKeyStream(encrypted.Encrypted.ContactPublicRendezvousSeed, link.WeshId.PublicRendezvousSeed)
		stream.XORKeyStream(encrypted.Encrypted.ContactAccountPk, link.WeshId.AccountPk)
		encrypted.Encrypted.DisplayName = link.WeshId.DisplayName

	case messengertypes.WeshLink_GroupV1Kind:
		if link.WeshGroup == nil || link.WeshGroup.Group == nil {
			return nil, errcode.ErrCode_ErrInvalidInput
		}
		// only field that stays clear
		encrypted.Encrypted.GroupType = link.WeshGroup.Group.GroupType

		// encrypt fields (order is important)
		encrypted.Encrypted.GroupPublicKey = make([]byte, len(link.WeshGroup.Group.PublicKey))
		encrypted.Encrypted.GroupSecret = make([]byte, len(link.WeshGroup.Group.Secret))
		encrypted.Encrypted.GroupSecretSig = make([]byte, len(link.WeshGroup.Group.SecretSig))
		encrypted.Encrypted.GroupSignPub = make([]byte, len(link.WeshGroup.Group.SignPub))
		encrypted.Encrypted.GroupLinkKeySig = make([]byte, len(link.WeshGroup.Group.LinkKeySig))
		stream.XORKeyStream(encrypted.Encrypted.GroupPublicKey, link.WeshGroup.Group.PublicKey)
		stream.XORKeyStream(encrypted.Encrypted.GroupSecret, link.WeshGroup.Group.Secret)
		stream.XORKeyStream(encrypted.Encrypted.GroupSecretSig, link.WeshGroup.Group.SecretSig)
		stream.XORKeyStream(encrypted.Encrypted.GroupSignPub, link.WeshGroup.Group.SignPub)
		stream.XORKeyStream(encrypted.Encrypted.GroupLinkKeySig, link.WeshGroup.Group.LinkKeySig)
		encrypted.Encrypted.DisplayName = link.WeshGroup.DisplayName

	default:
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	if len(encrypted.Encrypted.Checksum) > 0 {
		err := clearLinkChecksum(link, encrypted.Encrypted.Checksum)
		if err != nil {
			return nil, errcode.ErrCode_ErrInternal.Wrap(err)
		}
	}

	return &encrypted, nil
}

func clearLinkChecksum(link *messengertypes.WeshLink, dest []byte) error {
	// compute the hash based on every field that can be encrypted, even if they're empty.
	// the order is important.
	hasher := sha3.NewShake256()

	// contact v1
	if link.WeshId != nil {
		for _, b := range [][]byte{
			link.WeshId.PublicRendezvousSeed,
			link.WeshId.AccountPk,
		} {
			_, err := hasher.Write(b)
			if err != nil {
				return errcode.ErrCode_ErrInternal.Wrap(err)
			}
		}
	}

	// group v1
	if link.WeshGroup != nil && link.WeshGroup.Group != nil {
		for _, b := range [][]byte{
			link.WeshGroup.Group.PublicKey,
			link.WeshGroup.Group.Secret,
			link.WeshGroup.Group.SecretSig,
			link.WeshGroup.Group.SignPub,
		} {
			_, err := hasher.Write(b)
			if err != nil {
				return errcode.ErrCode_ErrInternal.Wrap(err)
			}
		}
	}

	_, err := hasher.Read(dest)
	if err != nil {
		return errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return nil
}

func InternalLinkToMessage(accountID, groupPK, cid string) (string, error) {
	if accountID == "" {
		return "", errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("account id should not be empty"))
	}

	if groupPK == "" {
		return "", errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("group pk should not be empty"))
	}

	if cid == "" {
		return "", errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("message cid should not be empty"))
	}

	internal, _, err := MarshalLink(&messengertypes.WeshLink{
		WeshMessageRef: &messengertypes.WeshLink_WeshMessageRef{
			AccountId: accountID,
			GroupPk:   groupPK,
			MessageId: cid,
		},
		Kind: messengertypes.WeshLink_MessageV1Kind,
	})
	if err != nil {
		return "", errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	return internal, nil
}

const (
	LinkWebPrefix       = "https://tone.onl/id#"
	LinkInternalPrefix  = "WESH://"
	DefaultChecksumSize = 1 // 1-byte length by default (should have ~1/256 false-positive in case of invalid password)
)

// from https://www.swisseduc.ch/informatik/theoretische_informatik/qr_codes/docs/qr_standard.pdf
//
// Alphanumeric Mode encodes data from a set of 45 characters, i.e.
// - 10 numeric digits (0 - 9) (ASCII values 30 to 39),
// - 26 alphabetic characters (A - Z) (ASCII values 41 to 5A),
// - and 9 symbols (SP, $, %, *, +, -, ., /, :) (ASCII values 20, 24, 25, 2A, 2B, 2D to 2F, 3A).
//
// we remove SP, %, +, which changes when passed through url.Encode.
//
// the generated string is longer than a base58 one, but the generated QR code is smaller which is best for scanning.
var qrBaseEncoder, _ = basex.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$*-.:/")
