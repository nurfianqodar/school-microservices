package hasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/nurfianqodar/school-microservices/utils/errs"
	"golang.org/x/crypto/argon2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Hash config
type Config struct {
	Memory, Iterations, SaltLength, KeyLength uint32
	Parallelism                               uint8
}

var DefaultConfig = &Config{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func GenerateFromPassword(password string, c *Config) (string, error) {
	saltBytes, err := genRandomBytes(c.KeyLength)
	if err != nil {
		return "", err
	}
	hashBytes := argon2.IDKey([]byte(password), saltBytes, c.Iterations, c.Memory, c.Parallelism, c.KeyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(saltBytes)
	b64Hash := base64.RawStdEncoding.EncodeToString(hashBytes)
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, c.Memory, c.Iterations, c.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func CompareHashWithPassword(hash string, password string) error {
	c, saltBytes, hashBytes, err := decodeHash(hash)
	if err != nil {
		return err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), saltBytes, c.Iterations, c.Memory, c.Parallelism, c.KeyLength)

	if subtle.ConstantTimeCompare(hashBytes, otherHash) == 1 {
		return nil
	}

	log.Println("error: password is incompatible")
	return errs.ErrInvalidCredential
}

func genRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("error: failed creating random bytes. %s\n", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return b, nil
}

func decodeHash(encodedHashString string) (c *Config, salt, hash []byte, err error) {
	vals := strings.Split(encodedHashString, "$")
	if len(vals) != 6 {
		log.Println("error: invalid hash format")
		return nil, nil, nil, errs.ErrInvalidCredential
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		log.Println("error: incompatible argon2 version")
		return nil, nil, nil, errs.ErrInvalidCredential
	}

	c = new(Config)
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &c.Memory, &c.Iterations, &c.Parallelism)
	if err != nil {
		log.Println("error: unable to parse argon2 configuration")
		return nil, nil, nil, errs.ErrInvalidCredential
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		log.Printf("error: unable to get or decode salt %s\n", err.Error())
		return nil, nil, nil, errs.ErrInvalidCredential
	}
	c.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	c.KeyLength = uint32(len(hash))

	return c, salt, hash, nil
}
