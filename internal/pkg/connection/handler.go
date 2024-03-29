package connection

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func PerformHandshake(conn net.Conn, isServer bool) error {
	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)

	if isServer {
		// Write greeting
		w.WriteString("HI")
		w.WriteByte(0x0)
		w.Flush()
	}

	// Wait for HI contact or timeout
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	s, err := r.ReadString(byte(0x0))
	if err != nil {
		log.Println("received error", s, err.Error())
		return err
	}

	_ = conn.SetReadDeadline(time.Time{})

	if cmd, err := Decode(s); err != nil {
		return err
	} else if msg := cmd[0]; msg != "HI" {
		return errors.New(fmt.Sprintf("expected HI, received: %s", msg))
	}

	if !isServer {
		// Respond to greeting
		w.WriteString("HI")
		w.WriteByte(0x0)
		w.Flush()
	}

	return nil
}

func ReadConn(conn net.Conn, msg, errChan chan string) {
	r := bufio.NewReader(conn)

	for {
		s, err := r.ReadString(byte(0x0))
		if err != nil {
			errChan <- err.Error()
			close(errChan)
			close(msg)
			return
		}
		msg <- s
	}
}

func WriteConn(conn net.Conn, msg, errChan chan string) {
	w := bufio.NewWriter(conn)

	for {
		select {
		case s, ok := <-msg:
			if !ok {
				close(errChan)
				return
			}
			_, err := w.WriteString("MSG ")
			if err != nil {
				errChan <- err.Error()
			}
			_, err = w.WriteString(s)
			if err != nil {
				errChan <- err.Error()
			}
			err = w.WriteByte(0x0)
			if err != nil {
				errChan <- err.Error()
			}
			err = w.Flush()
			if err != nil {
				errChan <- err.Error()
			}
		}
	}
}

func Decode(s string) ([]string, error) {
	// Trim 0x0
	s = s[:len(s)-1]

	//log.Println("Decoding", s)

	cmd := strings.Split(s, " ")
	if len(cmd) < 1 {
		return []string{}, errors.New(fmt.Sprintf("Invalid cmd received: %s", cmd))
	}
	return cmd, nil
}

func KeepAlive(conn net.Conn, ka chan struct{}, errChan chan string, isServer bool) {
	var rxDelay, txDelay, netDelay time.Duration
	if isServer {
		rxDelay = time.Duration(0)
		txDelay = time.Second * 5
		netDelay = time.Second * 1
	} else {
		rxDelay = time.Second * 5
		txDelay = time.Duration(0)
		netDelay = time.Second * 1
	}

	ticker := time.NewTicker(rxDelay + netDelay)
	defer ticker.Stop()
	w := bufio.NewWriter(conn)

	writeKeepAlive := func(delay time.Duration) {
		ticker.Stop()
		time.Sleep(delay)
		//log.Println("Sending KA")
		w.WriteString("KA")
		w.WriteByte(0x0)
		w.Flush()
		ticker.Reset(rxDelay + netDelay)
	}

	if isServer {
		//log.Println("Server sending initial KA")
		writeKeepAlive(txDelay)
	}

	for {
		select {
		case <-ticker.C:
			errChan <- "timed out waiting for KA"
			close(errChan)
			return
		case _, ok := <-ka:
			if !ok {
				errChan <- "received on closed keep alive channel"
				close(errChan)
				return
			}
			//log.Println("Received KA")
			writeKeepAlive(txDelay)
		}
	}
}

func HandleConnection(conn net.Conn, isServer bool) *Connection {

	r := make(chan string)
	w := make(chan string)
	e := make(chan string)

	go func() {
		err := PerformHandshake(conn, isServer)
		if err != nil {
			e <- err.Error()
			return
		}

		keepAliveChan := make(chan struct{})
		keepAliveErrChan := make(chan string)
		go KeepAlive(conn, keepAliveChan, keepAliveErrChan, isServer)

		readChan := make(chan string)
		readErrChan := make(chan string)
		go ReadConn(conn, readChan, readErrChan)

		writeChan := make(chan string)
		writeErrChan := make(chan string)
		go WriteConn(conn, writeChan, writeErrChan)

		// Wait for cmd or failed keepAlive
		var s string
		var ok bool
		for {
			select {
			case s, ok := <-keepAliveErrChan:
				if !ok {
					e <- "closed keep alive error channel from error"
					close(keepAliveChan)
					return
				}
				e <- s
				return
			case msg, ok := <-w:
				if !ok {
					e <- "write message channel closed but shouldn't be"
					return
				}
				writeChan <- msg
				continue
			case s, ok := <-writeErrChan:
				if !ok {
					e <- "closed write error channel from error"
					return
				}
				log.Println("error writing message", s)
				continue
			case s, ok = <-readChan:
				if !ok {
					e <- "closed message channel from error"
					return
				}
				break
				//log.Println("reading...")
			case s, ok := <-readErrChan:
				if !ok {
					e <- "closed message error channel from error"
					return
				}
				e <- s
				return
			}

			cmd, err := Decode(s)
			if err != nil {
				log.Println(err)
				continue
			}

			switch cmd[0] {
			case "HI":
				log.Println("unexpected HI after handshake")
			case "KA":
				keepAliveChan <- struct{}{}
			case "MSG":
				if len(cmd) > 1 {
					msg := strings.Join(cmd[1:], " ")
					r <- msg
				}
			default:
				log.Println(s, "is unhandled")
			}
		}
	}()

	return &Connection{r, w, e}
}

func ReadInput(input chan string) {
	r := bufio.NewReader(os.Stdin)
	for {
		msg, _ := r.ReadString('\n')
		// convert CRLF to LF
		msg = strings.Replace(msg, "\n", "", -1)
		input <- msg
	}
}

type Connection struct {
	R, W, E chan string
}
