package weshmessengertesting

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/core-laboratories/weshnet/mobile/v2/go/internal/initutil"
	"github.com/core-laboratories/weshnet/mobile/v2/go/pkg/messengertypes"
	"github.com/core-laboratories/weshnet/v2"
	"github.com/core-laboratories/weshnet/v2/pkg/protocoltypes"
)

func NonMockedTestingInfra(t testing.TB, accountAmount int) ([]messengertypes.MessengerServiceClient, []*weshnet.TestingProtocol) {
	messengers := make([]messengertypes.MessengerServiceClient, accountAmount)
	tps := make([]*weshnet.TestingProtocol, accountAmount)

	for i := 0; i < accountAmount; i++ {
		tempDir, err := os.MkdirTemp("", fmt.Sprintf("wesh-main-%d", i))

		require.NoError(t, err)
		t.Cleanup(func() {
			_ = os.RemoveAll(tempDir)
		})

		man, err := initutil.New(nil)
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = man.Close(nil)
		})

		fs := flag.NewFlagSet("man1", flag.ExitOnError)
		man.SetupLoggingFlags(fs)              // also available at root level
		man.SetupMetricsFlags(fs)              // add flags to enable metrics
		man.SetupLocalMessengerServerFlags(fs) // add flags to allow creating a full node in the same process
		man.SetupEmptyGRPCListenersFlags(fs)   // by default, we don't want to expose gRPC server for mini
		man.SetupRemoteNodeFlags(fs)           // mini can be run against an already running server
		man.SetupInitTimeout(fs)

		err = fs.Parse([]string{"-store.dir", tempDir})
		require.NoError(t, err)

		_, _, err = man.GetLocalIPFS()
		require.NoError(t, err)

		protocolClient, err := man.GetProtocolClient()
		require.NoError(t, err)

		messengers[i], err = man.GetMessengerClient()
		require.NotNil(t, messengers[i])
		require.NoError(t, err)

		tps[i] = &weshnet.TestingProtocol{
			Client: &WrappedMessengerClient{protocolClient},
		}
	}

	return messengers, tps
}

type WrappedMessengerClient struct {
	protocoltypes.ProtocolServiceClient
}

func (*WrappedMessengerClient) Close() error {
	return nil
}
