package main

import (
	"strconv"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"COMMAND": func(args []Value) Value { return Value{typ: "null"} }, // func to handle redis-cli's 'COMMAND' command during initialization
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"DEL":     del,
	"HDEL":    hdel,
	"EXISTS":  exists,
	"HEXISTS": hexists,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}
	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}

}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}
	key := args[0].bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	_, ok := HSETs[hash]
	if !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	vals, ok := HSETs[hash]

	if !ok {
		return Value{typ: "null"}
	}

	values := []Value{}
	for k, v := range vals {
		values = append(values, Value{typ: "bulk", bulk: k})
		values = append(values, Value{typ: "bulk", bulk: v})
	}
	HSETsMu.RUnlock()

	return Value{typ: "array", array: values}

}

func del(args []Value) Value {
	deleted := 0
	SETsMu.Lock()
	defer SETsMu.Unlock()
	for _, arg := range args {
		key := arg.bulk
		_, ok := SETs[key]
		if !ok {
			continue
		}
		delete(SETs, key)
		deleted += 1
	}
	return Value{typ: "string", str: "(integer) " + strconv.Itoa(deleted)}

}

func hdel(args []Value) Value {
	if len(args) < 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hdel' command"}
	}
	hash := args[0].bulk
	deleted := 0
	HSETsMu.Lock()
	defer HSETsMu.Unlock()
	for _, arg := range args[1:] {
		key := arg.bulk
		_, ok := HSETs[hash][key]
		if !ok {
			continue
		}
		delete(HSETs[hash], key)
		deleted += 1
	}
	return Value{typ: "string", str: "(integer) " + strconv.Itoa(deleted)}
}

func exists(args []Value) Value {
	key_exists := 0
	SETsMu.RLock()
	defer SETsMu.RUnlock()
	for _, arg := range args {
		key := arg.bulk
		_, ok := SETs[key]
		if !ok {
			continue
		}
		key_exists += 1
	}
	return Value{typ: "string", str: "(integer) " + strconv.Itoa(key_exists)}
}

func hexists(args []Value) Value {
	if len(args) < 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hexists' command"}
	}
	hash := args[0].bulk
	key_exists := 0
	HSETsMu.RLock()
	defer HSETsMu.RUnlock()
	for _, arg := range args[1:] {
		key := arg.bulk
		_, ok := HSETs[hash][key]
		if !ok {
			continue
		}
		key_exists += 1
	}
	return Value{typ: "string", str: "(integer) " + strconv.Itoa(key_exists)}
}
