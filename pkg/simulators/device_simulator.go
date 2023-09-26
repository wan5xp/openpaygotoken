package simulators

import (
	"fmt"
	"strconv"
	"time"

	"github.com/wan5xp/openpaygotoken/pkg/openpaygotoken"
)

// ErrTokenEntryBlocked is returned when the token entry is blocked.
type ErrTokenEntryBlocked struct {
}

func (e *ErrTokenEntryBlocked) Error() string {
	return "Token entry blocked"
}

// ErrOldToken is returned when the token is old.
type ErrOldToken struct {
}

func (e *ErrOldToken) Error() string {
	return "Old Token"
}

// DeviceSimulator is a simulator for a device.
type DeviceSimulator struct {
	StartingCode           int
	Key                    [16]byte
	TimeDivider            int
	RestrictedDigitSet     bool
	WaitingPeriodEnabled   bool
	PaygEnabled            bool
	Count                  int
	ExpirationTimestamp    time.Time
	InvalidTokenCount      int
	TokenEntryBlockedUntil time.Time
	UsedCounts             []int
	decoder                *openpaygotoken.TokenDecoder
}

// EnterToken enters a token in the device.
func (d *DeviceSimulator) EnterToken(token string) error {
	if len(token) == 9 {
		tokenInt, err := strconv.Atoi(token)
		if err != nil {
			return err
		}
		return d.updateDeviceStatusFromToken(tokenInt)
	} else {
		tokenInt, err := strconv.Atoi(token)
		if err != nil {
			return err
		}
		return d.updateDeviceStatusFromToken(tokenInt)
	}
}

// updateDeviceStatusFromToken updates the device status from a token.
func (d *DeviceSimulator) updateDeviceStatusFromToken(token int) error {
	if d.TokenEntryBlockedUntil.After(time.Now()) && d.WaitingPeriodEnabled {
		return &ErrTokenEntryBlocked{}
	}
	value, count, tokenType, err := d.decoder.GetActivationValueCountAndTypeFromToken(token, d.StartingCode, &d.Key, d.Count, d.RestrictedDigitSet, &d.UsedCounts)
	if err != nil {
		d.InvalidTokenCount++
		d.TokenEntryBlockedUntil = time.Now().Add(2 * time.Minute)
		for xn := 0; xn < d.InvalidTokenCount-1; xn++ {
			d.TokenEntryBlockedUntil = d.TokenEntryBlockedUntil.Add(time.Duration(2*d.InvalidTokenCount) * time.Minute)
		}
		return err
	} else if value == -2 {
		return &ErrOldToken{}
	} else {
		if count > d.Count || value == openpaygotoken.CounterSyncValue {
			d.Count = count
		}
		d.UsedCounts = d.decoder.UpdateUsedCounts(&d.UsedCounts, value, count, tokenType)
		d.InvalidTokenCount = 0
		if value <= openpaygotoken.MaxActivationValue {
			if !d.PaygEnabled && tokenType == openpaygotoken.SetTime {
				d.PaygEnabled = true
			}
			if d.PaygEnabled {
				if tokenType == openpaygotoken.SetTime {
					d.ExpirationTimestamp = time.Now().Add(time.Duration(value/d.TimeDivider) * 24 * time.Hour)
				} else {
					d.ExpirationTimestamp = d.ExpirationTimestamp.Add(time.Duration(value/d.TimeDivider) * 24 * time.Hour)
				}
			}
		} else if value == openpaygotoken.PAYGDisableValue {
			d.PaygEnabled = false
			// } else if value != CounterSyncValue {
			// 	// Unknown
			// } else {
			// 	// Unknown
		}
	}
	return nil
}

// IsActive returns true if the device is active.
func (d *DeviceSimulator) IsActive() bool {
	return time.Now().Before(d.ExpirationTimestamp)
}

// NewDeviceSimulator creates a new device simulator.
func NewDeviceSimulator(startingCode int, key *[16]byte, startingCount int, restrictedDigit bool, waitingPeriodEnabled bool, timeDivider int) (*DeviceSimulator, error) {
	decoder, err := openpaygotoken.NewDecoder()

	if err != nil {
		return nil, err
	}

	return &DeviceSimulator{
		StartingCode:           startingCode,
		Key:                    *key,
		Count:                  startingCount,
		RestrictedDigitSet:     restrictedDigit,
		WaitingPeriodEnabled:   waitingPeriodEnabled,
		TimeDivider:            timeDivider,
		decoder:                decoder,
		ExpirationTimestamp:    time.Now(),
		InvalidTokenCount:      0,
		UsedCounts:             make([]int, 0),
		TokenEntryBlockedUntil: time.Now(),
		PaygEnabled:            true,
	}, nil
}

// PrintStatus prints the status of the device.
func (d *DeviceSimulator) PrintStatus() {
	fmt.Println("-------------------------")
	fmt.Println("Expiration Date:", d.ExpirationTimestamp)
	fmt.Println("Current count:", d.Count)
	fmt.Println("PAYG Enabled:", d.PaygEnabled)
	fmt.Println("Active:", d.IsActive())
	fmt.Println("-------------------------")
}
