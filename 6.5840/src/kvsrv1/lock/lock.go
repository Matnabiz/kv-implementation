package lock

import (
	"math/rand"
	"time"

	"6.5840/kvsrv1/rpc"
	kvtest "6.5840/kvtest1"
)

type Lock struct {
	ck  kvtest.IKVClerk
	key string
	id  string
}

func MakeLock(ck kvtest.IKVClerk, lockKey string) *Lock {
	return &Lock{
		ck:  ck,
		key: lockKey,
		id:  kvtest.RandValue(8),
	}
}

func (lk *Lock) Acquire() {
	backoff := 5 * time.Millisecond

	for {
		val, ver, err := lk.ck.Get(lk.key)

		var baseVer rpc.Tversion
		switch err {
		case rpc.OK:
			if val == lk.id {
				return
			}
			if val != "" {
				lk.sleep(backoff)
				backoff = lk.nextBackoff(backoff)
				continue
			}
			baseVer = ver

		case rpc.ErrNoKey:
			baseVer = 0

		default:
			lk.sleep(backoff)
			backoff = lk.nextBackoff(backoff)
			continue
		}
		putErr := lk.ck.Put(lk.key, lk.id, baseVer)
		switch putErr {
		case rpc.OK:
			return

		case rpc.ErrVersion, rpc.ErrNoKey:
			lk.sleep(backoff)
			backoff = lk.nextBackoff(backoff)

		case rpc.ErrMaybe:
			confirmVal, _, confirmErr := lk.ck.Get(lk.key)
			if confirmErr == rpc.OK && confirmVal == lk.id {
				return
			}
			lk.sleep(backoff)
			backoff = lk.nextBackoff(backoff)

		default:
			lk.sleep(backoff)
			backoff = lk.nextBackoff(backoff)
		}
	}
}

func (lk *Lock) Release() {
	val, ver, err := lk.ck.Get(lk.key)
	if err != rpc.OK || val != lk.id {
		return
	}
	_ = lk.ck.Put(lk.key, "", ver)
}

func (lk *Lock) sleep(d time.Duration) {
	jitter := time.Duration(rand.Intn(5)) * time.Millisecond
	time.Sleep(d + jitter)
}

func (lk *Lock) nextBackoff(curr time.Duration) time.Duration {
	if curr < 100*time.Millisecond {
		return curr * 2
	}
	return curr
}
