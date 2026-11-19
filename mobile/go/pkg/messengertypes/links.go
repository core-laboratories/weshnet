package messengertypes

import (
	fmt "fmt"

	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/errcode"
	protocoltypes "github.com/core-laboratories/weshnet/v2/pkg/protocoltypes"
)

func (link *WeshLink) IsValid() error {
	if link == nil {
		return errcode.ErrCode_ErrMissingInput
	}
	switch link.Kind {
	case WeshLink_ContactInviteV1Kind:
		if link.WeshId == nil ||
			link.WeshId.AccountPk == nil ||
			link.WeshId.PublicRendezvousSeed == nil {
			return errcode.ErrCode_ErrMissingInput
		}
		return nil

	case WeshLink_GroupV1Kind:
		if link.WeshGroup == nil {
			return errcode.ErrCode_ErrMissingInput
		}
		if groupType := link.WeshGroup.Group.GroupType; groupType != protocoltypes.GroupType_GroupTypeMultiMember {
			return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("can't share a %q group type", groupType))
		}
		return nil

	case WeshLink_EncryptedV1Kind:
		if link.Encrypted == nil {
			return errcode.ErrCode_ErrMissingInput
		}
		switch link.Encrypted.Kind {
		case WeshLink_ContactInviteV1Kind:
			if link.Encrypted.ContactAccountPk == nil ||
				link.Encrypted.ContactPublicRendezvousSeed == nil {
				return errcode.ErrCode_ErrMissingInput
			}
		case WeshLink_GroupV1Kind:
			if groupType := link.Encrypted.GroupType; groupType != protocoltypes.GroupType_GroupTypeMultiMember {
				return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("can't share a %q group type", groupType))
			}
		default:
			return errcode.ErrCode_ErrInvalidInput
		}
		return nil
	case WeshLink_MessageV1Kind:
		if link.WeshMessageRef == nil {
			return errcode.ErrCode_ErrMissingInput
		}
		return nil
	}

	return errcode.ErrCode_ErrInvalidInput
}

func (link *WeshLink) IsContact() bool {
	return link.Kind == WeshLink_ContactInviteV1Kind &&
		link.IsValid() == nil
}

func (link *WeshLink) IsGroup() bool {
	return link.Kind == WeshLink_GroupV1Kind &&
		link.IsValid() == nil
}

func (id *WeshID) GetWeshLink() *WeshLink {
	return &WeshLink{
		Kind:    WeshLink_ContactInviteV1Kind,
		WeshId: id,
	}
}

func (group *WeshGroup) GetWeshLink() *WeshLink {
	return &WeshLink{
		Kind:       WeshLink_GroupV1Kind,
		WeshGroup: group,
	}
}
