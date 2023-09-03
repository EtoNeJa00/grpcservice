package memcache

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"

	"github.com/google/uuid"
)

var (
	ErrNotStored = errors.New("not stored")
	ErrNotFound  = errors.New("not found")

	stored   = []byte("STORED\r\n")
	end      = []byte("END\r\n")
	eol      = "\r\n"
	notFound = []byte("NOT_FOUND\r\n")
)

type MemCache interface {
	Set(data []byte) (uuid.UUID, error)
	Get(id uuid.UUID) ([][]byte, error)
	Delete(id uuid.UUID) error
}

type memCache struct {
	addr *net.TCPAddr
}

func NewMemcache(servAddr string) (MemCache, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return nil, err
	}

	return memCache{addr: tcpAddr}, err
}

func (mc memCache) connect() (*net.TCPConn, error) {
	return net.DialTCP("tcp", nil, mc.addr)
}

func (mc memCache) Set(data []byte) (uuid.UUID, error) {
	conn, err := mc.connect()
	if err != nil {
		return uuid.UUID{}, err
	}

	defer func() {
		errC := conn.Close()
		if err == nil {
			err = errC
		}
	}()

	id := uuid.New()

	_, err = fmt.Fprintf(conn, "set %s 0 0 %d\r\n%s\r\n", id.String(), len(data), data)
	if err != nil {
		return uuid.UUID{}, err
	}

	connReader := bufio.NewReader(conn)

	reply, err := connReader.ReadSlice('\n')
	if err != nil {
		return uuid.UUID{}, err
	}

	if bytes.Equal(reply, stored) {
		return id, nil
	}

	return uuid.UUID{}, ErrNotStored
}

func (mc memCache) Get(id uuid.UUID) ([][]byte, error) {
	conn, err := mc.connect()
	if err != nil {
		return nil, err
	}

	defer func() {
		errC := conn.Close()
		if err == nil {
			err = errC
		}
	}()

	_, err = fmt.Fprintf(conn, "get %s\r\n", id.String())
	if err != nil {
		return nil, err
	}

	var results [][]byte

	connReader := bufio.NewReader(conn)
	firstLine := true

	for {
		line, err := connReader.ReadSlice('\n')
		if err != nil {
			return nil, err
		}

		if bytes.Equal(line, end) {
			break
		}

		if firstLine {
			firstLine = false

			continue
		}

		results = append(results, bytes.Trim(line, eol))
	}

	if len(results) == 0 {
		return nil, ErrNotFound
	}

	return results, nil
}

func (mc memCache) Delete(id uuid.UUID) error {
	conn, err := mc.connect()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(conn, "delete %s\r\n", id.String())
	if err != nil {
		return err
	}

	connReader := bufio.NewReader(conn)

	reply, err := connReader.ReadSlice('\n')
	if err != nil {
		return err
	}

	if bytes.Equal(reply, notFound) {
		return ErrNotFound
	}

	return nil
}
