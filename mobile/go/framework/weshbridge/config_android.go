//go:build android
// +build android

package weshbridge

import "github.com/core-laboratories/weshnet/mobile/v2/go/pkg/tempdir"

func (c *BridgeConfig) SetAndroidCacheDir(path string) { tempdir.SetAndroidCacheDir(path) }
