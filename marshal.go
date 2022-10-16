package jwkset

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

const (
	// ALGEdDSA is the EdDSA algorithm.
	ALGEdDSA ALG = "EdDSA"

	// KeyTypeEC is the key type for ECDSA.
	KeyTypeEC KTY = "EC"
	// KeyTypeOKP is the key type for EdDSA.
	KeyTypeOKP KTY = "OKP"
	// KeyTypeRSA is the key type for RSA.
	KeyTypeRSA KTY = "RSA"
	// KeyTypeOct is the key type for octet sequences, such as HMAC.
	KeyTypeOct KTY = "oct"

	// CurveEd25519 is a curve for EdDSA.
	CurveEd25519 CRV = "Ed25519"
	// CurveP256 is a curve for ECDSA.
	CurveP256 CRV = "P-256"
	// CurveP384 is a curve for ECDSA.
	CurveP384 CRV = "P-384"
	// CurveP521 is a curve for ECDSA.
	CurveP521 CRV = "P-521"
)

var (
	// ErrKeyUnmarshalParameter indicates that a JWK's attributes are invalid and cannot be unmarshaled.
	ErrKeyUnmarshalParameter = errors.New("unable to unmarshal JWK due to invalid attributes")
	// ErrUnsupportedKeyType indicates a key type is not supported.
	ErrUnsupportedKeyType = errors.New("unsupported key type")
)

// ALG is a set of "JSON Web Signature and Encryption Algorithms" types from
// https://www.iana.org/assignments/jose/jose.xhtml(JWA) as defined in
// https://www.rfc-editor.org/rfc/rfc7518#section-7.1
type ALG string

// CRV is a set of "JSON Web Key Elliptic Curve" types from https://www.iana.org/assignments/jose/jose.xhtml as
// mentioned in https://www.rfc-editor.org/rfc/rfc7518.html#section-6.2.1.1.
type CRV string

func (crv CRV) String() string {
	return string(crv)
}

// KTY is a set of "JSON Web Key Types" from https://www.iana.org/assignments/jose/jose.xhtml as mentioned in
// https://www.rfc-editor.org/rfc/rfc7517#section-4.1
type KTY string

func (kty KTY) String() string {
	return string(kty)
}

// OtherPrimes is for RSA private keys that have more than 2 primes.
// https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7
type OtherPrimes struct {
	D string `json:"d,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7.2
	R string `json:"r,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7.1
	T string `json:"t,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7.3
}

// JWKMarshal is used to marshal or unmarshal a JSON Web Key.
// https://www.rfc-editor.org/rfc/rfc7517
// https://www.rfc-editor.org/rfc/rfc7518
// https://www.rfc-editor.org/rfc/rfc8037
type JWKMarshal struct {
	// TODO Use ALG field.
	ALG ALG    `json:"alg,omitempty"` // https://www.rfc-editor.org/rfc/rfc7517#section-4.4 and https://www.rfc-editor.org/rfc/rfc7518#section-4.1
	CRV CRV    `json:"crv,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.2.1.1 and https://www.rfc-editor.org/rfc/rfc8037.html#section-2
	D   string `json:"d,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.1 and https://www.rfc-editor.org/rfc/rfc7518#section-6.2.2.1 and https://www.rfc-editor.org/rfc/rfc8037.html#section-2
	DP  string `json:"dp,omitempty"`  // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.4
	DQ  string `json:"dq,omitempty"`  // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.5
	E   string `json:"e,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.1.2
	K   string `json:"k,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.4.1
	// TODO Use KEYOPS field.
	// KEYOPTS []string `json:"key_ops,omitempty"` // https://www.rfc-editor.org/rfc/rfc7517#section-4.3
	KID string        `json:"kid,omitempty"` // https://www.rfc-editor.org/rfc/rfc7517#section-4.5
	KTY KTY           `json:"kty,omitempty"` // https://www.rfc-editor.org/rfc/rfc7517#section-4.1
	N   string        `json:"n,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.1.1
	OTH []OtherPrimes `json:"oth,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.7
	P   string        `json:"p,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.2
	Q   string        `json:"q,omitempty"`   // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.3
	QI  string        `json:"qi,omitempty"`  // https://www.rfc-editor.org/rfc/rfc7518#section-6.3.2.6
	// TODO Use USE field.
	// USE USE        `json:"use,omitempty"` // https://www.rfc-editor.org/rfc/rfc7517#section-4.2
	X string `json:"x,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.2.1.2 and https://www.rfc-editor.org/rfc/rfc8037.html#section-2
	// TODO X.509 related fields.
	Y string `json:"y,omitempty"` // https://www.rfc-editor.org/rfc/rfc7518#section-6.2.1.3
}

// JWKSMarshal is used to marshal or unmarshal a JSON Web Key Set.
type JWKSMarshal struct {
	Keys []JWKMarshal `json:"keys"`
}

// KeyMarshalOptions are used to specify options for marshaling a JSON Web Key.
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
		jwk.CRV = CRV(pub.Curve.Params().Name)
		jwk.X = bigIntToBase64RawURL(pub.X)
		jwk.Y = bigIntToBase64RawURL(pub.Y)
		jwk.KTY = KeyTypeEC
		if options.AsymmetricPrivate {
			jwk.D = bigIntToBase64RawURL(key.D)
		}
	case *ecdsa.PublicKey:
		jwk.CRV = CRV(key.Curve.Params().Name)
		jwk.X = bigIntToBase64RawURL(key.X)
		jwk.Y = bigIntToBase64RawURL(key.Y)
		jwk.KTY = KeyTypeEC
	case ed25519.PrivateKey:
		pub := key.Public().(ed25519.PublicKey)
		jwk.ALG = ALGEdDSA
		jwk.CRV = CurveEd25519
		jwk.X = base64.RawURLEncoding.EncodeToString(pub)
		jwk.KTY = KeyTypeOKP
		if options.AsymmetricPrivate {
			jwk.D = base64.RawURLEncoding.EncodeToString(key[:32])
		}
	case ed25519.PublicKey:
		jwk.ALG = ALGEdDSA
		jwk.CRV = CurveEd25519
		jwk.X = base64.RawURLEncoding.EncodeToString(key)
		jwk.KTY = KeyTypeOKP
	case *rsa.PrivateKey:
		pub := key.PublicKey
		jwk.E = bigIntToBase64RawURL(big.NewInt(int64(pub.E)))
		jwk.N = bigIntToBase64RawURL(pub.N)
		jwk.KTY = KeyTypeRSA
		if options.AsymmetricPrivate {
			jwk.D = bigIntToBase64RawURL(key.D)
			jwk.P = bigIntToBase64RawURL(key.Primes[0])
			jwk.Q = bigIntToBase64RawURL(key.Primes[1])
			jwk.DP = bigIntToBase64RawURL(key.Precomputed.Dp)
			jwk.DQ = bigIntToBase64RawURL(key.Precomputed.Dq)
			jwk.QI = bigIntToBase64RawURL(key.Precomputed.Qinv)
			if len(key.Precomputed.CRTValues) > 0 {
				jwk.OTH = make([]OtherPrimes, len(key.Precomputed.CRTValues))
				for i := 0; i < len(key.Precomputed.CRTValues); i++ {
					jwk.OTH[i] = OtherPrimes{
						D: bigIntToBase64RawURL(key.Precomputed.CRTValues[i].Exp),
						T: bigIntToBase64RawURL(key.Precomputed.CRTValues[i].Coeff),
						R: bigIntToBase64RawURL(key.Precomputed.CRTValues[i].R),
					}
				}
			}
		}
	case *rsa.PublicKey:
		jwk.E = bigIntToBase64RawURL(big.NewInt(int64(key.E)))
		jwk.N = bigIntToBase64RawURL(key.N)
		jwk.KTY = KeyTypeRSA
	case []byte:
		if options.Symmetric {
			jwk.KTY = KeyTypeOct
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

// KeyUnmarshalOptions are used to specify options for unmarshaling a JSON Web Key.
type KeyUnmarshalOptions struct {
	AsymmetricPrivate bool
	Symmetric         bool
}

// KeyUnmarshal transforms a JWKMarshal into a KeyWithMeta, which contains the correct Go type for the cryptographic
// key.
func KeyUnmarshal(jwk JWKMarshal, options KeyUnmarshalOptions) (KeyWithMeta, error) {
	meta := KeyWithMeta{}
	switch jwk.KTY {
	case KeyTypeEC:
		if jwk.CRV == "" || jwk.X == "" || jwk.Y == "" {
			return KeyWithMeta{}, fmt.Errorf(`%w: %s requires parameters "crv", "x", and "y"`, ErrKeyUnmarshalParameter, KeyTypeEC)
		}
		x, err := base64urlTrailingPadding(jwk.X)
		if err != nil {
			return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "x": %w`, KeyTypeEC, err)
		}
		y, err := base64urlTrailingPadding(jwk.Y)
		if err != nil {
			return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "y": %w`, KeyTypeEC, err)
		}
		publicKey := &ecdsa.PublicKey{
			X: new(big.Int).SetBytes(x),
			Y: new(big.Int).SetBytes(y),
		}
		switch jwk.CRV {
		case CurveP256:
			publicKey.Curve = elliptic.P256()
		case CurveP384:
			publicKey.Curve = elliptic.P384()
		case CurveP521:
			publicKey.Curve = elliptic.P521()
		default:
			return KeyWithMeta{}, fmt.Errorf("%w: unsupported curve type %q", ErrKeyUnmarshalParameter, jwk.CRV)
		}
		if options.AsymmetricPrivate && jwk.D != "" {
			d, err := base64urlTrailingPadding(jwk.D)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "d": %w`, KeyTypeEC, err)
			}
			privateKey := &ecdsa.PrivateKey{
				PublicKey: *publicKey,
				D:         new(big.Int).SetBytes(d),
			}
			meta.Key = privateKey
		} else {
			meta.Key = publicKey
		}
	case KeyTypeOKP:
		if jwk.CRV != CurveEd25519 {
			return KeyWithMeta{}, fmt.Errorf("%w: %s key type should have %q curve", ErrKeyUnmarshalParameter, KeyTypeOKP, CurveEd25519)
		}
		if jwk.X == "" {
			return KeyWithMeta{}, fmt.Errorf(`%w: %s requires parameter "x"`, ErrKeyUnmarshalParameter, KeyTypeOKP)
		}
		public, err := base64urlTrailingPadding(jwk.X)
		if err != nil {
			return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "x": %w`, KeyTypeOKP, err)
		}
		if len(public) != ed25519.PublicKeySize {
			return KeyWithMeta{}, fmt.Errorf("%w: %s key should be %d bytes", ErrKeyUnmarshalParameter, KeyTypeOKP, ed25519.PublicKeySize)
		}
		if options.AsymmetricPrivate && jwk.D != "" {
			private, err := base64urlTrailingPadding(jwk.D)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "d": %w`, KeyTypeOKP, err)
			}
			private = append(private, public...)
			if len(private) != ed25519.PrivateKeySize {
				return KeyWithMeta{}, fmt.Errorf("%w: %s key should be %d bytes", ErrKeyUnmarshalParameter, KeyTypeOKP, ed25519.PrivateKeySize)
			}
			meta.Key = ed25519.PrivateKey(private)
		} else {
			meta.Key = ed25519.PublicKey(public)
		}
	case KeyTypeRSA:
		if jwk.N == "" || jwk.E == "" {
			return KeyWithMeta{}, fmt.Errorf(`%w: %s requires parameters "n" and "e"`, ErrKeyUnmarshalParameter, KeyTypeRSA)
		}
		n, err := base64urlTrailingPadding(jwk.N)
		if err != nil {
			return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "n": %w`, KeyTypeRSA, err)
		}
		e, err := base64urlTrailingPadding(jwk.E)
		if err != nil {
			return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "e": %w`, KeyTypeRSA, err)
		}
		publicKey := rsa.PublicKey{
			N: new(big.Int).SetBytes(n),
			E: int(new(big.Int).SetBytes(e).Uint64()),
		}
		if options.AsymmetricPrivate && jwk.D != "" && jwk.P != "" && jwk.Q != "" && jwk.DP != "" && jwk.DQ != "" && jwk.QI != "" {
			d, err := base64urlTrailingPadding(jwk.D)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "d": %w`, KeyTypeRSA, err)
			}
			p, err := base64urlTrailingPadding(jwk.P)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "p": %w`, KeyTypeRSA, err)
			}
			q, err := base64urlTrailingPadding(jwk.Q)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "q": %w`, KeyTypeRSA, err)
			}

			dp, err := base64urlTrailingPadding(jwk.DP)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "dp": %w`, KeyTypeRSA, err)
			}
			dq, err := base64urlTrailingPadding(jwk.DQ)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "dq": %w`, KeyTypeRSA, err)
			}
			qi, err := base64urlTrailingPadding(jwk.QI)
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "qi": %w`, KeyTypeRSA, err)
			}
			var oth []rsa.CRTValue
			var primes []*big.Int
			if len(jwk.OTH) > 0 {
				primes = make([]*big.Int, 2+len(jwk.OTH))
				primes[0] = new(big.Int).SetBytes(p)
				primes[1] = new(big.Int).SetBytes(q)
				// TODO Does each extra multi-prime need to be added to the slice of primes on the private key?
				oth = make([]rsa.CRTValue, len(jwk.OTH))
				for i, otherPrimes := range jwk.OTH {
					if otherPrimes.R == "" || otherPrimes.D == "" || otherPrimes.T == "" {
						return KeyWithMeta{}, fmt.Errorf(`%w: %s requires parameters "r", "d", and "t" for each "oth"`, ErrKeyUnmarshalParameter, KeyTypeRSA)
					}
					othD, err := base64urlTrailingPadding(otherPrimes.D)
					if err != nil {
						return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "d": %w`, KeyTypeRSA, err)
					}
					othT, err := base64urlTrailingPadding(otherPrimes.T)
					if err != nil {
						return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "t": %w`, KeyTypeRSA, err)
					}
					othR, err := base64urlTrailingPadding(otherPrimes.R)
					if err != nil {
						return KeyWithMeta{}, fmt.Errorf(`failed to decode %s key parameter "r": %w`, KeyTypeRSA, err)
					}
					primes[i+2] = new(big.Int).SetBytes(othR) // TODO This is incorrect
					oth[i] = rsa.CRTValue{
						Exp:   new(big.Int).SetBytes(othD),
						Coeff: new(big.Int).SetBytes(othT),
						R:     new(big.Int).SetBytes(othR),
					}
				}
			} else {
				primes = []*big.Int{
					new(big.Int).SetBytes(p),
					new(big.Int).SetBytes(q),
				}
			}
			privateKey := &rsa.PrivateKey{
				PublicKey: publicKey,
				D:         new(big.Int).SetBytes(d),
				Primes:    primes,
				Precomputed: rsa.PrecomputedValues{
					Dp:        new(big.Int).SetBytes(dp),
					Dq:        new(big.Int).SetBytes(dq),
					Qinv:      new(big.Int).SetBytes(qi),
					CRTValues: oth,
				},
			}
			err = privateKey.Validate()
			if err != nil {
				return KeyWithMeta{}, fmt.Errorf(`failed to validate %s key: %w`, KeyTypeRSA, err)
			}
			meta.Key = privateKey
		} else if !options.AsymmetricPrivate {
			meta.Key = &publicKey
		}
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
