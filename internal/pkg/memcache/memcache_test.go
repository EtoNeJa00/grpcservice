package memcache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type memcacheTestSuite struct {
	suite.Suite
	ctx          context.Context
	mamcacheCont testcontainers.Container
	mc           MemCache
}

func (t *memcacheTestSuite) SetupSuite() {
	t.ctx = context.Background()

	defer func() {
		if err := recover(); err != nil {
			t.T().Error(err)
		}
	}()

	mcReq := testcontainers.ContainerRequest{
		Image:        "memcached",
		ExposedPorts: []string{"11211/tcp"},
		Name:         "memcached-test",
	}

	c, err := testcontainers.GenericContainer(t.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: mcReq,
		Started:          true,
	})

	t.Require().NoError(err)
	time.Sleep(5 * time.Second)

	ip, err := c.Host(t.ctx)
	t.Require().NoError(err)

	port, err := c.MappedPort(t.ctx, "11211")
	t.Require().NoError(err)

	t.mamcacheCont = c

	t.mc, err = NewMemcache(ip + ":" + port.Port())
	t.Require().NoError(err)
}

func (t *memcacheTestSuite) TearDownSuite() {
	if t.mamcacheCont != nil {
		_ = t.mamcacheCont.Terminate(t.ctx)
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(memcacheTestSuite))
}

func (t *memcacheTestSuite) TestMC() {
	data := []byte("123")

	resSet, err := t.mc.Set(data)
	t.Require().NoError(err)

	resGet, err := t.mc.Get(resSet)
	t.Require().NoError(err)
	t.Require().Equal(resGet[0], data)

	err = t.mc.Delete(resSet)
	t.Require().NoError(err)

	_, err = t.mc.Get(resSet)
	t.ErrorIs(err, ErrNotFound)
}
