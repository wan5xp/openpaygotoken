package openpaygotoken

import (
	"fmt"
	"strconv"
)

// GenerateStandardToken generates a token with the given parameters.
// The token is generated from the starting code, the key, the value, the count and the mode.
// This function returns the count, the token and an error if there is one.
func GenerateStandardToken(startingCode int, key *[16]byte, value int, count int, mode TokenType, restrictedDigitSet bool) (int, string, error) {
	// We get the first 3 digits with encoded value
	startingCodeBase := getTokenBase(startingCode)
	tokenBase := encodeBase(startingCodeBase, value)
	currentToken, err := putBaseInToken(startingCode, tokenBase)
	if err != nil {
		return 0, "", err
	}
	currentCountOdd := count%2 == 1
	var newCount int
	if mode == SetTime {
		if currentCountOdd { // Odd numbers are for SetTime
			newCount = count + 2
		} else {
			newCount = count + 1
		}
	} else {
		if currentCountOdd { // Even numbers are for AddTime
			newCount = count + 1
		} else {
			newCount = count + 2
		}
	}
	for xn := 0; xn < newCount; xn++ {
		currentToken = generateNextToken(currentToken, key)
	}
	finalToken, err := putBaseInToken(currentToken, tokenBase)
	if err != nil {
		return 0, "", err
	}
	if restrictedDigitSet {
		finalToken = convertTo4DigitToken(finalToken)
		return newCount, fmt.Sprintf("%015d", finalToken), nil
	} else {
		return newCount, fmt.Sprintf("%09d", finalToken), nil
	}
}

// GenerateExtendedToken generates a token with the given parameters.
// The token is generated from the starting code, the key, the value, the count and the mode.
// This function returns the count, the token and an error if there is one.
func GenerateExtendedToken(startingCode int, key *[16]byte, value int, count int, restrictedDigitSet bool) (int, string, error) {
	startingCodeBase := getTokenBaseExtended(startingCode)
	tokenBase := encodeBaseExtended(startingCodeBase, value)
	currentToken, err := putBaseInTokenExtended(startingCode, tokenBase)
	if err != nil {
		return 0, "", err
	}
	newCount := count + 1
	for xn := 0; xn < newCount; xn++ {
		currentToken = generateNextTokenExtended(currentToken, key)
	}
	finalToken, err := putBaseInTokenExtended(currentToken, tokenBase)
	if err != nil {
		return 0, "", err
	}
	if restrictedDigitSet {
		finalToken = convertTo4DigitToken(finalToken)
		return newCount, fmt.Sprintf("%020d", finalToken), nil
	} else {
		return newCount, fmt.Sprintf("%012d", finalToken), nil
	}
}

// Encode base for token
func encodeBase(base int, number int) int {
	if number+base > maxBase {
		return number + base - tokenValueOffset
	} else {
		return number + base
	}
}

// Encode base for extended token
func encodeBaseExtended(base int, number int) int {
	if number+base > maxBaseExtended {
		return number + base - tokenValueOffsetExtended
	} else {
		return number + base
	}
}

// Convert token to restricted digit set
func convertTo4DigitToken(token int) int {
	encoded := 0
	for _, digit := range strconv.FormatInt(int64(token), 4) {
		encoded = encoded*10 + int(digit-'0') + 1
	}

	return encoded
}
