package throttle

import "time"

const (
	defaultInitialSleep     = 100 * time.Millisecond
	defaultMinSleep         = 10 * time.Millisecond
	defaultMaxSleep         = 10 * time.Second
	defaultIncreaseModifier = float64(1.2)
	defaultDecreaseModifier = float64(0.8)
)

type Config struct {
	InitialSleep     time.Duration
	MinSleep         time.Duration
	MaxSleep         time.Duration
	IncreaseModifier float64
	DecreaseModifier float64
}

type Throttler struct {
	currentSleep     time.Duration
	minSleep         time.Duration
	maxSleep         time.Duration
	increaseModifier float64
	decreaseModifier float64
}

// NewThrottler creates a Throttler with the given Config.
// If a configuration is not set, a default will be used.
func NewThrottler(config Config) *Throttler {
	initialSleep := defaultInitialSleep
	if config.InitialSleep > 0 { initialSleep = config.InitialSleep }

	minSleep := defaultMinSleep
	if config.MinSleep > 0 { minSleep = config.MinSleep }

	maxSleep := defaultMaxSleep
	if config.MaxSleep > 0 { maxSleep = config.MaxSleep }

	increaseModifier := defaultIncreaseModifier
	if config.IncreaseModifier > 0 { increaseModifier = config.IncreaseModifier }

	decreaseModifier := defaultDecreaseModifier
	if config.DecreaseModifier > 0 { decreaseModifier = config.DecreaseModifier }

	return &Throttler{
		currentSleep:     initialSleep,
		minSleep:         minSleep,
		maxSleep:         maxSleep,
		increaseModifier: increaseModifier,
		decreaseModifier: decreaseModifier,
	}
}

// Increase the currentSleep duration by the increaseModifier percentage.
// currentSleep cannot be higher than maxSleep
func (t *Throttler) Increase() time.Duration {
	sleep := time.Duration(float64(t.currentSleep) * t.increaseModifier)
	if sleep > t.maxSleep {
		t.currentSleep = t.maxSleep
	} else {
		t.currentSleep = sleep
	}
	return t.currentSleep
}

// Decrease the currentSleep duration by the decreaseModifier percentage.
// currentSleep cannot be lower than minSleep
func (t *Throttler) Decrease() time.Duration {
	sleep := time.Duration(float64(t.currentSleep) * t.decreaseModifier)
	if sleep < t.minSleep {
		t.currentSleep = t.minSleep
	} else {
		t.currentSleep = sleep
	}
	return t.currentSleep
}

// Reset currentSleep to the minSleep duration
func (t *Throttler) Reset() time.Duration {
	t.currentSleep = t.minSleep
	return t.currentSleep
}

// Duration returns the currentSleep duration.
func (t *Throttler) Duration() time.Duration {
	return t.currentSleep
}
