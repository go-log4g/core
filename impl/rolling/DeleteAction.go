package rolling

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/optional"
)

type DeleteAction struct {
	dir          string
	pattern      string
	maxAge       time.Duration
	maxFiles     int
	maxTotalSize int64
}

type deleteFile struct {
	path    string
	modTime time.Time
	size    int64
}

func NewDeleteAction(filePattern string, maxAge time.Duration, maxFiles int, maxTotalSize int64) *DeleteAction {
	dir, pattern := FilePatternMatcher(filePattern)

	return &DeleteAction{
		dir:          dir,
		pattern:      pattern,
		maxAge:       maxAge,
		maxFiles:     maxFiles,
		maxTotalSize: maxTotalSize,
	}
}

func (this *DeleteAction) Execute() {
	files := this.files()

	if this.maxAge > 0 {
		files = this.deleteExpired(files)
	}
	if this.maxFiles > 0 {
		files = this.deleteExcessFiles(files)
	}
	if this.maxTotalSize > 0 {
		this.deleteExcessSize(files)
	}
}

func (this *DeleteAction) files() []deleteFile {
	entries, e := os.ReadDir(this.dir)
	lang.Assert(e == nil, "failed to read rollover directory %q", this.dir)

	result := make([]deleteFile, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		match, e := filepath.Match(this.pattern, entry.Name())
		lang.Assert(e == nil, "invalid rollover file pattern %q", this.pattern)
		if !match {
			continue
		}

		info := optional.OfCommaErr(entry.Info()).OrElsePanic("failed to stat rollover file %q", entry.Name())
		result = append(result, deleteFile{
			path:    filepath.Join(this.dir, entry.Name()),
			modTime: info.ModTime(),
			size:    info.Size(),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].modTime.Before(result[j].modTime)
	})

	return result
}

func (this *DeleteAction) deleteExpired(files []deleteFile) []deleteFile {
	threshold := time.Now().Add(-this.maxAge)
	index := 0

	for index < len(files) && files[index].modTime.Before(threshold) {
		this.delete(files[index])
		index++
	}

	return files[index:]
}

func (this *DeleteAction) deleteExcessFiles(files []deleteFile) []deleteFile {
	excess := len(files) - this.maxFiles
	if excess <= 0 {
		return files
	}

	for i := 0; i < excess; i++ {
		this.delete(files[i])
	}

	return files[excess:]
}

func (this *DeleteAction) deleteExcessSize(files []deleteFile) {
	var total int64
	for _, file := range files {
		total += file.size
	}

	for _, file := range files {
		if total <= this.maxTotalSize {
			return
		}

		this.delete(file)
		total -= file.size
	}
}

func (this *DeleteAction) delete(file deleteFile) {
	lang.Assert(os.Remove(file.path) == nil, "failed to delete rollover file %q", file.path)
}
