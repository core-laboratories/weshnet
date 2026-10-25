package rendezvous

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sync"
	"time"
)

var (
	RotationGracePeriod  = time.Minute * 10
	MinimumDelayRotation = time.Minute
)

type RotationInterval struct {
	interval time.Duration

	cacheTopics    map[string][]*Point
	cacheRotations map[string]*Point
	muCache        sync.RWMutex
}

func NewStaticRotationInterval() *RotationInterval {
	// from https://stackoverflow.com/a/32620397
	maxTime := time.Unix(1<<63-62135596801, 999999999)
	return NewRotationInterval(time.Until(maxTime))
}

func NewRotationInterval(interval time.Duration) *RotationInterval {
	return &RotationInterval{
		interval:       interval,
		cacheTopics:    make(map[string][]*Point),
		cacheRotations: make(map[string]*Point),
	}
}

func (r *RotationInterval) RegisterRotation(at time.Time, topic string, seed []byte) {
	point := r.NewRendezvousPointForPeriod(at, topic, append([]byte(nil), seed...))
	r.muCache.Lock()
	r.removeSeedLocked(topic, seed)
	r.registerPoint(point)
	r.muCache.Unlock()
}

// UnregisterRotation removes a seed from a topic. If multiple seeds are
// registered, PointForTopic falls back to the previously registered seed.
// Existing rotation aliases are removed immediately so retired rendezvous
// points can no longer be used to decode new envelopes.
func (r *RotationInterval) UnregisterRotation(topic string, seed []byte) bool {
	r.muCache.Lock()
	defer r.muCache.Unlock()

	return r.removeSeedLocked(topic, seed)
}

func (r *RotationInterval) RoundTimePeriod(at time.Time) time.Time {
	return RoundTimePeriod(at, r.interval)
}

func (r *RotationInterval) NextTimePeriod(at time.Time) time.Time {
	return NextTimePeriod(at, r.interval)
}

func (r *RotationInterval) PointForRawRotation(rotation []byte) (*Point, error) {
	return r.PointForRotation(base64.StdEncoding.EncodeToString(rotation))
}

func (r *RotationInterval) PointForRotation(rotation string) (*Point, error) {
	r.muCache.Lock()
	defer r.muCache.Unlock()

	if point, ok := r.cacheRotations[rotation]; ok {
		if point.IsExpired() {
			point = r.rotate(point, DefaultRotationInterval)
		}

		return point, nil
	}

	return nil, fmt.Errorf("unable to get rendez vous, no matching rotation registered")
}

func (r *RotationInterval) PointForTopic(topic string) (*Point, error) {
	r.muCache.Lock()
	defer r.muCache.Unlock()

	if points := r.cacheTopics[topic]; len(points) > 0 {
		point := points[len(points)-1]
		if point.IsExpired() {
			point = r.rotate(point, DefaultRotationInterval)
		}

		return point, nil
	}

	return nil, fmt.Errorf("unable to get rendezvous point, no matching topic registered")
}

func (r *RotationInterval) NewRendezvousPointForPeriod(at time.Time, topic string, seed []byte) (point *Point) {
	at = r.RoundTimePeriod(at)
	rotation := GenerateRendezvousPointForPeriod([]byte(topic), seed, at)

	next := r.NextTimePeriod(at)
	return &Point{
		rp:       r,
		rotation: rotation,
		topic:    topic,
		seed:     seed,
		deadline: next,
	}
}

func (r *RotationInterval) registerPoint(point *Point) {
	keytopic, keyrotation := point.keys()
	r.cacheTopics[keytopic] = append(r.cacheTopics[keytopic], point)
	r.cacheRotations[keyrotation] = point
}

func (r *RotationInterval) removeSeedLocked(topic string, seed []byte) bool {
	points := r.cacheTopics[topic]
	removed := false
	kept := points[:0]
	for _, point := range points {
		if bytes.Equal(point.seed, seed) {
			_, rotation := point.keys()
			delete(r.cacheRotations, rotation)
			removed = true
			continue
		}
		kept = append(kept, point)
	}
	for rotation, point := range r.cacheRotations {
		if point.topic == topic && bytes.Equal(point.seed, seed) {
			delete(r.cacheRotations, rotation)
			removed = true
		}
	}
	if len(kept) == 0 {
		delete(r.cacheTopics, topic)
	} else {
		r.cacheTopics[topic] = kept
	}
	return removed
}

func (r *RotationInterval) rotate(old *Point, graceperiod time.Duration) *Point {
	newPoint := old.NextPoint()

	points := r.cacheTopics[old.topic]
	for i, point := range points {
		if point == old {
			points[i] = newPoint
			break
		}
	}
	r.cacheTopics[old.topic] = points
	_, newRotation := newPoint.keys()
	r.cacheRotations[newRotation] = newPoint

	cleanuptime := max(time.Until(newPoint.Deadline().Add(graceperiod)), 0)
	// cleanup after the grace period
	time.AfterFunc(cleanuptime, func() {
		r.muCache.Lock()
		_, keyrotation := old.keys()
		delete(r.cacheRotations, keyrotation)
		r.muCache.Unlock()
	})

	return newPoint
}

type Point struct {
	rp       *RotationInterval
	topic    string
	rotation []byte
	seed     []byte
	deadline time.Time
}

func (p *Point) NextPoint() *Point {
	if p.IsExpired() {
		return p.rp.NewRendezvousPointForPeriod(time.Now(), p.topic, p.seed)
	}

	return p.rp.NewRendezvousPointForPeriod(p.deadline.Add(time.Second), p.topic, p.seed)
}

func (p *Point) Seed() []byte {
	return p.seed
}

func (p *Point) keys() (topic string, rotation string) {
	topic = p.Topic()
	rotation = p.RotationTopic()
	return
}

func (p *Point) Topic() string {
	return p.topic
}

func (p *Point) RawTopic() []byte {
	return []byte(p.topic)
}

func (p *Point) RotationTopic() string {
	return base64.StdEncoding.EncodeToString(p.rotation)
}

func (p *Point) RawRotationTopic() []byte {
	return p.rotation
}

func (p *Point) Deadline() time.Time {
	return p.deadline
}

func (p *Point) TTL() time.Duration {
	return time.Until(p.deadline)
}

func (p *Point) IsExpired() bool {
	return p.TTL() <= 0
}
