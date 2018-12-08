package the_platinum_searcher

import (
	"io"
	"sync"
)

type printer struct {
	in        chan match
	mu        *sync.Mutex
	opts      Option
	formatter formatPrinter
	done      chan struct{}
}

func newPrinter(
	pattern pattern,
	out,
	errorWriter io.Writer,
	opts Option,
) printer {
	p := printer{
		in:        make(chan match, 200),
		mu:        new(sync.Mutex),
		opts:      opts,
		formatter: newFormatPrinter(pattern, out, errorWriter, opts),
		done:      make(chan struct{}),
	}

	go p.loop()
	return p
}

func (p printer) print(match match) {
	if match.size() == 0 {
		return
	}

	p.in <- match
}

func (p printer) loop() {
	defer func() {
		p.done <- struct{}{}
	}()

	for match := range p.in {
		p.formatter.print(match)
	}
}

func (p printer) printError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.formatter.printError(err)
}
