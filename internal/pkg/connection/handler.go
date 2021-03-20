package handler

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

func ReadConn(conn net.Conn, msg chan string) {
	r := bufio.NewReader(conn)

	for {
		s, err := r.ReadString(byte(0x0))
		if err != nil {
			log.Println("received error", s, err.Error())
			close(msg)
			return
		}
		msg <- s
	}
}


func Decode(s string) ([]string, error) {
	// Trim 0x0
	s = s[:len(s)-1]

	log.Println("Decoding", s)

	cmd := strings.Split(s, " ")
	if len(cmd) < 1 {
		return []string{}, errors.New(fmt.Sprintf("Invalid cmd received: %s", cmd))
	}
	return cmd, nil
}

func KeepAlive(conn net.Conn, ka chan struct{}, rxDelay, txDelay, networkDelay time.Duration) {
	ticker := time.NewTicker(rxDelay + networkDelay)
	defer ticker.Stop()
	w := bufio.NewWriter(conn)

	writeKeepAlive := func(delay time.Duration) {
		ticker.Stop()
		time.Sleep(delay)
		log.Println("Sending KA")
		w.WriteString("KA")
		w.WriteByte(0x0)
		w.Flush()
		ticker.Reset(rxDelay + networkDelay)
	}

	for {
		select {
		case <- ticker.C:
			ka <- struct{}{}
			return
		case <- ka:
			log.Println("Received KA")
			writeKeepAlive(txDelay)
		}
	}
}

func HandleConnection(conn net.Conn, isServer bool) (
	chan string,
	chan string,
	chan string,
) {

	r := make(chan string)
	w := make(chan string)
	e := make(chan string)

	go func() {
		err := PerformHandshake(conn, isServer)
		if err != nil {
			log.Println(err)
			return
		}

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
		keepAlive := make(chan struct{})
		go KeepAlive(conn, keepAlive, rxDelay, txDelay, netDelay)
		if isServer {
			keepAlive <- struct{}{}
		}

		msgChan := make(chan string)
		go ReadConn(conn, msgChan)

		w := bufio.NewWriter(conn)
		go SendMessage(w)

		// Wait for cmd or failed keepAlive
		var s string
		var ok bool
		for {
			select {
			case _, ok := <-keepAlive:
				if !ok {
					log.Println("Keep Alive channel closed but shouldn't be!")
					return
				}
				close(keepAlive)
				log.Println("failed to receive KA")
				return
			case s, ok = <-msgChan:
				log.Println("reading...")
				if !ok {
					log.Println("closed from error")
					return
				}
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
				keepAlive <- struct{}{}
			case "MSG":
				if len(cmd) > 1 {
					msg := strings.Join(cmd[1:], " ")
					r <- msg
					//b.TellEveryone(fmt.Sprintf("MSG %s %s", conn.RemoteAddr().String(), msg))
				}
				//if len(cmd) > 2 {
				//	from := cmd[1]
				//	msg := strings.Join(cmd[2:], " ")
				//	log.Printf("[%s]: %s", from, msg)
				//}
			default:
				log.Println(s, "is unhandled")
			}
		}
	}()

	return r, w, e
}

func SendMessage(w *bufio.Writer) {
	r := bufio.NewReader(os.Stdin)
	for {
		msg, _ := r.ReadString('\n')
		// convert CRLF to LF
		msg = strings.Replace(msg, "\n", "", -1)
		w.WriteString("MSG ")
		w.WriteString(msg)
		w.WriteByte(0x0)
		w.Flush()
	}
}
