package firethorn

import (
	"hash/fnv"
	_ "strings"
)

func (s *Storage) Shard(pool []interface{}, key string) interface{} {
	// from looking at the source, creating a new hash is basically free
	hashFn := fnv.New32a()
	hashFn.Write([]byte(key))

	index := uint32(hashFn.Sum32())
	index = index % uint32(len(pool))

	return pool[int(index)]
}
