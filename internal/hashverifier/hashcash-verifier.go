package hashverifier

import (
	"strings"
	"time"

	"github.com/smirzoavliyoev/word_of_wisdom_test/internal/hashverifier/repo"
	"github.com/smirzoavliyoev/word_of_wisdom_test/pkg/utils"
)

// http://www.hashcash.org/dev/
const (
	bytesToRead      int    = 8              // Bytes to read for random token
	bitsPerHexChar   int    = 4              // Each hex character takes 4 bits
	zero             rune   = 48             // ASCII code for number zero
	hashcashV1Length int    = 7              // Number of items in a V1 hashcash header
	timeFormat       string = "060102150405" // YYMMDDhhmmss
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
	// Storage underlying storage where hashcash tokens are stored and retrieved.
	Storage *repo.Repository
}

// DefaultConfig default hashcash configuration
var DefaultConfig = &Config{
	Bits:    20,
	Future:  time.Now().AddDate(0, 0, 2),
	Expired: time.Now().AddDate(0, 0, -30),
}

// Hashcash instance
type Hashcash struct {
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
	// validatorFunc user supplied function which validates resource
	validatorFunc func(string) bool
	// expired expiry time for headers
	expired time.Time
	// future tolerance for clock skew
	future time.Time
	// store the spent hashcash stamps
	storage *repo.Repository
}

// Verify that a hashcash header is valid. If the header is not in a valid
// format, ErrInvalidHeader error is returned.
func (h *Hashcash) Verify(header string) (bool, error) {
	vals := strings.Split(header, ":")
	if len(vals) != hashcashV1Length {
		return false, ErrInvalidHeader
	}
	// vals: [version bits date resource extension random counter]
	var (
		hash      = utils.Sha1Hash(header)
		wantZeros = h.bits / bitsPerHexChar
	)
	// test 1 - zero count
	if !acceptableHeader(hash, zero, wantZeros) {
		return false, ErrNoCollision
	}
	// test 2 - check token is not too far in the future or expired
	created, err := parseHashcashTime(vals[2])
	if err != nil {
		return false, err
	}
	if created.After(h.future) || created.Before(h.expired) {
		return false, ErrTimestamp
	}
	// test 3 - check resource is valid
	resource := vals[3]
	if !h.validatorFunc(resource) {
		return false, ErrResourceFail
	}
	// test 4 - check if hash is in spent storage
	if h.storage.Spent(hash) {
		return false, ErrSpent
	}
	h.storage.Add(hash)
	return true, nil
}

// New creates a new Hashcash instance
func New(res *Resource, config *Config) (*Hashcash, error) {
	if res == nil {
		return nil, ErrResourceEmpty
	}
	if config == nil {
		config = DefaultConfig
	}
	if config.Storage == nil {
		config.Storage = repo.NewRepo()
	}

	rand, err := utils.RandomBytes(bytesToRead)
	if err != nil {
		return nil, err
	}

	return &Hashcash{
		version:       1,
		bits:          config.Bits,
		created:       time.Now(),
		resource:      res.Data,
		validatorFunc: res.ValidatorFunc,
		extension:     "",
		rand:          utils.Base64EncodeBytes(rand),
		counter:       1,
		expired:       config.Expired,
		future:        config.Future,
		storage:       config.Storage,
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

// parseHashcashTime parses datetime in hashcash format
func parseHashcashTime(msgTime string) (date time.Time, err error) {
	// In a hashcash header the date parts year, month and day are mandatory but
	// hours, minutes and seconds are optional. So a valid date can be in format:
	//
	// - YYMMDD
	// - YYMMDDhhmm
	// - YYMMDDhhmmss
	//
	// Here we try find the format of the time, so it can be parsed.
	switch len(msgTime) {
	case 6:
		f := timeFormat[:6]
		date, err = time.Parse(f, msgTime)
	case 10:
		f := timeFormat[:10]
		date, err = time.Parse(f, msgTime)
	case 12:
		f := timeFormat[:12]
		date, err = time.Parse(f, msgTime)
	}
	return date, err
}
