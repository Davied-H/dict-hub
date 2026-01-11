package mdx

import (
	"os"
	"testing"
)

func TestOALE9(t *testing.T) {
	path := "testdata/mdx/testdict.mdx"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("Skipping test because test data file is missing: %s", path)
	}

	dict, err := New(path)
	if err != nil {
		t.Error(err)
		return
	}

	err = dict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}

	keywordEntries, err := dict.GetKeyWordEntries()
	if err != nil {
		t.Fatal(err)
	}
	for idx, entry := range keywordEntries {
		if idx > 10 {
			break
		}
		t.Logf("\n\n-----------\n\n")
		t.Logf("keyword: %s", entry.KeyWord)
		index, err := dict.KeywordEntryToIndex(entry)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("index: %+v", index)

		def, err := dict.LocateByKeywordIndex(index)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("def: %s", def)
	}
}
