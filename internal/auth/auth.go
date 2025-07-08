package auth

import (
	"errors"
	"sync"
	"time"

	"devices_api/internal/krypto"
)

var (
	ErrDuplicateUser      = errors.New("duplicate user")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// // Emailer is used to send templated emails.
// type Emailer interface {
// 	Send(ctx context.Context, template string, to email.Address, data interface{}) error
// }

// ErrFunc is a function that handles errors.
type ErrFunc func(error)

// ServiceConfig is the configuration for the Service. Some methods run in seperate goroutines,
// it is up to the caller to wait for these methods to finish. This can be done by calling the
// Wait method.
type ServiceConfig struct {
	// WorkerTimeout is the max duration worker goroutines are allowed
	// to take befor they are cancelled.
	WorkerTimeout time.Duration
	// TokenExpirty is the duration a token is valid.
	TokenExpiry time.Duration
}

// Service is the type that provides the main rules for
// authentication.
type Service struct {
	//store      Store
	//emailer    Emailer
	wg         *sync.WaitGroup
	errHandler ErrFunc
	cfg        ServiceConfig

	// comparisonHash is used to compare passwords when no user was found.
	comparisonHash krypto.Argon2Hash

	// NowFunc is used to get the current time.
	// Exposed for testing purposes.
	NowFunc func() time.Time
}

// NewService creates a new Service.
// func NewService(s Store, emailer Emailer, errHandler ErrFunc, cfg ServiceConfig) (*Service, error) {
// 	tok, err := krypto.GenerateToken()
// 	if err != nil {
// 		return nil, err
// 	}

// 	hash, err := krypto.HashArgon2(tok[:])
// 	if err != nil {
// 		return nil, err
// 	}

// 	svc := &Service{
// 		// store:          s,
// 		// emailer:        emailer,
// 		wg:             &sync.WaitGroup{},
// 		errHandler:     errHandler,
// 		cfg:            cfg,
// 		comparisonHash: hash,
// 		NowFunc:        time.Now,
// 	}

// 	return svc, nil
// }

// Wait waits for all open workers to finish.
func (s *Service) Wait() {
	s.wg.Wait()
}
