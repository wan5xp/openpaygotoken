package openpaygotoken

import (
	"encoding/binary"

	"github.com/aead/siphash"
)

// TokenType is the type of token.
type TokenType int

const (
	maxBase                  int = 999
	tokenValueOffset         int = 1000
	maxBaseExtended          int = 999999
	tokenValueOffsetExtended int = 1000000

	// PAYGDisableValue is the value of the PAYG disable token.
	PAYGDisableValue int = 998
	// CounterSyncValue is the value of the counter sync token.
	CounterSyncValue int = 999

	// SetTime is the token type for setting the time.
	SetTime TokenType = 1
	// AddTime is the token type for adding time.
	AddTime TokenType = 2
	// MaxActivationValue is the maximum value of an activation token.
	MaxActivationValue int = 995
)

// GetTokenBase returns the base of the token.
func getTokenBase(code int) int {
	return code % tokenValueOffset
}

// PutBaseInToken returns the token with the given base.
func putBaseInToken(token int, tokenBase int) (int, error) {
	if tokenBase > maxBase {
		return 0, &ErrInvalidTokenBase{Value: tokenBase}
	}
	return token - getTokenBase(token) + tokenBase, nil
}

// GetTokenBaseExtended returns the base of the extended token.
func getTokenBaseExtended(code int) int {
	return code % tokenValueOffsetExtended
}

// PutBaseInTokenExtended returns the extended token with the given base.
func putBaseInTokenExtended(token int, tokenBase int) (int, error) {
	if tokenBase > maxBaseExtended {
		return 0, &ErrInvalidTokenBase{Value: tokenBase}
	}
	return token - getTokenBaseExtended(token) + tokenBase, nil
}

// GenerateNextToken generates a token with the given parameters.
func generateNextToken(lastCode int, key *[16]byte) int {
	conformedToken := make([]byte, 8)
	binary.BigEndian.PutUint32(conformedToken, uint32(lastCode))     // We convert the token to bytes
	binary.BigEndian.PutUint32(conformedToken[4:], uint32(lastCode)) // We duplicate it to fit the minimum length
	tokenHash := siphash.Sum64(conformedToken, key)                  // We hash it
	newToken := convertHashToToken(tokenHash)                        // We convert to token and return
	return newToken
}

// GenerateNextTokenExtended generates an extended token with the given parameters.
func generateNextTokenExtended(lastCode int, key *[16]byte) int {
	conformedToken := make([]byte, 8)
	binary.BigEndian.PutUint64(conformedToken, uint64(lastCode)) // We convert the token to bytes
	tokenHash := siphash.Sum64(conformedToken, key)              // We hash it
	newToken := convertHashToTokenExtended(tokenHash)            // We convert to token and return
	return newToken
}

// convertHashToToken converts hashed value to token.
func convertHashToToken(thisHash uint64) int {
	hashInt := make([]byte, 8)
	binary.BigEndian.PutUint64(hashInt, thisHash)  // We convert the hash to bytes
	hiHash := binary.BigEndian.Uint32(hashInt[:4]) // We split it in two 32bits INT
	loHash := binary.BigEndian.Uint32(hashInt[4:])
	resultHash := hiHash ^ loHash           // We XOR the two together to get a single 32bits INT
	token := convertoTo29_5Bits(resultHash) // We convert the 32bits value to an INT no greater than 9 digits
	return int(token)
}

// convertTo29_5Bits converts a 32bits value to an INT no greater than 9 digits.
func convertoTo29_5Bits(source uint32) uint32 {
	var mask uint32 = 0xFFFFFFFC
	temp := (source & mask) >> 2
	if temp > 999999999 {
		temp = temp - 73741825
	}
	return temp
}

// convertHashToTokenExtended converts hashed value to token.
func convertHashToTokenExtended(thisHash uint64) int {
	token := convertoTo40BitsExtended(thisHash) // We convert the 64bits value to an INT no greater than 12 digits
	return int(token)
}

// convertTo40BitsExtended converts a 64bits value to an INT no greater than 12 digits.
func convertoTo40BitsExtended(source uint64) uint64 {
	var mask uint64 = 0xFFFFFFFFFF000000
	temp := (source & mask) >> 24
	if temp > 999999999999 {
		temp = temp - 99511627777
	}
	return temp
}
