// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backendjwt

import (
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/net/mwjwt"
	"github.com/corestoreio/csfw/util/csjwt"
	"github.com/corestoreio/csfw/util/errors"
)

// Default creates new mwjwt.Option slice with the default configuration
// structure and a noop encryptor/decryptor IF no option arguments have been
// provided. It panics on error, so us it only during the app init phase.
func Default(opts ...cfgmodel.Option) mwjwt.ScopedOptionFunc {
	cfgStruct, err := NewConfigStructure()
	if err != nil {
		panic(err)
	}
	if len(opts) == 0 {
		opts = append(opts, cfgmodel.WithEncryptor(cfgmodel.NoopEncryptor{}))
	}

	return PrepareOptions(New(cfgStruct, opts...))
}

// PrepareOptions creates a closure around the type Backend. The closure will
// be used during a scoped request to figure out the configuration depending on
// the incoming scope. An option array will be returned by the closure.
func PrepareOptions(be *Backend) mwjwt.ScopedOptionFunc {

	return func(sg config.ScopedGetter) []mwjwt.Option {
		var opts [6]mwjwt.Option
		var i int
		scp, id := sg.Scope()

		off, err := be.NetJwtDisabled.Get(sg)
		if err != nil {
			opts[i] = func(s *mwjwt.Service) error {
				return errors.Wrap(err, "[backendcors] NetJwtDisabled.Get")
			}
			return opts[:]
		}
		opts[i] = mwjwt.WithDisable(scp, id, off)
		i++

		exp, err := be.NetJwtExpiration.Get(sg)
		if err != nil {
			opts[i] = func(s *mwjwt.Service) error {
				return errors.Wrap(err, "[backendjwt] NetJwtExpiration.Get")
			}
			return opts[:]
		}
		opts[i] = mwjwt.WithExpiration(scp, id, exp)
		i++

		skew, err := be.NetJwtSkew.Get(sg)
		if err != nil {
			opts[i] = func(s *mwjwt.Service) error {
				return errors.Wrap(err, "[backendjwt] NetJwtSkew.Get")
			}
			return opts[:]
		}
		opts[i] = mwjwt.WithSkew(scp, id, skew)
		i++

		isJTI, err := be.NetJwtEnableJTI.Get(sg)
		if err != nil {
			opts[i] = func(s *mwjwt.Service) error {
				return errors.Wrap(err, "[backendjwt] NetJwtEnableJTI.Get")
			}
			return opts[:]
		}
		opts[i] = mwjwt.WithTokenID(scp, id, isJTI)
		i++

		signingMethod, err := be.NetJwtSigningMethod.Get(sg)
		if err != nil {
			opts[i] = func(s *mwjwt.Service) error {
				return errors.Wrap(err, "[backendjwt] NetJwtSigningMethod.Get")
			}
			return opts[:]
		}

		var key csjwt.Key

		switch signingMethod.Alg() {
		case csjwt.RS256, csjwt.RS384, csjwt.RS512:

			rsaKey, err := be.NetJwtRSAKey.Get(sg)
			if err != nil {
				opts[i] = func(s *mwjwt.Service) error {
					return errors.Wrap(err, "[backendjwt] NetJwtRSAKey.Get")
				}
				return opts[:]
			}
			rsaPW, err := be.NetJwtRSAKeyPassword.Get(sg)
			if err != nil {
				opts[i] = func(s *mwjwt.Service) error {
					return errors.Wrap(err, "[backendjwt] NetJwtRSAKeyPassword.Get")
				}
				return opts[:]
			}
			key = csjwt.WithRSAPrivateKeyFromPEM(rsaKey, rsaPW)

		case csjwt.ES256, csjwt.ES384, csjwt.ES512:

			ecdsaKey, err := be.NetJwtECDSAKey.Get(sg)
			if err != nil {
				opts[i] = func(s *mwjwt.Service) error {
					return errors.Wrap(err, "[backendjwt] NetJwtECDSAKey.Get")
				}
				return opts[:]
			}
			ecdsaPW, err := be.NetJwtECDSAKeyPassword.Get(sg)
			if err != nil {
				opts[i] = func(s *mwjwt.Service) error {
					return errors.Wrap(err, "[backendjwt] NetJwtECDSAKeyPassword.Get")
				}
				return opts[:]
			}
			key = csjwt.WithECPrivateKeyFromPEM(ecdsaKey, ecdsaPW)

		case csjwt.HS256, csjwt.HS384, csjwt.HS512:

			password, err := be.NetJwtHmacPassword.Get(sg)
			if err != nil {
				opts[i] = func(s *mwjwt.Service) error {
					return errors.Wrap(err, "[backendjwt] NetJwtHmacPassword.Get")
				}
				return opts[:]
			}
			key = csjwt.WithPassword(password)

		default:
			opts[i] = func(s *mwjwt.Service) error {
				return errors.Errorf("[mwjwt] Unknown signing method: %q", signingMethod.Alg())
			}
			return opts[:]
		}

		// WithSigningMethod must be added at the end of the slice to overwrite default signing methods
		opts[i] = mwjwt.WithKey(scp, id, key)
		i++
		opts[i] = mwjwt.WithSigningMethod(scp, id, signingMethod)
		return opts[:]
	}
}
