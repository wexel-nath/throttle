package throttle

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewThrottler(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   *Throttler
	}{
		{
			name:   "Default Configs",
			config: Config{},
			want:   &Throttler{
				currentSleep:     defaultInitialSleep,
				minSleep:         defaultMinSleep,
				maxSleep:         defaultMaxSleep,
				increaseModifier: defaultIncreaseModifier,
				decreaseModifier: defaultDecreaseModifier,
			},
		},
		{
			name:   "Some Configs Set",
			config: Config{
				InitialSleep:     500 * time.Millisecond,
				DecreaseModifier: float64(1.05),
			},
			want:   &Throttler{
				currentSleep:     500 * time.Millisecond,
				minSleep:         defaultMinSleep,
				maxSleep:         defaultMaxSleep,
				increaseModifier: defaultIncreaseModifier,
				decreaseModifier: float64(1.05),
			},
		},
		{
			name:   "All Configs Set",
			config: Config{
				InitialSleep:     500 * time.Millisecond,
				MinSleep:         100 * time.Millisecond,
				MaxSleep:         5 * time.Second,
				IncreaseModifier: float64(1.1),
				DecreaseModifier: float64(1.05),
			},
			want:   &Throttler{
				currentSleep:     500 * time.Millisecond,
				minSleep:         100 * time.Millisecond,
				maxSleep:         5 * time.Second,
				increaseModifier: float64(1.1),
				decreaseModifier: float64(1.05),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			got := NewThrottler(test.config)
			assert.Equal(st, test.want, got)
		})
	}
}

func TestThrottler_Increase(t *testing.T) {
	tests := []struct {
		name     string
		initial  time.Duration
		max      time.Duration
		modifier float64
		want     time.Duration
	}{
		{
			name:     "Default Configs",
			initial:  defaultInitialSleep,
			max:      defaultMaxSleep,
			modifier: defaultIncreaseModifier,
			want:     120 * time.Millisecond,
		},
		{
			name:     "50ms Initial Sleep 1.5 Increase Modifier",
			initial:  50 * time.Millisecond,
			max:      defaultMaxSleep,
			modifier: float64(1.5),
			want:     75 * time.Millisecond,
		},
		{
			name:     "Increase Above Max",
			initial:  800 * time.Millisecond,
			max:      time.Second,
			modifier: float64(1.5),
			want:     time.Second,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			throttler := NewThrottler(Config{
				InitialSleep:     test.initial,
				MaxSleep:         test.max,
				IncreaseModifier: test.modifier,
			})

			got := throttler.Increase()
			assert.Equal(st, test.want, got)
			assert.Equal(st, test.want, throttler.Duration())
		})
	}
}

func TestThrottler_Decrease(t *testing.T) {
	tests := []struct {
		name     string
		initial  time.Duration
		min      time.Duration
		modifier float64
		want     time.Duration
	}{
		{
			name:     "Default Configs",
			initial:  defaultInitialSleep,
			min:      defaultMinSleep,
			modifier: defaultDecreaseModifier,
			want:     80 * time.Millisecond,
		},
		{
			name:     "160ms InitialSleep 0.75 Increase Modifier",
			initial:  160 * time.Millisecond,
			min:      defaultMinSleep,
			modifier: float64(0.75),
			want:     120 * time.Millisecond,
		},
		{
			name:     "Decrease Below Min",
			initial:  125 * time.Millisecond,
			min:      100 * time.Millisecond,
			modifier: float64(0.6),
			want:     100 * time.Millisecond,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			throttler := NewThrottler(Config{
				InitialSleep:     test.initial,
				MinSleep:         test.min,
				DecreaseModifier: test.modifier,
			})

			got := throttler.Decrease()
			assert.Equal(st, test.want, got)
			assert.Equal(st, test.want, throttler.Duration())
		})
	}
}

func TestThrottler_Reset(t *testing.T) {
	tests := []struct {
		name string
		min  time.Duration
		want time.Duration
	}{
		{
			name: "Default Configs",
			min:  defaultMinSleep,
			want: 10 * time.Millisecond,
		},
		{
			name: "1s MinSleep",
			min:  1 * time.Second,
			want: 1 * time.Second,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			throttler := NewThrottler(Config{
				MinSleep: test.min,
			})

			got := throttler.Reset()
			assert.Equal(st, test.want, got)
			assert.Equal(st, test.want, throttler.Duration())
		})
	}
}

func TestThrottler(t *testing.T) {
	throttler := NewThrottler(Config{
		InitialSleep:     200 * time.Millisecond,
		MinSleep:         20 * time.Millisecond,
		MaxSleep:         time.Second,
		IncreaseModifier: float64(1.5),
		DecreaseModifier: float64(0.8),
	})

	throttler.Increase()
	throttler.Increase()

	assert.Equal(t, 450 * time.Millisecond, throttler.Duration())

	throttler.Decrease()

	assert.Equal(t, 360 * time.Millisecond, throttler.Duration())

	throttler.Increase()
	throttler.Increase()
	throttler.Increase()

	assert.Equal(t, time.Second, throttler.Duration())

	throttler.Decrease()

	assert.Equal(t, 800 * time.Millisecond, throttler.Duration())
}
