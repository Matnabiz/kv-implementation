package kvsrv

import (
	"time"

	"6.5840/kvsrv1/rpc"
	kvtest "6.5840/kvtest1"
	tester "6.5840/tester1"
)

type Clerk struct {
	clnt   *tester.Clnt
	server string
}

func MakeClerk(clnt *tester.Clnt, server string) kvtest.IKVClerk {
	return &Clerk{clnt: clnt, server: server}
}

func (ck *Clerk) Get(key string) (string, rpc.Tversion, rpc.Err) {
	args := rpc.GetArgs{Key: key}
	var reply rpc.GetReply

	for {
		ok := ck.clnt.Call(ck.server, "KVServer.Get", &args, &reply)
		if ok {
			if reply.Err == rpc.OK || reply.Err == rpc.ErrNoKey {
				return reply.Value, reply.Version, reply.Err
			}
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (ck *Clerk) Put(key, value string, version rpc.Tversion) rpc.Err {
	args := rpc.PutArgs{Key: key, Value: value, Version: version}
	var reply rpc.PutReply

	attemptCount := 0

	for {
		ok := ck.clnt.Call(ck.server, "KVServer.Put", &args, &reply)
		attemptCount++

		if ok {
			switch reply.Err {
			case rpc.OK:
				return rpc.OK
			case rpc.ErrVersion:
				if attemptCount == 1 {
					return rpc.ErrVersion
				}
				return rpc.ErrMaybe
			case rpc.ErrNoKey:
				return rpc.ErrNoKey
			case rpc.ErrWrongLeader, rpc.ErrWrongGroup:
			default:
			}
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}
