package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	fDebug := flag.Bool("v", false, "show debug log")
	fIP := flag.String("ip", "127.0.0.1", "set IP address to scan")
	fMin := flag.Int("min", 1, "minimum port")
	fMax := flag.Int("max", 1024, "maximum port")
	fRead := flag.Int64("limit", 4096, "max bytes to read from each port")
	fTimeout := flag.Int("timeout", 1, "seconds to scan for")
	flag.Parse()

	debugOut := ioutil.Discard
	if *fDebug {
		debugOut = os.Stderr
	}

	s := &scanner{
		log:       log.New(os.Stderr, "scan: ", 0),
		debug:     log.New(debugOut, "debug: ", 0),
		dialer:    &net.Dialer{},
		readLimit: *fRead,
	}

	if *fDebug {
		s.debugListener()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*fTimeout)*time.Second)
	defer cancel()

	s.scan(ctx, *fIP, *fMin, *fMax)
}

type scanner struct {
	debug     *log.Logger
	log       *log.Logger
	dialer    *net.Dialer
	readLimit int64
}

func (s *scanner) scan(ctx context.Context, ip string, min, max int) {
	var wg sync.WaitGroup
	for n := min; n <= max; n++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			addr := ip + ":" + strconv.Itoa(port)
			res, err := s.test(ctx, addr)
			if err != nil {
				s.debug.Printf("%s: %v", addr, err)
				return
			}
			s.log.Printf("read from %s: %s", addr, string(res))
		}(n)
	}
	wg.Wait()
}

func (s *scanner) test(ctx context.Context, addr string) ([]byte, error) {
	conn, err := s.dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetReadDeadline(deadline); err != nil {
			return nil, err
		}
	}
	s.log.Printf("connected to %s", addr)
	return ioutil.ReadAll(&io.LimitedReader{conn, s.readLimit})
}

func (s *scanner) debugListener() {
	l, err := net.Listen("tcp", ":123")
	if err != nil {
		s.debug.Fatal(err)
	}
	go func() {
		c, err := l.Accept()
		if err != nil {
			s.debug.Fatal(err)
		}
		defer c.Close()
		if _, err := fmt.Fprintln(c, "scan debug listener"); err != nil {
			s.debug.Fatal(err)
		}
	}()
}
