//
// Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package mdx

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("default")

// Mdict is a high-level wrapper for mdx/mdd dictionary files.
// It embeds MdictBase to handle the underlying parsing logic and provides a user-facing API.
type Mdict struct {
	*MdictBase
	rangeTreeRoot *RecordBlockRangeTreeNode
}

// New creates a new Mdict instance.
// It automatically determines the dictionary type based on the file extension (.mdx or .mdd).
func New(filename string) (*Mdict, error) {
	dictType := MdictTypeMdx
	if strings.ToLower(filepath.Ext(filename)) == ".mdd" {
		dictType = MdictTypeMdd
	}

	mdict := &Mdict{
		MdictBase: &MdictBase{
			filePath:      filename,
			fileType:      dictType,
			rangeTreeRoot: new(RecordBlockRangeTreeNode),
		},
	}
	return mdict, mdict.init()
}

// init initializes the dictionary, mainly responsible for reading and parsing the file header and key block metadata.
func (mdict *Mdict) init() error {
	// Read dictionary header
	if err := mdict.readDictHeader(); err != nil {
		return err
	}

	// Read key block metadata
	if err := mdict.readKeyBlockMeta(); err != nil {
		return err
	}

	return nil
}

// BuildIndex builds the complete dictionary index.
// This process can consume significant memory and time as it needs to read all keyword and record block information.
// It is recommended to call this once during program initialization.
func (mdict *Mdict) BuildIndex() error {
	if err := mdict.readKeyBlockInfo(); err != nil {
		return err
	}

	if err := mdict.readKeyEntries(); err != nil {
		return err
	}

	if err := mdict.readRecordBlockMeta(); err != nil {
		return err
	}

	if err := mdict.readRecordBlockInfo(); err != nil {
		return err
	}

	mdict.buildRecordRangeTree()

	return nil
}

// Name returns the name of the dictionary, usually the filename without the extension.
func (mdict *Mdict) Name() string {
	_, rawpath := filepath.Split(mdict.filePath)
	rawpath = strings.TrimSuffix(rawpath, ".mdx")
	rawpath = strings.TrimSuffix(rawpath, ".mdd")
	return rawpath
}

// Title returns the title of the dictionary.
func (mdict *Mdict) Title() string {
	return mdict.meta.title
}

// Description returns the description of the dictionary.
func (mdict *Mdict) Description() string {
	return mdict.meta.description
}

// GeneratedByEngineVersion returns the engine version that generated the dictionary.
func (mdict *Mdict) GeneratedByEngineVersion() string {
	return mdict.meta.generatedByEngineVersion
}

// CreationDate returns the creation date of the dictionary.
func (mdict *Mdict) CreationDate() string {
	return mdict.meta.creationDate
}

// Version returns the version number of the dictionary.
func (mdict *Mdict) Version() string {
	return fmt.Sprintf("%f", mdict.meta.version)
}

// IsMDD checks if the dictionary is an MDD file.
func (mdict *Mdict) IsMDD() bool {
	return mdict.fileType == MdictTypeMdd
}

// IsRecordEncrypted checks if the dictionary's record blocks are encrypted.
func (mdict *Mdict) IsRecordEncrypted() bool {
	return mdict.meta.encryptType == EncryptRecordEnc
}

// IsUTF16 checks if the dictionary's encoding is UTF-16.
func (mdict *Mdict) IsUTF16() bool {
	return mdict.meta.encoding == EncodingUtf16
}

// Lookup finds the definition for a given word.
// It uses binary search for efficiency, with a time complexity of O(log n).
func (mdict *Mdict) Lookup(word string) ([]byte, error) {
	word = strings.TrimSpace(word)
	entries := mdict.keyBlockData.keyEntries
	i := sort.Search(len(entries), func(i int) bool {
		// Note: This comparison depends on the keyword sorting rules in the MDict file.
		// For most dictionaries, this should be effective.
		return entries[i].KeyWord >= word
	})

	if i < len(entries) && entries[i].KeyWord == word {
		log.Infof("mdict.Lookup hit entry [%d/%d] key:(%s)", i, len(entries), word)
		return mdict.LocateByKeywordEntry(entries[i])
	}

	return nil, fmt.Errorf("word not found: (%s)", word)
}

// LocateByKeywordEntry locates and returns the definition by keyword entry.
func (mdict *Mdict) LocateByKeywordEntry(entry *MDictKeywordEntry) ([]byte, error) {
	if entry == nil {
		return nil, errors.New("invalid mdict keyword entry")
	}
	return mdict.MdictBase.locateByKeywordEntry(entry)
}

// LocateByKeywordIndex locates and returns the definition by keyword index.
func (mdict *Mdict) LocateByKeywordIndex(index *MDictKeywordIndex) ([]byte, error) {
	if index == nil {
		return nil, errors.New("invalid mdict keyword index")
	}
	return mdict.MdictBase.locateByKeywordIndex(index)
}

// GetKeyWordEntries returns all keyword entries in the dictionary.
func (mdict *Mdict) GetKeyWordEntries() ([]*MDictKeywordEntry, error) {
	return mdict.MdictBase.GetKeyWordEntries()
}

// GetKeyWordEntriesSize returns the total number of keyword entries in the dictionary.
func (mdict *Mdict) GetKeyWordEntriesSize() int64 {
	return mdict.keyBlockData.keyEntriesSize
}

// KeywordEntryToIndex converts a keyword entry to a more detailed keyword index.
func (mdict *Mdict) KeywordEntryToIndex(item *MDictKeywordEntry) (*MDictKeywordIndex, error) {
	return mdict.MdictBase.keywordEntryToIndex(item)
}
