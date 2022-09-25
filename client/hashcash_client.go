package main

import (
	"fmt"
	"time"

	"github.com/smirzoavliyoev/word_of_wisdom_test/pkg/utils"
)

// https://pkg.go.dev/github.com/umahmood/hashcash
// http://www.hashcash.org/dev/
const (
	maxIterations  int    = 1 << 20        // Max iterations to find a solution
	bytesToRead    int    = 8              // Bytes to read for random token
	bitsPerHexChar int    = 4              // Each hex character takes 4 bits
	zero           rune   = 48             // ASCII code for number zero
	timeFormat     string = "060102150405" // YYMMDDhhmmss
)

// Resource represents a hashcash resource
type Resource struct {
	// Data email, IP address, etc...
	Data string
	// ValidatorFunc user supplied function which validates Data
	ValidatorFunc func(string) bool
}

// Config for a hashcash instance
type Config struct {
	// Bits recommended default collision sizes are 20-bits
	Bits int
	// Expiry time before hashcash tokens are considered expired. Recommended
	// expiry time is 28 days
	Expired time.Time
	// Future hashcash in the future that should be rejected. Recommended
	// tolerance for clock skew is 48 hours
	Future time.Time
}

// DefaultConfig default hashcash configuration
var DefaultConfig = &Config{
	Bits:    40,
	Future:  time.Now().AddDate(0, 0, 2),
	Expired: time.Now().AddDate(0, 0, -30),
}

// Hashcash instance
type HashcashClient struct {
	// version hashcash format version, 1 (which supersedes version 0).
	version int
	// bits number of "partial pre-image" (zero) bits in the hashed code.
	bits int
	// created date The time that the message was sent.
	created time.Time
	// resource data string being transmitted, e.g., an IP address or email address.
	resource string
	// extension (optional; ignored in version 1).
	extension string
	// rand characters, encoded in base-64 format.
	rand string
	// counter (up to 2^20), encoded in base-64 format.
	counter int
}

// Compute a new hashcash header. If no solution can be found 'ErrSolutionFail'
// error is returned.
func (h *HashcashClient) Compute() (string, error) {
	// hex char: 0    0    0    0    0
	// binary  : 0000 0000 0000 0000 0000 = 4 bits per char = 20 bits total
	var (
		wantZeros = h.bits / bitsPerHexChar
		header    = h.createHeader()
		hash      = utils.Sha1Hash(header)
	)
	for !acceptableHeader(hash, zero, wantZeros) {
		h.counter++
		header = h.createHeader()
		hash = utils.Sha1Hash(header)
		if h.counter >= maxIterations {
			return "", ErrSolutionFail
		}
	}
	return header, nil
}

// New creates a new Hashcash instance
func New(res *Resource, config *Config) (*HashcashClient, error) {
	if res == nil {
		return nil, ErrResourceEmpty
	}
	if config == nil {
		config = DefaultConfig
	}

	rand, err := utils.RandomBytes(bytesToRead)
	if err != nil {
		return nil, err
	}

	return &HashcashClient{
		version:   1,
		bits:      config.Bits,
		created:   time.Now(),
		resource:  res.Data,
		extension: "",
		rand:      utils.Base64EncodeBytes(rand),
		counter:   1,
	}, nil
}

// acceptableHeader determines if the string 'hash' is prefixed with 'n',
// 'char' characters.
func acceptableHeader(hash string, char rune, n int) bool {
	for _, val := range hash[:n] {
		if val != char {
			return false
		}
	}
	return true
}

// createHeader creates a new hashcash header
func (h *HashcashClient) createHeader() string {
	return fmt.Sprintf("%d:%d:%s:%s:%s:%s:%s", h.version,
		h.bits,
		h.created.Format(timeFormat),
		h.resource,
		h.extension,
		h.rand,
		utils.Base64EncodeInt(h.counter))
}
