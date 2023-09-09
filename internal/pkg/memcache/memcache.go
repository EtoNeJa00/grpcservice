package memcache

import (
	"bytes"
	"context"
	"errors"
	"net"
	"strconv"

	"github.com/buraksezer/connpool"
	"github.com/google/uuid"
)

var (
	ErrNotStored = errors.New("not stored")
	ErrNotFound  = errors.New("not found")

	stored   = []byte("STORED\r\n")
	end      = []byte("END\r\n")
	eol      = []byte("\r\n")
	eolnb    = []byte("\n")
	notFound = []byte("NOT_FOUND\r\n")
)

type MemCache interface {
	Set(data []byte) (uuid.UUID, error)
	Get(id uuid.UUID) ([][]byte, error)
	Delete(id uuid.UUID) error
}

type memCache struct {
	p connpool.Pool
}

func NewMemcache(servAddr string) (MemCache, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		return nil, err
	}

	p, err := connpool.NewChannelPool(5, 30, func() (net.Conn, error) {
		return net.DialTCP("tcp", nil, tcpAddr)
	})

	return memCache{p: p}, err
}

func (mc memCache) Set(data []byte) (uuid.UUID, error) {
	conn, err := mc.p.Get(context.Background())
	if err != nil {
		return uuid.UUID{}, err
	}

	defer func() {
		errC := conn.Close()
		if errC != nil {
			err = errC
		}
	}()

	id := uuid.New()

	_, err = conn.Write([]byte("set " + id.String() + " 0 0 " + strconv.Itoa(len(data)) + "\r\n" + string(data) + "\r\n"))
	if err != nil {
		return uuid.UUID{}, err
	}

	msg, err := mc.readMsg(conn, make([]byte, 8))
	if err != nil {
		return uuid.UUID{}, err
	}

	if bytes.Equal(msg, stored) {
		return id, nil
	}

	return uuid.UUID{}, ErrNotStored
}

func (mc memCache) readMsg(conn net.Conn, buffer []byte) ([]byte, error) {
	var msg []byte

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return nil, err
		}

		msg = append(msg, buffer[:n]...)

		if bytes.Equal(msg[len(msg)-2:], eol) {
			break
		}
	}

	return msg, nil
}

func (mc memCache) readMultilineMsg(conn net.Conn, buffer []byte) ([]byte, error) {
	var msg []byte

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return nil, err
		}

		msg = append(msg, buffer[:n]...)

		if bytes.Equal(msg[len(msg)-5:], end) {
			break
		}
	}

	return msg, nil
}

func (mc memCache) Get(id uuid.UUID) ([][]byte, error) {
	conn, err := mc.p.Get(context.Background())
	if err != nil {
		return nil, err
	}

	defer func() {
		errC := conn.Close()
		if err == nil {
			err = errC
		}
	}()

	_, err = conn.Write([]byte("get " + id.String() + "\r\n"))
	if err != nil {
		return nil, err
	}

	msg, err := mc.readMultilineMsg(conn, make([]byte, 5))
	if err != nil {
		return nil, err
	}

	response := bytes.Split(msg, eol)

	if len(response) == 0 || len(response) < 3 {
		return nil, ErrNotFound
	}

	return response[1 : len(response)-1], nil
}

func (mc memCache) Delete(id uuid.UUID) error {
	conn, err := mc.p.Get(context.Background())
	if err != nil {
		return err
	}

	defer func() {
		errC := conn.Close()
		if err == nil {
			err = errC
		}
	}()

	_, err = conn.Write([]byte("delete " + id.String() + "\r\n"))
	if err != nil {
		return err
	}

	msg, err := mc.readMsg(conn, make([]byte, 8))
	if err != nil {
		return err
	}

	if bytes.Equal(msg, notFound) {
		return ErrNotFound
	}

	return nil
}
