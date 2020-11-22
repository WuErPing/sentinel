package sentinel_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/WuErPing/sentinel"
	"github.com/gomodule/redigo/redis"
)

func TestSentinelPool(t *testing.T) {
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
	pool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   64,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			masterAddr, err := sntnl.MasterAddr()
			assert.NoError(t, err)
			if err != nil {
				return nil, err
			}
			c, err := redis.Dial("tcp",
				masterAddr,
				redis.DialUsername(""),
				redis.DialPassword("test@hago"))
			assert.NoError(t, err)
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

	c, err := pool.Dial()
	assert.NoError(t, err)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(c)
	}
	t.Log(pool.ActiveCount())
	slaves, err := sntnl.Slaves()
	assert.NoError(t, err)
	for _, slave := range slaves {
		t.Log(slave)
	}
	r, err := redis.String(c.Do("get", "key_andy"))
	assert.NoError(t, err)
	defer c.Close()
	t.Log(r)
}
