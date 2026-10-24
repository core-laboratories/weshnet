//go:build android && !noproximitytransport

package androidnearby

import (
	"go.uber.org/zap"

	proximity "github.com/core-laboratories/weshnet/v2/pkg/proximitytransport"
)

// Supported is used by main package as default value for enable this driver.
// While UI actually enable or not the Java Android Nearby driver.
// TODO: remove this when UI will be able to handle this for the first App launching.
const Supported = true

// Noop implementation for Android
// The host Android application supplies the real driver through this bridge.
func NewDriver(logger *zap.Logger) proximity.ProximityDriver {
	logger = logger.Named("Nearby")
	logger.Info("NewDriver(): Java driver not found")

	return proximity.NewNoopProximityDriver(ProtocolCode, ProtocolName, DefaultAddr)
}
