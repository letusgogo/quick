package listener

import (
	"io"
	"log"
	"net"
	"runtime/debug"
	"time"
)

func IoBind(dst io.ReadWriteCloser, src io.ReadWriteCloser) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("bind crashed %s", err)
		}
	}()
	errCh := make(chan error, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("bind crashed %s", err)
			}
		}()
		err := ioCopy(src, dst)
		errCh <- err
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("bind crashed %s", err)
			}
		}()
		err := ioCopy(dst, src)
		errCh <- err
	}()

	if err := <-errCh; err != nil && err != io.EOF {
		return err
	}
	return nil
}

func ioCopy(dst io.ReadWriter, src io.ReadWriter) (err error) {
	defer func() {
		if e := recover(); e != nil {
		}
	}()
	buf := LeakyBuffer.Get()
	defer LeakyBuffer.Put(buf)
	n := 0
	for {
		n, err = src.Read(buf)
		if n > 0 {
			if n > len(buf) {
				n = len(buf)
			}
			if _, e := dst.Write(buf[0:n]); e != nil {
				return e
			}
		}
		if err != nil {
			return
		}
	}
}

func Close(conn net.Conn) {
	_ = conn.SetDeadline(time.Now().Add(time.Millisecond * 100))
	if err := conn.Close(); err != nil {
		log.Printf("http close error, err: %v \nstack: %v", err, string(debug.Stack()))
	}
}
