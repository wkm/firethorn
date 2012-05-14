package main

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

type Query struct {
}

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

// computes the number of keys referenced by a query
func (s *Schema) QueryKeyCount(q *Query) uint {

}

// creates a generator that iteratively creates all keys referenced
// in the execution of a query. Keys are generated depth-first on each
// dimension, starting from the dimension with the highest index
func (s *Schema) QueryKeyIterator(q *Query) []string {

}
