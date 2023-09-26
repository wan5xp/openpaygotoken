package simulators

import (
	"math"
	"time"

	"github.com/wan5xp/openpaygotoken/pkg/openpaygotoken"
)

// ErrTooManyDays is returned when the number of days is too high.
type ErrTooManyDays struct {
}

func (e *ErrTooManyDays) Error() string {
	return "Too many days"
}

// SingleDeviceServerSimulator is a simulator for a single device server.
type SingleDeviceServerSimulator struct {
	StartingCode           int
	Key                    [16]byte
	Count                  int
	ExpirationDate         time.Time
	FurthestExpirationDate time.Time
	PaygEnabled            bool
	TimeDivider            int
	RestrictedDigitSet     bool
}

// NewSingleDeviceServerSimulator creates a new SingleDeviceServerSimulator.
func NewSingleDeviceServerSimulator(startingCode int, key *[16]byte, startingCount int, restrictedDigitSet bool, timeDivider int) *SingleDeviceServerSimulator {
	return &SingleDeviceServerSimulator{
		StartingCode:           startingCode,
		Key:                    *key,
		Count:                  startingCount,
		PaygEnabled:            true,
		TimeDivider:            timeDivider,
		RestrictedDigitSet:     restrictedDigitSet,
		ExpirationDate:         time.Now(),
		FurthestExpirationDate: time.Now(),
	}
}

// GeneratePaygDisableToken generates a PAYG disable token.
func (s *SingleDeviceServerSimulator) GeneratePaygDisableToken() (string, error) {
	count, token, err := openpaygotoken.GenerateStandardToken(s.StartingCode, &s.Key, openpaygotoken.PAYGDisableValue, s.Count, openpaygotoken.SetTime, s.RestrictedDigitSet)
	if err != nil {
		return "", err
	}
	s.Count = count
	return token, nil
}

// GenerateTokenFromDate generates a token from a date
func (s *SingleDeviceServerSimulator) GenerateTokenFromDate(newExpirationDate time.Time, force bool) (string, error) {
	var value int
	var err error
	furthestExpirationDate := s.FurthestExpirationDate
	if newExpirationDate.After(s.FurthestExpirationDate) {
		s.FurthestExpirationDate = newExpirationDate
	}
	if newExpirationDate.After(furthestExpirationDate) {
		value, err = s.getValueToActivate(newExpirationDate, s.ExpirationDate, force)
		if err != nil {
			return "", err
		}
		s.ExpirationDate = newExpirationDate
		return s.GenerateTokenFromValue(value, openpaygotoken.AddTime)
	} else {
		value, err = s.getValueToActivate(newExpirationDate, time.Now(), force)
		if err != nil {
			return "", err
		}
		s.ExpirationDate = newExpirationDate
		return s.GenerateTokenFromValue(value, openpaygotoken.SetTime)
	}
}

// GenerateTokenFromValue generates a token from a value
func (s *SingleDeviceServerSimulator) GenerateTokenFromValue(value int, mode openpaygotoken.TokenType) (string, error) {
	count, token, err := openpaygotoken.GenerateStandardToken(s.StartingCode, &s.Key, value, s.Count, mode, s.RestrictedDigitSet)
	if err != nil {
		return "", err
	}
	s.Count = count
	return token, nil
}

// GetValueToActivate returns the value to activate.
func (s *SingleDeviceServerSimulator) getValueToActivate(newTime time.Time, referenceTime time.Time, forceMaximum bool) (int, error) {
	if !newTime.After(referenceTime) {
		return 0, nil
	} else {
		days := math.Round(newTime.Sub(referenceTime).Hours() / 24)
		value := int(days) * s.TimeDivider
		if value > openpaygotoken.MaxActivationValue {
			if !forceMaximum {
				return 0, &ErrTooManyDays{}
			} else {
				return openpaygotoken.MaxActivationValue, nil
			}
		}
		return value, nil
	}
}
