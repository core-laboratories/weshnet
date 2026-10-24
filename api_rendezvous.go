package weshnet

import (
	"context"
	"fmt"
	"time"

	"github.com/core-laboratories/weshnet/v2/pkg/errcode"
	"github.com/core-laboratories/weshnet/v2/pkg/protocoltypes"
)

func (s *service) GroupRendezvousSeedAdd(_ context.Context, req *protocoltypes.GroupRendezvousSeedAdd_Request) (*protocoltypes.GroupRendezvousSeedAdd_Reply, error) {
	group, err := s.GetContextGroupForID(req.GetGroupPk())
	if err != nil {
		return nil, errcode.ErrCode_ErrGroupMemberUnknownGroupID.Wrap(err)
	}
	if err := validateAdditionalRendezvousSeed(req.GetSeed()); err != nil {
		return nil, err
	}

	s.registerGroupRendezvousSeed(group, req.GetSeed())
	return &protocoltypes.GroupRendezvousSeedAdd_Reply{}, nil
}

func (s *service) GroupRendezvousSeedRemove(_ context.Context, req *protocoltypes.GroupRendezvousSeedRemove_Request) (*protocoltypes.GroupRendezvousSeedRemove_Reply, error) {
	group, err := s.GetContextGroupForID(req.GetGroupPk())
	if err != nil {
		return nil, errcode.ErrCode_ErrGroupMemberUnknownGroupID.Wrap(err)
	}
	if err := validateAdditionalRendezvousSeed(req.GetSeed()); err != nil {
		return nil, err
	}

	s.unregisterGroupRendezvousSeed(group, req.GetSeed())
	return &protocoltypes.GroupRendezvousSeedRemove_Reply{}, nil
}

func validateAdditionalRendezvousSeed(seed []byte) error {
	if len(seed) != protocoltypes.RendezvousSeedLength {
		return errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("rendezvous seed must be %d bytes", protocoltypes.RendezvousSeedLength))
	}
	return nil
}

func (s *service) registerGroupRendezvousSeed(group *GroupContext, seed []byte) {
	groupID := string(group.Group().GetPublicKey())
	s.rendezvousSeedsMu.Lock()
	defer s.rendezvousSeedsMu.Unlock()
	if s.rendezvousSeeds == nil {
		s.rendezvousSeeds = make(map[string][]byte)
	}
	if previous := s.rendezvousSeeds[groupID]; len(previous) != 0 {
		for _, topic := range groupRendezvousTopics(group) {
			s.odb.rotationInterval.UnregisterRotation(topic, previous)
		}
	}
	for _, topic := range groupRendezvousTopics(group) {
		s.odb.rotationInterval.RegisterRotation(time.Now(), topic, seed)
	}
	s.rendezvousSeeds[groupID] = append([]byte(nil), seed...)
}

func (s *service) unregisterGroupRendezvousSeed(group *GroupContext, seed []byte) {
	groupID := string(group.Group().GetPublicKey())
	s.rendezvousSeedsMu.Lock()
	defer s.rendezvousSeedsMu.Unlock()
	if active := s.rendezvousSeeds[groupID]; string(active) != string(seed) {
		return
	}
	for _, topic := range groupRendezvousTopics(group) {
		s.odb.rotationInterval.UnregisterRotation(topic, seed)
	}
	delete(s.rendezvousSeeds, groupID)
}

func groupRendezvousTopics(group *GroupContext) []string {
	return []string{
		group.MetadataStore().Address().String(),
		group.MessageStore().Address().String(),
	}
}
