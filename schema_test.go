package firethorn

import (
	"github.com/bmizerany/assert"
	"testing"
	f "tumblr/firethorn"
)

func TestInsertKeyGenerator(t *testing.T) {
	ins := f.Insert{"foo", 1,
		[][]string{
			[]string{"a", "b", "c"},
			[]string{"1", "23"},
			[]string{"p"},
		},
	}

	// 1) validate NextKey
	gen := ins.KeyGenerator()
	assert.Equal(t, 3, len(gen))

	assert.Equal(t, f.InsertKeyGenerator{0, 0, 0}, gen)
	assert.Equal(t, f.InsertKeyGenerator{1, 0, 0}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{2, 0, 0}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{3, 0, 0}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{0, 1, 0}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{1, 1, 0}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{2, 1, 0}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{3, 1, 0}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{0, 2, 0}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{1, 2, 0}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{2, 2, 0}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{3, 2, 0}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{0, 0, 1}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{1, 0, 1}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{2, 0, 1}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{3, 0, 1}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{0, 1, 1}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{1, 1, 1}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{2, 1, 1}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{3, 1, 1}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{0, 2, 1}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{1, 2, 1}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{2, 2, 1}, *gen.NextKey(ins))
	assert.Equal(t, f.InsertKeyGenerator{3, 2, 1}, *gen.NextKey(ins))

	// 2) validate GetKey
	assert.Equal(t, "", (&f.InsertKeyGenerator{0, 0, 0}).GetKey(ins))
	assert.Equal(t, "p", (&f.InsertKeyGenerator{0, 0, 1}).GetKey(ins))
	assert.Equal(t, "a.1/23.p", (&f.InsertKeyGenerator{1, 2, 1}).GetKey(ins))
	assert.Equal(t, "a/b/c.p", (&f.InsertKeyGenerator{3, 0, 1}).GetKey(ins))
	assert.Equal(t, "a/b/c", (&f.InsertKeyGenerator{3, 0, 0}).GetKey(ins))
}
