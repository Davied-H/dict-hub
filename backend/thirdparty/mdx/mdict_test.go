package mdx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMdict_Lookup(t *testing.T) {
	mdict, err := New("testdata/现代汉语八百词.mdx")
	if err != nil {
		t.Fatal(err)
	}
	err = mdict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}
	word := "一律"
	definition, err := mdict.Lookup(word)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, definition, "The definition for '%s' should not be empty", word)
	// The exact content depends on the file, so we do a basic check.
	assert.True(t, strings.Contains(string(definition), "一律"), "The definition for '%s' should contain the word itself", word)

	t.Logf("Lookup result for '%s': %s", word, string(definition))

	// Test a non-existent word
	_, err = mdict.Lookup("一个不存在的词")
	assert.Error(t, err, "Looking up a non-existent word should return an error")
}

func BenchmarkMdict_Lookup(b *testing.B) {
	mdict, err := New("testdata/现代汉语八百词.mdx")
	if err != nil {
		b.Fatal(err)
	}
	err = mdict.BuildIndex()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = mdict.Lookup("一律")
	}
}
