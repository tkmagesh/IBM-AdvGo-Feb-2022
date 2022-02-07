package pool

import "io"

type Pool struct {
}

func New(factory func() (io.Closer, error), size int) (*Pool, error) {
	return &Pool{}, nil
}

func (p *Pool) Acquire() (io.Closer, error) {
	return nil, nil
}

func (p *Pool) Release(io.Closer) {

}

func (p *Pool) Close() {

}
