package firethorn

import (
	"bytes"
)

type Schema struct {
	Dimensions map[string]Dimension
}

type Dimension struct {
	Id        uint
	Key       string
	Hierarchy string
	Enums     []string
}

type Insert struct {
	Namespace  string
	Delta      int
	Dimensions [][]string
}
type InsertKeyGenerator []uint

type Query struct {
}
type QueryKeyGenerator []uint

// computes the number of keys referenced in an insert
func (s *Schema) InsertKeyCount(i *Insert) uint {
	count := uint(0)
	for _, column := range i.Dimensions {
		// "+1" to account for the `root` of every dimension, representing
		// the lack of a specification for that dimension
		//
		// note how this correctly accounts for dimensions which are left
		// unspecified by the insert as well as dimensions which are only
		// partially specified
		dimensionCount := len(s.Dimensions) * (1 + len(column))
		count *= uint(dimensionCount)
	}

	return count
}

// creates a generator that iteratively creates all keys referenced in an
// insert
func (ins *Insert) KeyGenerator() InsertKeyGenerator {
	return InsertKeyGenerator(make([]uint, len(ins.Dimensions)))
}

// updates the generator to give the next key on GetKey
func (igen *InsertKeyGenerator) NextKey(insert Insert) *InsertKeyGenerator {
	// start at the first dimension, increment
	(*igen)[0]++

	// now iterate over the rest, resolving overflows
	lastOverflow := -1
	for i := range insert.Dimensions {
		if (*igen)[i] <= uint(len(insert.Dimensions[i])) {
			break
		} else {
			(*igen)[i+1]++
			lastOverflow = i
		}
	}

	// zero out up to the overflow
	for i := 0; i <= lastOverflow; i++ {
		(*igen)[i] = 0
	}

	return igen
}

func (igen *InsertKeyGenerator) GetKey(insert Insert) string {
	var buff bytes.Buffer

	var needDot = false
	for dimindx, dist := range *igen {
		// check dist to not write '.' by themselves for zero'd
		// dimensions
		if dist > 0 && needDot {
			buff.WriteByte(byte('.'))
			needDot = false
		}

		for i := 0; i < int(dist); i++ {
			buff.WriteString(insert.Dimensions[dimindx][i])
			needDot = true

			if i < int(dist)-1 {
				buff.WriteByte(byte('/'))
			}
		}
	}

	return buff.String()
}

// computes the number of keys referenced by a query
func (s *Schema) QueryKeyCount(q *Query) uint {
	return 0
}

// creates a generator that iteratively creates all keys referenced
// in the execution of a query. Keys are generated depth-first on each
// dimension, starting from the dimension with the highest index
func (s *Schema) QueryKeyGenerator(q *Query) []string {
	return []string{}
}
