package openpaygotoken_test

import (
	"testing"

	"github.com/wan5xp/openpaygotoken/pkg/openpaygotoken"
	"golang.org/x/exp/slices"
)

var (
	key = [16]byte{162, 154, 184, 46, 220, 95, 187, 196, 30, 201, 83, 15, 109, 172, 134, 177}
)

const (
	startingCode = 123456789
)

func TestGenerateStandardToken(t *testing.T) {
	count, token, err := openpaygotoken.GenerateStandardToken(startingCode, &key, openpaygotoken.PAYGDisableValue, 1, openpaygotoken.SetTime, false)
	if err != nil {
		t.Error(err)
	}
	if count != 3 {
		t.Errorf("Expected count to be 3, got %d", count)
	}
	if token != "312690787" {
		t.Errorf("Expected token to be 312690787, got %s", token)
	}
}

func TestGenerateStandardTokenRestricted(t *testing.T) {
	count, token, err := openpaygotoken.GenerateStandardToken(startingCode, &key, openpaygotoken.PAYGDisableValue, 1, openpaygotoken.SetTime, true)
	if err != nil {
		t.Error(err)
	}
	if count != 3 {
		t.Errorf("Expected count to be 3, got %d", count)
	}
	if token != "213331421312314" {
		t.Errorf("Expected token to be 213331421312314, got %s", token)
	}
}

func TestGenerateExtendedToken(t *testing.T) {
	count, token, err := openpaygotoken.GenerateExtendedToken(startingCode, &key, 1000, 1, false)
	if err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Errorf("Expected count to be 2, got %d", count)
	}
	if token != "315154457789" {
		t.Errorf("Expected token to be 315154457789, got %s", token)
	}
}

func TestDecodeStandardToken(t *testing.T) {
	decoder, err := openpaygotoken.NewDecoder()
	if err != nil {
		t.Error(err)
	}
	var usedCount []int = make([]int, 0)
	value, count, tokenType, err := decoder.GetActivationValueCountAndTypeFromToken(312690787, startingCode, &key, 0, false, &usedCount)
	if err != nil {
		t.Error(err)
	}
	if count != 3 {
		t.Errorf("Expected count to be 3, got %d", count)
	}
	if value != openpaygotoken.PAYGDisableValue {
		t.Errorf("Expected value to be %d, got %d", openpaygotoken.PAYGDisableValue, value)
	}
	if tokenType != openpaygotoken.SetTime {
		t.Errorf("Expected tokenType to be %d, got %d", openpaygotoken.SetTime, tokenType)
	}
}

func TestDecodeStandardTokenRestricted(t *testing.T) {
	decoder, err := openpaygotoken.NewDecoder()
	if err != nil {
		t.Error(err)
	}
	var usedCount []int = make([]int, 0)
	value, count, tokenType, err := decoder.GetActivationValueCountAndTypeFromToken(213331421312314, startingCode, &key, 0, true, &usedCount)
	if err != nil {
		t.Error(err)
	}
	if count != 3 {
		t.Errorf("Expected count to be 3, got %d", count)
	}
	if value != openpaygotoken.PAYGDisableValue {
		t.Errorf("Expected value to be %d, got %d", openpaygotoken.PAYGDisableValue, value)
	}
	if tokenType != openpaygotoken.SetTime {
		t.Errorf("Expected tokenType to be %d, got %d", openpaygotoken.SetTime, tokenType)
	}
}

func TestDecodeExtendedToken(t *testing.T) {
	decoder, err := openpaygotoken.NewDecoder()
	if err != nil {
		t.Error(err)
	}
	var usedCount []int = make([]int, 0)
	value, count, err := decoder.GetActivationValueCountAndTypeFromExtendedToken(315154457789, startingCode, &key, 1, false, &usedCount)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Errorf("Expected count to be 1, got %d", count)
	}
	if value != 1000 {
		t.Errorf("Expected value to be 1000, got %d", value)
	}

}

func TestUpdateUsedCount(t *testing.T) {
	decoder, err := openpaygotoken.NewDecoder()
	if err != nil {
		t.Error(err)
	}
	var usedCount []int = make([]int, 0)
	usedCount = decoder.UpdateUsedCounts(&usedCount, 1, 3, openpaygotoken.SetTime)
	if len(usedCount) != 17 {
		t.Errorf("Expected usedCount to be 17, got %d", len(usedCount))
	}
	usedCount = decoder.UpdateUsedCounts(&usedCount, 1, 5, openpaygotoken.SetTime)
	if len(usedCount) != 17 {
		t.Errorf("Expected usedCount to be 17, got %d", len(usedCount))
	}
	usedCount = decoder.UpdateUsedCounts(&usedCount, 1, 6, openpaygotoken.AddTime)
	if len(usedCount) != 17 {
		t.Errorf("Expected usedCount to be 17, got %d", len(usedCount))
	}
	usedCount = decoder.UpdateUsedCounts(&usedCount, 1, 100, openpaygotoken.AddTime)
	if len(usedCount) != 1 {
		t.Errorf("Expected usedCount to be 1, got %d", len(usedCount))
	}
	if slices.Contains(usedCount, 100) == false {
		t.Errorf("Expected usedCount to contain 100")
	}
	usedCount = decoder.UpdateUsedCounts(&usedCount, 1, 98, openpaygotoken.AddTime)
	if len(usedCount) != 2 {
		t.Errorf("Expected usedCount to be 2, got %d", len(usedCount))
	}
	if slices.Contains(usedCount, 98) == false || slices.Contains(usedCount, 100) == false {
		t.Errorf("Expected usedCount to contain 100 and 98")
	}
}
