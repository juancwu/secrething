package utils

import (
	"crypto/rand"
	"errors"
	"io"
)

// RandomBytes generates a cryptographically secure random byte array of given size.
func RandomBytes(size uint32) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// digitsTable is just a reference table which the RandomDigits
// function selects the digits to use in the generated sequence.
// The digits are characters so that they are easily transformed into string.
var digitsTable = [10]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

var ErrWrongSequenceLengthRead error = errors.New("Wrong sequence length read from rand.Reader")

// RandomDigits generates a random digit (0-9) sequence of the given size.
// The digit sequence contains the ascii value of digits 0-9 and not the
// actual binary value representation of each digit.
func RandomDigits(size uint32) ([]byte, error) {
	sequence := make([]byte, size)
	n, err := io.ReadFull(rand.Reader, sequence)
	if err != nil {
		return nil, err
	}
	if n != int(size) {
		return nil, ErrWrongSequenceLengthRead
	}
	for i := range sequence {
		sequence[i] = digitsTable[int(sequence[i])%10]
	}
	return sequence, nil
}
