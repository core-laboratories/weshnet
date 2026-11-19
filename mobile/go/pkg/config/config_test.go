package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenWorks(t *testing.T) {
	// we have at least welcomebot and testbot, and each contact has a link
	require.NotNil(t, Config.Wesh.Contacts)
	require.True(t, len(Config.Wesh.Contacts) >= 1)
	for _, contact := range []string{"welcomebot", "testbot"} {
		require.NotNil(t, Config.Wesh.Contacts[contact])
	}
	for idx := range Config.Wesh.Contacts {
		require.NotEmpty(t, Config.Wesh.Contacts[idx].Link)
	}

	// we have at least one rdvp and each one has a maddr
	require.NotNil(t, Config.P2P.RDVP)
	require.True(t, len(Config.P2P.RDVP) >= 1)
	for _, peer := range Config.P2P.RDVP {
		require.NotEmpty(t, peer.Maddr)
	}
}
