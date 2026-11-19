package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"moul.io/u"

	"github.com/core-laboratories/weshnet/v2/pkg/logutil"
	"github.com/core-laboratories/weshnet/v2/pkg/testutil"
)

func TestPersistentIdentity(t *testing.T) {
	testutil.FilterStability(t, testutil.Flappy)

	if runtime.GOOS == "windows" {
		t.Skip("disabled on windows")
	}

	// @FIXME(gfanton): avoid write error caused by grpc logger
	logutil.ReplaceGRPCLogger(zap.NewNop())

	// create tempdir for the test
	var tempdir string
	{
		var err error
		tempdir, err = os.MkdirTemp("", "wesh-main")
		require.NoError(t, err)
		defer os.RemoveAll(tempdir)
	}
	var (
		pathWesh1 = filepath.Join(tempdir, "wesh1")
		pathWesh2 = filepath.Join(tempdir, "wesh2")
		pathExport = filepath.Join(tempdir, "export.tar")
	)

	// wesh1: init a new account
	var key1 string
	{
		closer, err := u.CaptureStdout()
		require.NoError(t, err)
		err = runMain([]string{
			"share-invite",
			"-log.filters=none",
			"-store.dir", pathWesh1,
			"-no-qr",
		})
		require.NoError(t, err)
		key1 = strings.TrimSpace(closer())
		require.NotEmpty(t, key1)
		require.Contains(t, key1, "https://tone.onl/id")
	}

	// wesh1: export account
	{
		err := runMain([]string{
			"export",
			"-log.filters=none",
			"-store.dir", pathWesh1,
			"-export-path", pathExport,
		})
		require.NoError(t, err)
	}

	// wesh2: init a new account from export
	{
		closer, err := u.CaptureStdout()
		require.NoError(t, err)
		err = runMain([]string{
			"share-invite",
			"-log.filters=none",
			"-store.dir", pathWesh2,
			"-node.restore-export-path", pathExport,
			"-no-qr",
		})
		require.NoError(t, err)
		key2 := strings.TrimSpace(closer())
		require.NotEmpty(t, key2)
		require.Equal(t, key1, key2)
	}
}
