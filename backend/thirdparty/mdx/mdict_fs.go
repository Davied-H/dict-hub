package mdx

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"time"
)

// ErrWordNotFound is returned when a word is not found in the dictionary.
var ErrWordNotFound = errors.New("word not found")

// MdictFS wraps an Mdict instance to implement the io/fs.FS interface.
// This allows an MDX/MDD file to be accessed like a regular file system, for example, for an HTTP file server.
type MdictFS struct {
	mdict *Mdict // The Mdict instance provides access to the dictionary data.
}

// NewMdictFS creates a new MdictFS instance.
func NewMdictFS(mdict *Mdict) *MdictFS {
	if mdict == nil {
		panic("MdictFS: Mdict instance cannot be nil")
	}
	return &MdictFS{
		mdict: mdict,
	}
}

// Open opens a file (a keyword or an MDD resource).
func (mfs *MdictFS) Open(name string) (fs.File, error) {
	log.Debugf("MdictFS: Open called with name: '%s'", name)

	if name == "." || name == "" || strings.HasSuffix(name, "/") {
		name = "."
	}

	modTime := time.Now()
	if mfs.mdict.meta != nil && mfs.mdict.meta.creationDate != "" {
		parsedTime, err := time.Parse("2006-01-02", mfs.mdict.meta.creationDate)
		if err != nil {
			parsedTime, err = time.Parse("2006.01.02 15:04:05", mfs.mdict.meta.creationDate)
			if err != nil {
				log.Warningf("MdictFS: Could not parse CreationDate '%s' for ModTime, using current time.", mfs.mdict.meta.creationDate)
			} else {
				modTime = parsedTime
			}
		} else {
			modTime = parsedTime
		}
	}

	if name == "." {
		log.Debugf("MdictFS: Opening root directory '.'")
		rootInfo := &MdictFileInfo{
			name:    ".",
			isDir:   true,
			modTime: modTime,
		}
		return &MdictFile{
			fs:       mfs,
			name:     ".",
			isDir:    true,
			fileInfo: rootInfo,
		}, nil
	}

	var fileContent []byte
	var lookupErr error

	if mfs.mdict.IsMDD() {
		log.Debugf("MdictFS: MDD file, attempting to find resource: '%s'", name)
		actualName := strings.ReplaceAll(name, "/", "\\")
		if !strings.HasPrefix(actualName, "\\") {
			actualName = "\\" + actualName
		}

		var foundEntry *MDictKeywordEntry
		entries, _ := mfs.mdict.GetKeyWordEntries()
		for _, entry := range entries {
			if strings.EqualFold(entry.KeyWord, actualName) {
				foundEntry = entry
				break
			}
		}

		if foundEntry != nil {
			log.Debugf("MdictFS: Found MDD entry for '%s' (keyword '%s')", name, foundEntry.KeyWord)
			fileContent, lookupErr = mfs.mdict.LocateByKeywordEntry(foundEntry)
		} else {
			log.Debugf("MdictFS: MDD resource '%s' (normalized: '%s') not found in keyword entries.", name, actualName)
			lookupErr = fs.ErrNotExist
		}
	} else { // MDX file
		log.Debugf("MdictFS: MDX file, looking up keyword: '%s'", name)
		definition, err := mfs.mdict.Lookup(name)
		if err != nil {
			if errors.Is(err, ErrWordNotFound) || strings.Contains(err.Error(), "not found") {
				log.Debugf("MdictFS: Keyword '%s' not found in MDX.", name)
				return nil, fs.ErrNotExist
			}
			log.Errorf("MdictFS: Error looking up keyword '%s' in MDX: %v", name, err)
			return nil, fmt.Errorf("error looking up keyword '%s': %w", name, err)
		}
		if len(definition) == 0 {
			log.Debugf("MdictFS: Keyword '%s' found but has no definition.", name)
			return nil, fs.ErrNotExist
		}
		fileContent = definition
		lookupErr = nil
		log.Debugf("MdictFS: Found MDX keyword '%s', content length: %d", name, len(fileContent))
	}

	if lookupErr != nil {
		if errors.Is(lookupErr, fs.ErrNotExist) {
			return nil, fs.ErrNotExist
		}
		log.Errorf("MdictFS: Error getting content for '%s': %v", name, lookupErr)
		return nil, fmt.Errorf("error getting content for '%s': %w", name, lookupErr)
	}

	if fileContent == nil {
		log.Warningf("MdictFS: Content for '%s' is nil after successful lookup, treating as not found.", name)
		return nil, fs.ErrNotExist
	}

	fileInfo := &MdictFileInfo{
		name:    path.Base(name),
		size:    int64(len(fileContent)),
		isDir:   false,
		modTime: modTime,
	}

	return &MdictFile{
		fs:       mfs,
		name:     name,
		isDir:    false,
		content:  fileContent,
		reader:   bytes.NewReader(fileContent),
		fileInfo: fileInfo,
	}, nil
}

// MdictFile implements the fs.File interface.
type MdictFile struct {
	fs       *MdictFS
	name     string
	isDir    bool
	reader   *bytes.Reader
	content  []byte
	fileInfo fs.FileInfo
}

// Stat returns the FileInfo for the file.
func (mf *MdictFile) Stat() (fs.FileInfo, error) {
	if mf.fileInfo == nil {
		log.Warningf("MdictFile.Stat: fileInfo is nil for '%s', creating default.", mf.name)
		modTime := time.Now()
		if mf.fs.mdict.meta != nil && mf.fs.mdict.meta.creationDate != "" {
			parsedTime, err := time.Parse("2006-01-02", mf.fs.mdict.meta.creationDate)
			if err == nil {
				modTime = parsedTime
			}
		}
		mf.fileInfo = &MdictFileInfo{
			name:    path.Base(mf.name),
			size:    int64(len(mf.content)),
			isDir:   mf.isDir,
			modTime: modTime,
		}
	}
	return mf.fileInfo, nil
}

// Read reads up to len(b) bytes from the file.
func (mf *MdictFile) Read(b []byte) (int, error) {
	if mf.isDir {
		log.Debugf("MdictFile.Read: Attempt to read directory '%s'", mf.name)
		return 0, &fs.PathError{Op: "read", Path: mf.name, Err: errors.New("is a directory")}
	}
	if mf.reader == nil {
		log.Warningf("MdictFile.Read: No reader available for file '%s' (might be closed or not a regular file).", mf.name)
		return 0, &fs.PathError{Op: "read", Path: mf.name, Err: fs.ErrClosed}
	}
	return mf.reader.Read(b)
}

// Close closes the file.
func (mf *MdictFile) Close() error {
	log.Debugf("MdictFile.Close: Closing file '%s'", mf.name)
	mf.reader = nil
	mf.content = nil
	mf.fileInfo = nil
	return nil
}

// Seek sets the offset for the next Read or Write on the file.
func (mf *MdictFile) Seek(offset int64, whence int) (int64, error) {
	if mf.isDir {
		return 0, &fs.PathError{Op: "seek", Path: mf.name, Err: errors.New("is a directory")}
	}
	if mf.reader == nil {
		return 0, &fs.PathError{Op: "seek", Path: mf.name, Err: fs.ErrClosed}
	}
	return mf.reader.Seek(offset, whence)
}

// MdictFileInfo implements the fs.FileInfo interface.
type MdictFileInfo struct {
	name    string
	size    int64
	isDir   bool
	modTime time.Time
}

// Name returns the base name of the file.
func (mfi *MdictFileInfo) Name() string { return mfi.name }

// Size returns the length in bytes for regular files.
func (mfi *MdictFileInfo) Size() int64 { return mfi.size }

// IsDir reports whether mfi describes a directory.
func (mfi *MdictFileInfo) IsDir() bool { return mfi.isDir }

// ModTime returns the modification time.
func (mfi *MdictFileInfo) ModTime() time.Time { return mfi.modTime }

// Sys returns underlying data source (can be nil).
func (mfi *MdictFileInfo) Sys() interface{} { return nil }

// Info returns the FileInfo for the file.
func (mfi *MdictFileInfo) Info() (fs.FileInfo, error) { return mfi, nil }

// Type returns the file's type.
func (mfi *MdictFileInfo) Type() fs.FileMode { return mfi.Mode().Type() }

// Mode returns the file mode bits.
func (mfi *MdictFileInfo) Mode() fs.FileMode {
	if mfi.isDir {
		return fs.ModeDir | 0555
	}
	return 0444
}

var _ fs.File = (*MdictFile)(nil)
var _ fs.ReadDirFile = (*MdictFile)(nil)
var _ fs.FS = (*MdictFS)(nil)

// ReadDir reads the contents of the directory.
func (mf *MdictFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if !mf.isDir || mf.name != "." {
		log.Warningf("ReadDir called on non-root or non-directory MdictFile: %s", mf.name)
		return nil, &fs.PathError{Op: "readdir", Path: mf.name, Err: errors.New("not a directory or not root")}
	}

	log.Debugf("ReadDir called for MdictFS root: %s", mf.fs.mdict.filePath)

	keywords, err := mf.fs.mdict.GetKeyWordEntries()
	if err != nil {
		log.Errorf("ReadDir: Error getting keyword entries for %s: %v", mf.fs.mdict.filePath, err)
		return nil, fmt.Errorf("could not get keyword entries: %w", err)
	}

	modTime := time.Now()
	if mf.fs.mdict.meta != nil && mf.fs.mdict.meta.creationDate != "" {
		parsedTime, ptErr := time.Parse("2006-01-02", mf.fs.mdict.meta.creationDate)
		if ptErr == nil {
			modTime = parsedTime
		} else {
			parsedTime, ptErr = time.Parse("2006.01.02 15:04:05", mf.fs.mdict.meta.creationDate)
			if ptErr == nil {
				modTime = parsedTime
			}
		}
	}

	entries := make([]fs.DirEntry, 0, len(keywords))
	for _, kw := range keywords {
		entryName := kw.KeyWord
		isDir := false

		if mf.fs.mdict.IsMDD() {
			entryName = strings.TrimLeft(kw.KeyWord, "\\/")
		}

		dirEntryInfo := &MdictFileInfo{
			name:    path.Base(entryName),
			size:    0,
			isDir:   isDir,
			modTime: modTime,
		}
		entries = append(entries, dirEntryInfo)
	}

	if n > 0 && n < len(entries) {
		entries = entries[:n]
	}

	log.Debugf("ReadDir for '%s' returning %d entries", mf.name, len(entries))
	return entries, nil
}
