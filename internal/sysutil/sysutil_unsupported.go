//go:build !linux && !darwin

package sysutil

import "github.com/core-laboratories/weshnet/v2/pkg/protocoltypes"

func appendCustomSystemInfo(reply *protocoltypes.SystemInfo_Process) error {
	return nil
}
