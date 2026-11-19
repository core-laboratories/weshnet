package rendezvous

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAdditionalRotationFallsBackWhenRemoved(t *testing.T) {
	rotations := NewRotationInterval(time.Hour)
	now := time.Now()
	defaultSeed := []byte("default-seed")
	privateSeed := []byte("private-seed")

	rotations.RegisterRotation(now, "group-topic", defaultSeed)
	rotations.RegisterRotation(now, "group-topic", privateSeed)

	point, err := rotations.PointForTopic("group-topic")
	require.NoError(t, err)
	require.Equal(t, privateSeed, point.Seed())
	require.True(t, rotations.UnregisterRotation("group-topic", privateSeed))

	point, err = rotations.PointForTopic("group-topic")
	require.NoError(t, err)
	require.Equal(t, defaultSeed, point.Seed())
}

func TestUnregisterRotationRejectsUnknownSeed(t *testing.T) {
	rotations := NewRotationInterval(time.Hour)
	rotations.RegisterRotation(time.Now(), "group-topic", []byte("default-seed"))

	require.False(t, rotations.UnregisterRotation("group-topic", []byte("unknown-seed")))
	point, err := rotations.PointForTopic("group-topic")
	require.NoError(t, err)
	require.Equal(t, []byte("default-seed"), point.Seed())
}
