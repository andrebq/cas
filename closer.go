package cas

import "io"

type (
	pCloser struct {
		c      io.Closer
		closed bool
	}
)

func (p *pCloser) Close() error {
	err := p.c.Close()
	p.closed = err == nil
	return err
}
