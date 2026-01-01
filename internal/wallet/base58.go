// ============================================================
// Base58 Encoding - Bitcoin/Solana alphabet
// ============================================================

package wallet

import (
	"math/big"
)

// Base58 alphabet used by Bitcoin and Solana
const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var (
	bigZero  = big.NewInt(0)
	bigRadix = big.NewInt(58)
)

// Base58Encode encodes a byte slice to a base58 string
func Base58Encode(input []byte) string {
	if len(input) == 0 {
		return ""
	}

	// Count leading zeros
	leadingZeros := 0
	for _, b := range input {
		if b == 0 {
			leadingZeros++
		} else {
			break
		}
	}

	// Convert to big integer
	num := new(big.Int).SetBytes(input)

	// Build base58 string in reverse
	var encoded []byte
	mod := new(big.Int)

	for num.Cmp(bigZero) > 0 {
		num.DivMod(num, bigRadix, mod)
		encoded = append(encoded, alphabet[mod.Int64()])
	}

	// Add leading '1's for each leading zero byte
	for i := 0; i < leadingZeros; i++ {
		encoded = append(encoded, alphabet[0])
	}

	// Reverse the result
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}

	return string(encoded)
}

// Base58Decode decodes a base58 string to bytes
func Base58Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil
	}

	// Build alphabet index map
	alphabetMap := make(map[rune]int64)
	for i, c := range alphabet {
		alphabetMap[c] = int64(i)
	}

	// Count leading '1's
	leadingOnes := 0
	for _, c := range input {
		if c == '1' {
			leadingOnes++
		} else {
			break
		}
	}

	// Convert from base58
	num := big.NewInt(0)
	for _, c := range input {
		idx, ok := alphabetMap[c]
		if !ok {
			return nil, &Base58Error{Char: c}
		}
		num.Mul(num, bigRadix)
		num.Add(num, big.NewInt(idx))
	}

	// Convert to bytes
	decoded := num.Bytes()

	// Add leading zeros
	result := make([]byte, leadingOnes+len(decoded))
	copy(result[leadingOnes:], decoded)

	return result, nil
}

// Base58Error is returned when an invalid character is encountered
type Base58Error struct {
	Char rune
}

func (e *Base58Error) Error() string {
	return "invalid base58 character: " + string(e.Char)
}

