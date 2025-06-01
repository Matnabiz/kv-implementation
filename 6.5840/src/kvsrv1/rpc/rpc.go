package rpc

import "6.5840/labgob"

type Err string

const (
	OK         = "OK"
	ErrNoKey   = "ErrNoKey"
	ErrVersion = "ErrVersion"
	ErrMaybe   = "ErrMaybe"

	ErrWrongLeader = "ErrWrongLeader"
	ErrWrongGroup  = "ErrWrongGroup"
)

type Tversion uint64

type PutArgs struct {
	Key     string
	Value   string
	Version Tversion
}

type PutReply struct {
	Err Err
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Value   string
	Version Tversion
	Err     Err
}

func init() {
	labgob.Register(&GetArgs{})
	labgob.Register(&GetReply{})
	labgob.Register(&PutArgs{})
	labgob.Register(&PutReply{})
}
