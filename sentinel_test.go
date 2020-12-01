package sentinel_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/WuErPing/sentinel"
	"github.com/gomodule/redigo/redis"
)

func TestSentinel(t *testing.T) {
	sntnl := &sentinel.Sentinel{
		Addrs:      []string{":26379", ":26380", ":26381"},
		MasterName: "mymaster",
		Dial: func(addr string) (redis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := redis.DialTimeout("tcp", addr, timeout, timeout, timeout)
			assert.NoError(t, err)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	slaves, err := sntnl.Slaves()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(slaves), 2)
	for _, slave := range slaves {
		assert.NotEmpty(t, slave)
		t.Log(slave.Available(), slave)
	}

	slaveAddrs, err := sntnl.SlaveAddrs()
	assert.NoError(t, err)
	t.Log(slaveAddrs)

	availableSlaveAddrs, err := sntnl.AvailableSlaveAddrs()
	assert.NoError(t, err)
	t.Log(availableSlaveAddrs)
	assert.GreaterOrEqual(t, len(slaveAddrs), len(availableSlaveAddrs))
}

func newSentinelPool() *redis.Pool {
	sntnl := &sentinel.Sentinel{
		Addrs:      []string{":26379", ":26380", ":26381"},
		MasterName: "mymaster",
		Dial: func(addr string) (redis.Conn, error) {
			timeout := 500 * time.Millisecond
			c, err := redis.DialTimeout("tcp", addr, timeout, timeout, timeout)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   64,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			if err != nil {
				return nil, err
			}
			c, err := redis.Dial("tcp", masterAddr)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if !sentinel.TestRole(c, "master") {
				return errors.New("Role check failed")
			} else {
				return nil
			}
		},
	}
}

func TestPool(t *testing.T) {
	pool := newSentinelPool()
	fmt.Print(pool)
	c, err := pool.Dial()
	assert.NoError(t, err)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(c)
	}
	_, err = redis.String(c.Do("get", "key_non_existent"))
	assert.Error(t, err)
	defer c.Close()
}
