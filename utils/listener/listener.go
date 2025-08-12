package listener

import (
	"errors"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"time"
)

type TcpListenerArgs struct {
	Local string // 本地使用的地址
}

// TcpListener tcp 服务器
type TcpListener struct {
	cfg      *TcpListenerArgs
	quitChan chan interface{}
	wg       sync.WaitGroup
	Listener net.Listener
}

func NewTcpListener(cfg *TcpListenerArgs) *TcpListener {
	return &TcpListener{
		cfg:      cfg,
		quitChan: make(chan interface{}),
	}
}

// StartListen start tcp server. Notice: this method will not block
// callback will be called when new connection accepted
func (t *TcpListener) StartListen(callback func(conn net.Conn)) error {
	listen, err := net.Listen("tcp", t.cfg.Local)
	if err != nil {
		return err
	}

	t.Listener = listen

	t.wg.Add(1)
	go func() {
		defer t.wg.Done()

		for {
			conn, err := t.Listener.Accept()
			if err != nil {
				select {
				case <-t.quitChan:
					return
				default:
					log.Printf("TcpListener accept error: %v", err.Error())
					return
				}
			} else {
				t.wg.Add(1)
				go func() {
					defer t.wg.Done()
					defer func() {
						if e := recover(); e != nil {
							log.Printf("TcpListener connection handler crashed , acceptError : %v , \ntrace:%v", e, string(debug.Stack()))
						}
					}()
					// accept new connection, callback
					callback(conn)
				}()
			}
		}
	}()

	return nil
}

func (t *TcpListener) StopGracefully(wait time.Duration) error {
	close(t.quitChan)

	err := t.Listener.Close()
	if err != nil {
		log.Printf("TcpListener close tcp listener err: %v", err)
	}
	allExitChan := make(chan bool)
	go func() {
		// wait all goroutine exit
		t.wg.Wait()
		allExitChan <- true
	}()

	select {
	case <-time.After(wait):
		return errors.New("close tcp wait timeout")
	case <-allExitChan:
		return nil
	}
}
