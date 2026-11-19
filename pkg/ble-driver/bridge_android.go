//go:build android && !noproximitytransport

package ble

import (
	"go.uber.org/zap"

	proximity "github.com/core-laboratories/weshnet/v2/pkg/proximitytransport"
)

// Supported is used by main package as default value for enable the BLE  driver.
// While UI actually enable or not the Java BLE driver.
// TODO: remove this when UI will be able to handle this for the first App launching.
const Supported = true

// Noop implementation for Android
// The real driver is provided by the embedding Android application.
func NewDriver(logger *zap.Logger) proximity.ProximityDriver {
	logger = logger.Named("BLE")
	logger.Info("NewDriver(): Java driver not found")

	return proximity.NewNoopProximityDriver(ProtocolCode, ProtocolName, DefaultAddr)
}
