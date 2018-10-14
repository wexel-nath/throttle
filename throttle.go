package throttle

import "time"

const (
	defaultInitialSleep = 100 * time.Millisecond
	defaultMinSleep     = time.Duration(0)
	defaultMaxSleep     = 10 * time.Second
	defaultModifier     = float64(1.2)
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

	increaseModifier := defaultModifier
	if config.IncreaseModifier > 0 { increaseModifier = config.IncreaseModifier }

	decreaseModifier := defaultModifier
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
func (t *Throttler) Increase() {
	sleep := time.Duration(float64(t.currentSleep) * t.increaseModifier)
	if sleep > t.maxSleep {
		t.currentSleep = t.maxSleep
	} else {
		t.currentSleep = sleep
	}
}

// Decrease the currentSleep duration by the decreaseModifier percentage.
// currentSleep cannot be lower than minSleep
func (t *Throttler) Decrease() {
	sleep := time.Duration(float64(t.currentSleep) / t.decreaseModifier)
	if sleep < t.minSleep {
		t.currentSleep = t.minSleep
	} else {
		t.currentSleep = sleep
	}
}

// Wait sleeps and returns the sleep duration.
func (t *Throttler) Wait() time.Duration {
	time.Sleep(t.currentSleep)
	return t.currentSleep
}
