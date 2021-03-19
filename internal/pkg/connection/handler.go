package handler

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
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