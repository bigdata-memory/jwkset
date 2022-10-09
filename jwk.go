package jwkset

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const (
	// KeyTypeEC is the key type for ECDSA.
	KeyTypeEC KeyType = "EC"
	// KeyTypeOKP is the key type for EdDSA.
	KeyTypeOKP KeyType = "OKP"
	// KeyTypeRSA is the key type for RSA.
	KeyTypeRSA KeyType = "RSA"
	// KeyTypeOct is the key type for octet sequences, such as HMAC.
	KeyTypeOct KeyType = "oct"

	// CurveEd25519 is the curve for EdDSA.
	CurveEd25519 JWKCRV = "Ed25519"
	// CurveP256 is the curve for ECDSA.
	CurveP256 JWKCRV = "P-256"
	// CurveP384 is the curve for ECDSA.
	CurveP384 JWKCRV = "P-384"
	// CurveP521 is the curve for ECDSA.
	CurveP521 JWKCRV = "P-521"
)

var (
	// ErrKeyUnmarshalParameter indicates that a JWK's attributes are invalid and cannot be unmarshalled.
	ErrKeyUnmarshalParameter = errors.New("unable to unmarshal JWK due to invalid attributes")
	// ErrUnsupportedKeyType indicates a key type is not supported.
	ErrUnsupportedKeyType = errors.New("unsupported key type")
)

// JWKCRV is a set of "JSON Web Key Elliptic JWKCRV" types from https://www.iana.org/assignments/jose/jose.xhtml as
// mentioned in https://www.rfc-editor.org/rfc/rfc7518.html#section-6.2.1.1.
type JWKCRV string

func (crv JWKCRV) String() string {
	return string(crv)
}

// KeyType is a set of "JSON Web Key Types" from https://www.iana.org/assignments/jose/jose.xhtml as mentioned in
// https://www.rfc-editor.org/rfc/rfc7517#section-4.1
type KeyType string

func (kty KeyType) String() string {
	return string(kty)
}

// KeyWithMeta is holds a Key and its metadata.
type KeyWithMeta struct {
	Key   interface{}
	KeyID string
}

// NewKey creates a new KeyWithMeta.
func NewKey(key interface{}, keyID string) KeyWithMeta {
	return KeyWithMeta{
		Key:   key,
		KeyID: keyID,
	}
}

// OtherPrimes is for RSA private keys that have more than 2 primes.
// https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7
type OtherPrimes struct {
	CRTFactorExponent    string `json:"d,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7.2
	CRTFactorCoefficient string `json:"t,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7.3
	PrimeFactor          string `json:"r,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7.1
}

// JWKMarshal is used to marshal or unmarshal a JSON Web Key.
// https://www.rfc-editor.org/rfc/rfc7517
// https://www.rfc-editor.org/rfc/rfc7518
// https://www.rfc-editor.org/rfc/rfc8037
type JWKMarshal struct {
	CRV string        `json:"crv,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.2.1.1 and https://www.rfc-editor.org/rfc/rfc8037.html#section-2
	D   string        `json:"d,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.1 and https://www.rfc-editor.org/rfc/rfc7518#section-6.2.2.1 and https://www.rfc-editor.org/rfc/rfc8037.html#section-2
	DP  string        `json:"dp,omitempty"`  // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.4
	DQ  string        `json:"dq,omitempty"`  // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.5
	E   string        `json:"e,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.1.2
	K   string        `json:"k,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.4.1
	KID string        `json:"kid,omitempty"` // https://www.rfc-editor.org/rfc/rfc7517#section-4.5
	KTY string        `json:"kty,omitempty"` // https://www.rfc-editor.org/rfc/rfc7517#section-4.1
	N   string        `json:"n,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.1.1
	OTH []OtherPrimes `json:"oth,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7
	P   string        `json:"p,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.2
	Q   string        `json:"q,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.3
	QI  string        `json:"qi,omitempty"`  // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.6
	X   string        `json:"x,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.2.1.2 and https://www.rfc-editor.org/rfc/rfc8037.html#section-2
	Y   string        `json:"y,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.2.1.3
	// TODO Use ALG field.
	// ALG string        `json:"alg,omitempty"` // https://www.rfc-editor.org/rfc/rfc7517#section-4.4 and https://www.rfc-editor.org/rfc/rfc7518#section-4.1
	// TODO Use KEYOPS field.
	// KEYOPTS []string `json:"key_ops,omitempty"` // https://www.rfc-editor.org/rfc/rfc7517#section-4.3
	// TODO Use USE field.
	// USE string        `json:"use,omitempty"` // https://www.rfc-editor.org/rfc/rfc7517#section-4.2
	// TODO X.509 related fields.
}

// JWKSMarshal is used to marshal or unmarshal a JSON Web Key Set.
type JWKSMarshal struct {
	Keys []JWKMarshal `json:"keys"`
}

// JWKSet is a set of JSON Web Keys.
type JWKSet struct {
	Store Storage
}

// NewMemory creates a new in-memory JWKSet.
func NewMemory() JWKSet {
	return JWKSet{
		Store: NewMemoryStorage(),
	}
}

// JSON creates the JSON representation of the JWKSet.
func (j JWKSet) JSON(ctx context.Context) (json.RawMessage, error) {
	jwks := JWKSMarshal{}
	options := KeyMarshalOptions{
		AsymmetricPrivate: false,
	}

	keys, err := j.Store.SnapshotKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot of all keys from storage: %w", err)
	}

	for _, meta := range keys {
		jwk, err := KeyMarshal(meta, options)
		if err != nil {
			if errors.Is(err, ErrUnsupportedKeyType) {
				// Ignore the key.
				continue
			}
			return nil, fmt.Errorf("failed to marshal key: %w", err)
		}
		jwks.Keys = append(jwks.Keys, jwk)
	}

	return json.Marshal(jwks)
}

// KeyMarshalOptions are used to specify options for marshalling a JSON Web Key.
type KeyMarshalOptions struct {
	AsymmetricPrivate bool
	Symmetric         bool
}

// KeyMarshal transforms a KeyWithMeta into a JWKMarshal, which is used to marshal/unmarshal a JSON Web Key.
func KeyMarshal(meta KeyWithMeta, options KeyMarshalOptions) (JWKMarshal, error) {
	var jwk JWKMarshal
	switch key := meta.Key.(type) {
	case *ecdsa.PrivateKey:
		pub := key.PublicKey
		jwk.CRV = pub.Curve.Params().Name
		jwk.X = bigIntToBase64RawURL(pub.X)
		jwk.Y = bigIntToBase64RawURL(pub.Y)
		jwk.KTY = KeyTypeEC.String()
		if options.AsymmetricPrivate {
			jwk.D = bigIntToBase64RawURL(key.D)
		}
	case ecdsa.PublicKey: // TODO Make this a pointer. Maybe support value with reassignment and fallthrough.
		jwk.CRV = key.Curve.Params().Name
		jwk.X = bigIntToBase64RawURL(key.X)
		jwk.Y = bigIntToBase64RawURL(key.Y)
		jwk.KTY = KeyTypeEC.String()
	case ed25519.PrivateKey:
		pub := key.Public().(ed25519.PublicKey)
		jwk.CRV = CurveEd25519.String()
		jwk.X = base64.RawURLEncoding.EncodeToString(pub)
		jwk.KTY = KeyTypeOKP.String()
		if options.AsymmetricPrivate {
			jwk.D = base64.RawURLEncoding.EncodeToString(key)
		}
	case ed25519.PublicKey:
		jwk.CRV = CurveEd25519.String()
		jwk.X = base64.RawURLEncoding.EncodeToString(key)
		jwk.KTY = KeyTypeOKP.String()
	case *rsa.PrivateKey:
		pub := key.PublicKey
		jwk.E = bigIntToBase64RawURL(big.NewInt(int64(pub.E)))
		jwk.N = bigIntToBase64RawURL(pub.N)
		jwk.KTY = KeyTypeRSA.String()
		if options.AsymmetricPrivate {
			jwk.D = bigIntToBase64RawURL(key.D)
			jwk.P = bigIntToBase64RawURL(key.Primes[0])
			jwk.Q = bigIntToBase64RawURL(key.Primes[1])
			jwk.DP = bigIntToBase64RawURL(key.Precomputed.Dp)
			jwk.DQ = bigIntToBase64RawURL(key.Precomputed.Dq)
			jwk.QI = bigIntToBase64RawURL(key.Precomputed.Qinv)
			for i := 2; i < len(key.Primes); i++ {
				jwk.OTH = append(jwk.OTH, OtherPrimes{
					CRTFactorExponent:    bigIntToBase64RawURL(key.Precomputed.CRTValues[i].Exp),
					CRTFactorCoefficient: bigIntToBase64RawURL(key.Precomputed.CRTValues[i].Coeff),
					PrimeFactor:          bigIntToBase64RawURL(key.Precomputed.CRTValues[i].R),
				})
			}
		}
	case rsa.PublicKey: // TODO Make this a pointer. Maybe support value with reassignment and fallthrough.
		jwk.E = bigIntToBase64RawURL(big.NewInt(int64(key.E)))
		jwk.N = bigIntToBase64RawURL(key.N)
		jwk.KTY = KeyTypeRSA.String()
	case []byte:
		if options.Symmetric {
			jwk.KTY = KeyTypeOct.String()
			jwk.K = base64.RawURLEncoding.EncodeToString(key)
		} else {
			return JWKMarshal{}, fmt.Errorf("%w: incorrect options to marshal symmetric key (oct)", ErrUnsupportedKeyType)
		}
	default:
		return JWKMarshal{}, fmt.Errorf("%w: %T", ErrUnsupportedKeyType, key)
	}
	jwk.KID = meta.KeyID
	return jwk, nil
}

type KeyUnmarshalOptions struct {
	AsymmetricPrivate bool
	Symmetric         bool
}

func KeyUnmarshal(jwk JWKMarshal, options KeyUnmarshalOptions) (KeyWithMeta, error) {
	meta := KeyWithMeta{}
	switch KeyType(jwk.KTY) {
	case KeyTypeEC:
		if jwk.X == "" || jwk.Y == "" || jwk.CRV == "" {
			return KeyWithMeta{}, fmt.Errorf("%w: %s requires parameters x, y, and crv", ErrKeyUnmarshalParameter, KeyTypeEC)
		}
		x, err := base64urlTrailingPadding(jwk.X)
		if err != nil {
			return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "x": %w`, KeyTypeEC, err)
		}
		y, err := base64urlTrailingPadding(jwk.Y)
		if err != nil {
			return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "y": %w`, KeyTypeEC, err)
		}
		publicKey := ecdsa.PublicKey{
			X: big.NewInt(0).SetBytes(x),
			Y: big.NewInt(0).SetBytes(y),
		}
		switch JWKCRV(jwk.CRV) {
		case CurveP256:
			publicKey.Curve = elliptic.P256()
		case CurveP384:
			publicKey.Curve = elliptic.P384()
		case CurveP521:
			publicKey.Curve = elliptic.P521()
		default:
			return KeyWithMeta{}, fmt.Errorf("%w: unsupported curve type %q", ErrKeyUnmarshalParameter, jwk.CRV)
		}
		if options.AsymmetricPrivate {
			if jwk.D == "" {
				return KeyWithMeta{}, fmt.Errorf(`%w: %s requires parameter "d"`, ErrKeyUnmarshalParameter, KeyTypeEC)
			}
			d, err := base64urlTrailingPadding(jwk.D)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "d": %w`, KeyTypeEC, err)
			}
			privateKey := &ecdsa.PrivateKey{
				PublicKey: publicKey,
				D:         big.NewInt(0).SetBytes(d),
			}
			meta.Key = privateKey
		} else {
			meta.Key = &publicKey
		}
	case KeyTypeOKP:
		if JWKCRV(jwk.CRV) != CurveEd25519 {
			return KeyWithMeta{}, fmt.Errorf("%w: %s key type should have %q curve", ErrUnsupportedKeyType, KeyTypeOKP, CurveEd25519)
		}
		if options.AsymmetricPrivate {
			if jwk.D == "" {
				return KeyWithMeta{}, fmt.Errorf(`%w: %s requires parameter "d"`, ErrKeyUnmarshalParameter, KeyTypeOKP)
			}
			key, err := base64urlTrailingPadding(jwk.D)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "d": %w`, KeyTypeOKP, err)
			}
			if len(key) != ed25519.PrivateKeySize {
				return KeyWithMeta{}, fmt.Errorf("%w: %s key should be %d bytes", ErrUnsupportedKeyType, KeyTypeOKP, ed25519.PrivateKeySize)
			}
			meta.Key = ed25519.PrivateKey(key)
		} else if !options.AsymmetricPrivate {
			if jwk.X == "" {
				return KeyWithMeta{}, fmt.Errorf(`%w: %s requires parameter "x"`, ErrKeyUnmarshalParameter, KeyTypeOKP)
			}
			key, err := base64urlTrailingPadding(jwk.X)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "x": %w`, KeyTypeOKP, err)
			}
			if len(key) != ed25519.PublicKeySize {
				return KeyWithMeta{}, fmt.Errorf("%w: %s key should be %d bytes", ErrUnsupportedKeyType, KeyTypeOKP, ed25519.PublicKeySize)
			}
			meta.Key = ed25519.PublicKey(key)
		}
	case KeyTypeRSA:
		// TODO
	case KeyTypeOct:
		if options.Symmetric {
			if jwk.K == "" {
				return KeyWithMeta{}, fmt.Errorf(`%w: %s requires parameter "k"`, ErrKeyUnmarshalParameter, KeyTypeOct)
			}
			key, err := base64urlTrailingPadding(jwk.K)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "k": %w`, KeyTypeOct, err)
			}
			meta.Key = key
		} else {
			return KeyWithMeta{}, fmt.Errorf("%w: incorrect options to unmarshal symmetric key (%s)", ErrUnsupportedKeyType, KeyTypeOct)
		}
	default:
		return KeyWithMeta{}, fmt.Errorf("%w: %s", ErrUnsupportedKeyType, jwk.KTY)
	}
	meta.KeyID = jwk.KID
	return meta, nil
}

// base64urlTrailingPadding removes trailing padding before decoding a string from base64url. Some non-RFC compliant
// JWKS contain padding at the end values for base64url encoded public keys.
//
// Trailing padding is required to be removed from base64url encoded keys.
// RFC 7517 defines base64url the same as RFC 7515 Section 2:
// https://datatracker.ietf.org/doc/html/rfc7517#section-1.1
// https://datatracker.ietf.org/doc/html/rfc7515#section-2
func base64urlTrailingPadding(s string) ([]byte, error) {
	s = strings.TrimRight(s, "=")
	return base64.RawURLEncoding.DecodeString(s)
}

func bigIntToBase64RawURL(i *big.Int) string {
	return base64.RawURLEncoding.EncodeToString(i.Bytes())
}
