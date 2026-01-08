package proxy

import (
	"net"
)

func WriteToConn(conn net.Conn, buf []byte) error {
	_, err := conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func WriteStringToConn(conn net.Conn, msg string) error {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func ReadFromConn(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func ReadStringFromConn(conn net.Conn) (string, error) {
	buf, err := ReadFromConn(conn)
	return string(buf), err
}
