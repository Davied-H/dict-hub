package mdict

import (
	"os"
	"path/filepath"
	"testing"
)

// Test dictionary file path - adjust this to point to your test dictionary
var testDictPath = ""

func init() {
	// Try to find test dictionaries
	paths := []string{
		"../../dicts/现代英汉汉英综合大辞典.mdx",
		"../../dicts/[英-英] Cambridge Phrasal Verbs Dictionary.mdx",
	}
	
	for _, p := range paths {
		absPath, _ := filepath.Abs(p)
		if _, err := os.Stat(absPath); err == nil {
			testDictPath = absPath
			break
		}
	}
}

func TestNew(t *testing.T) {
	if testDictPath == "" {
		t.Skip("No test dictionary found")
	}
	
	mdict, err := New(testDictPath)
	if err != nil {
		t.Fatalf("Failed to create Mdict: %v", err)
	}
	
	t.Logf("Dictionary loaded: %s", mdict.Name())
	t.Logf("Title: %s", mdict.Title())
	t.Logf("Description: %s", mdict.Description())
	t.Logf("Version: %.1f", mdict.Header.Version)
	t.Logf("Encoding: %d", mdict.Header.Encoding)
	t.Logf("EncryptType: %d", mdict.Header.EncryptType)
	t.Logf("Expected entries: %d", mdict.KeyBlockMeta.EntriesNum)
	t.Logf("Key blocks: %d", mdict.KeyBlockMeta.KeyBlockNum)
}

func TestBuildIndex(t *testing.T) {
	if testDictPath == "" {
		t.Skip("No test dictionary found")
	}
	
	mdict, err := New(testDictPath)
	if err != nil {
		t.Fatalf("Failed to create Mdict: %v", err)
	}
	
	err = mdict.BuildIndex()
	if err != nil {
		t.Fatalf("Failed to build index: %v", err)
	}
	
	t.Logf("Index built successfully")
	t.Logf("Key entries: %d", len(mdict.KeyEntries))
	t.Logf("Record blocks: %d", len(mdict.RecordBlockInfos))
	
	// Show first few entries
	if len(mdict.KeyEntries) > 5 {
		t.Log("First 5 entries:")
		for i := 0; i < 5; i++ {
			t.Logf("  %d: %s", i, mdict.KeyEntries[i].Keyword)
		}
	}
}

func TestLookup(t *testing.T) {
	if testDictPath == "" {
		t.Skip("No test dictionary found")
	}
	
	mdict, err := New(testDictPath)
	if err != nil {
		t.Fatalf("Failed to create Mdict: %v", err)
	}
	
	err = mdict.BuildIndex()
	if err != nil {
		t.Fatalf("Failed to build index: %v", err)
	}
	
	// Test lookup with common words
	testWords := []string{"hello", "world", "test", "dictionary", "a", "the"}
	
	for _, word := range testWords {
		definition, err := mdict.Lookup(word)
		if err != nil {
			t.Logf("Word '%s' not found: %v", word, err)
			continue
		}
		
		// Show first 200 chars of definition
		defStr := string(definition)
		if len(defStr) > 200 {
			defStr = defStr[:200] + "..."
		}
		t.Logf("Found '%s': %s", word, defStr)
	}
}

func TestSuggest(t *testing.T) {
	if testDictPath == "" {
		t.Skip("No test dictionary found")
	}
	
	mdict, err := New(testDictPath)
	if err != nil {
		t.Fatalf("Failed to create Mdict: %v", err)
	}
	
	err = mdict.BuildIndex()
	if err != nil {
		t.Fatalf("Failed to build index: %v", err)
	}
	
	// Test suggestions
	prefixes := []string{"hel", "wor", "dic"}
	
	for _, prefix := range prefixes {
		suggestions := mdict.Suggest(prefix, 5)
		t.Logf("Suggestions for '%s': %v", prefix, suggestions)
	}
}

// Benchmark for index building
func BenchmarkBuildIndex(b *testing.B) {
	if testDictPath == "" {
		b.Skip("No test dictionary found")
	}
	
	for i := 0; i < b.N; i++ {
		mdict, err := New(testDictPath)
		if err != nil {
			b.Fatalf("Failed to create Mdict: %v", err)
		}
		
		err = mdict.BuildIndex()
		if err != nil {
			b.Fatalf("Failed to build index: %v", err)
		}
	}
}

// Benchmark for lookup
func BenchmarkLookup(b *testing.B) {
	if testDictPath == "" {
		b.Skip("No test dictionary found")
	}
	
	mdict, err := New(testDictPath)
	if err != nil {
		b.Fatalf("Failed to create Mdict: %v", err)
	}
	
	err = mdict.BuildIndex()
	if err != nil {
		b.Fatalf("Failed to build index: %v", err)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, _ = mdict.Lookup("hello")
	}
}
