package kvsrv
import (
	"log"
	"sync"

	"6.5840/kvsrv1/rpc"
	"6.5840/labrpc"
	tester "6.5840/tester1"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type valueEntry struct {
	Value   string
	Version rpc.Tversion
}

type KVServer struct {
	mu    sync.Mutex
	store map[string]valueEntry
}

func (kv *KVServer) ServiceName() string {
	return "KVServer"
}

func MakeKVServer() *KVServer {
	kv := &KVServer{
		store: make(map[string]valueEntry),
	}
	return kv
}

func (kv *KVServer) Get(args *rpc.GetArgs, reply *rpc.GetReply) {

	kv.mu.Lock()
	defer kv.mu.Unlock()

	entry, ok := kv.store[args.Key]
	if !ok {
		reply.Err = rpc.ErrNoKey
		return
	}

	reply.Value = entry.Value
	reply.Version = entry.Version
	reply.Err = rpc.OK
}

func (kv *KVServer) Put(args *rpc.PutArgs, reply *rpc.PutReply) {

	kv.mu.Lock()
	defer kv.mu.Unlock()

	entry, ok := kv.store[args.Key]

	if !ok {
		if args.Version == 0 {
			kv.store[args.Key] = valueEntry{
				Value:   args.Value,
				Version: 1,
			}
			reply.Err = rpc.OK
		} else {
			reply.Err = rpc.ErrNoKey
		}
		return
	}

	if entry.Version != args.Version {
		reply.Err = rpc.ErrVersion
		return
	}

	kv.store[args.Key] = valueEntry{
		Value:   args.Value,
		Version: entry.Version + 1,
	}
	reply.Err = rpc.OK
}

func (kv *KVServer) Kill() {
}

func StartKVServer(ends []*labrpc.ClientEnd, gid tester.Tgid, srv int, persister *tester.Persister) []tester.IService {
	kv := MakeKVServer()
	return []tester.IService{kv}
}
