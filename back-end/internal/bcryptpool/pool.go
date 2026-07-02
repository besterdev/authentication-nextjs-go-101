package bcryptpool

import "golang.org/x/crypto/bcrypt"

type Pool struct {
	sem chan struct{}
}

func New(maxConcurrent int) *Pool {
	if maxConcurrent < 1 {
		maxConcurrent = 1
	}
	return &Pool{sem: make(chan struct{}, maxConcurrent)}
}

func (p *Pool) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	p.sem <- struct{}{}
	defer func() { <-p.sem }()
	return bcrypt.GenerateFromPassword(password, cost)
}

func (p *Pool) CompareHashAndPassword(hashedPassword, password []byte) error {
	p.sem <- struct{}{}
	defer func() { <-p.sem }()
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
