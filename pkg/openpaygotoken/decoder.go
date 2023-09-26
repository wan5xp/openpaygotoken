package openpaygotoken

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

type TokenDecoder struct {
	maxTokenJump            int
	maxTokenJumpCounterSync int
	maxUnusedOlderToken     int
}

const (
	defaultMaxTokenJumpCounterSync int = 100
	defaultMaxTokenJump            int = 64
	defaultMaxUnusedOldToken       int = 8 * 2
)

// NewDecoder creates a new TokenDecoder with the given parameters.
// If no parameters are given, the default values are used.
// If only one parameter is given, it is used as maxTokenJump.
// If two parameters are given, the first is used as maxTokenJump and the second as maxTokenJumpCounterSync.
// If three or more parameters are given, the first is used as maxTokenJump, the second as maxTokenJumpCounterSync and the third as maxUnusedOlderToken.
func NewDecoder(args ...int) (*TokenDecoder, error) {
	if len(args) == 0 {
		return &TokenDecoder{maxTokenJump: defaultMaxTokenJump, maxTokenJumpCounterSync: defaultMaxTokenJumpCounterSync, maxUnusedOlderToken: defaultMaxUnusedOldToken}, nil
	} else if len(args) == 1 {
		return &TokenDecoder{maxTokenJump: args[0], maxTokenJumpCounterSync: defaultMaxTokenJumpCounterSync, maxUnusedOlderToken: defaultMaxUnusedOldToken}, nil
	} else if len(args) == 2 {
		return &TokenDecoder{maxTokenJump: args[0], maxTokenJumpCounterSync: args[1], maxUnusedOlderToken: defaultMaxUnusedOldToken}, nil
	} else {
		return &TokenDecoder{maxTokenJump: args[0], maxTokenJumpCounterSync: args[1], maxUnusedOlderToken: args[2]}, nil
	}
}

// GetActivationValueCountAndTypeFromToken returns the value, count and type of the token.
// If the token is not valid, an error is returned.
func (d *TokenDecoder) GetActivationValueCountAndTypeFromToken(token int, startingCode int, key *[16]byte, lastCount int, restrictedDigitSet bool, usedCounts *[]int) (int, int, TokenType, error) {
	if restrictedDigitSet {
		token = int(convertFrom4DigitToken(token))
	}
	validOlderToken := false
	tokenBase := getTokenBase(token)                            // We get the base of the token
	currentCode, err := putBaseInToken(startingCode, tokenBase) // We put the base in the starting code
	if err != nil {
		return 0, 0, 0, err
	}
	startingCodeBase := getTokenBase(startingCode)   // We get the base of the starting code
	value := decodeBase(startingCodeBase, tokenBase) // If there is a match we get the value from the token
	// We try all combination up until last_count + TOKEN_JUMP, or to the larger jump if syncing counter
	// We could start directly the loop at the last count if we kept the token value for the last count
	var maxCountTry int
	if value == CounterSyncValue {
		maxCountTry = lastCount + d.maxTokenJumpCounterSync + 1
	} else {
		maxCountTry = lastCount + d.maxTokenJump + 1
	}
	for count := 0; count < maxCountTry; count++ {
		maskedToken, err := putBaseInToken(currentCode, tokenBase)
		if err != nil {
			return 0, 0, 0, err
		}
		if maskedToken == token {
			var thisType TokenType
			if count%2 == 1 {
				thisType = SetTime
			} else {
				thisType = AddTime
			}
			if d.countIsValid(count, lastCount, value, thisType, usedCounts) {
				return value, count, thisType, nil
			} else {
				validOlderToken = true
			}
		}
		currentCode = generateNextToken(currentCode, key) // If not we go to the next token
	}
	if validOlderToken {
		return -2, 0, 0, nil
	}
	return 0, 0, 0, &ErrInvalidToken{}
}

// Check if count is valid
func (d *TokenDecoder) countIsValid(count int, lastCount int, value int, tokenType TokenType, usedCounts *[]int) bool {
	if value == CounterSyncValue {
		if count > lastCount-30 {
			return true
		}
	} else if count > lastCount {
		return true
	} else if d.maxUnusedOlderToken > 0 {
		if count > lastCount-d.maxUnusedOlderToken {
			if !slices.Contains(*usedCounts, count) && tokenType == AddTime {
				return true
			}
		}
	}
	return false
}

// UpdateUsedCounts returns the list of used counts.
func (d *TokenDecoder) UpdateUsedCounts(pastUsedCounts *[]int, value int, newCount int, tokenType TokenType) []int {
	highestCount := 0
	if pastUsedCounts != nil && len(*pastUsedCounts) > 0 {
		highestCount = slices.Max(*pastUsedCounts)
	}
	if newCount > highestCount {
		highestCount = newCount
	}
	bottomRange := highestCount - d.maxUnusedOlderToken
	var usedCounts []int
	if tokenType != AddTime || value == CounterSyncValue || value == PAYGDisableValue {
		// If it isnot an Add TIme token, we mark al the past tokens as used in the range
		for count := bottomRange; count <= highestCount; count++ {
			usedCounts = append(usedCounts, count)
		}
	} else {
		// If it is an Add Time token, we just mark the tokens actually used in the range
		for count := bottomRange; count <= highestCount; count++ {
			if count == newCount || slices.Contains(*pastUsedCounts, count) {
				usedCounts = append(usedCounts, count)
			}
		}
	}
	return usedCounts
}

// GetActivationValueCountAndTypeFromExtendedToken returns the value, count and type of the token.
// If the token is not valid, an error is returned.
func (d *TokenDecoder) GetActivationValueCountAndTypeFromExtendedToken(token int, startingCode int, key *[16]byte, lastCount int, restrictedDigitSet bool, usedCounts *[]int) (int, int, error) {
	if restrictedDigitSet {
		token = int(convertFrom4DigitToken(token))
	}
	tokenBase := getTokenBaseExtended(token)                            // We get the base of the token
	currentCode, err := putBaseInTokenExtended(startingCode, tokenBase) // We put the base in the starting code
	if err != nil {
		return 0, 0, err
	}
	startingCodeBase := getTokenBaseExtended(startingCode)   // We get the base of the starting code
	value := decodeBaseExtended(startingCodeBase, tokenBase) // If there is a match we get the value from the token
	for count := 0; count < 30; count++ {
		maskedToken, err := putBaseInTokenExtended(currentCode, tokenBase)
		if err != nil {
			return 0, 0, err
		}
		if maskedToken == token && count > lastCount {
			cleanCount := count - 1
			return value, cleanCount, nil
		}
		currentCode = generateNextTokenExtended(currentCode, key) // If not we go to the next token
	}
	return 0, 0, fmt.Errorf("token not found")
}

// Get decode base
func decodeBase(startingCodeBase int, tokenBase int) int {
	if tokenBase < startingCodeBase {
		return tokenBase + tokenValueOffset - startingCodeBase
	} else {
		return tokenBase - startingCodeBase
	}
}

// Get decode base for extended token
func decodeBaseExtended(startingCodeBase int, tokenBase int) int {
	if tokenBase < startingCodeBase {
		return tokenBase + tokenValueOffsetExtended - startingCodeBase
	} else {
		return tokenBase - startingCodeBase
	}
}

// Convert token for restricted digit set
func convertFrom4DigitToken(token int) int64 {
	var decoded int64 = 0
	for _, digit := range fmt.Sprintf("%d", token) {
		decoded = decoded*10 + (int64(digit-'0') - 1)
	}
	decoded, err := strconv.ParseInt(fmt.Sprintf("%d", decoded), 4, 64)
	if err != nil {
		return 0
	}
	return decoded
}
