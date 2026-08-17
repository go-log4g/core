package rolling

import (
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/optional"
)

type DefaultRolloverStrategy struct {
	max          int
	deleteAction *DeleteAction
}

func NewDefaultRolloverStrategy(max int, deleteAction *DeleteAction) *DefaultRolloverStrategy {
	lang.Assert(max > 0, "max must be positive")

	return &DefaultRolloverStrategy{
		max:          max,
		deleteAction: deleteAction,
	}
}

// Implements RolloverStrategy
func (this *DefaultRolloverStrategy) Rollover(file, filePattern string, at time.Time) {
	lang.Assert(strings.Contains(filePattern, "%i"), "filePattern must contain %%i")

	for i := 1; i <= this.max; i++ {
		target := FormatFilePattern(filePattern, i, at)

		if _, e := os.Stat(target); os.IsNotExist(e) {
			this.archive(file, target)
			this.cleanup()
			return
		}
	}

	oldest := FormatFilePattern(filePattern, 1, at)
	lang.Assert(os.Remove(oldest) == nil, "failed to remove rollover file %q", oldest)

	for i := 2; i <= this.max; i++ {
		source := FormatFilePattern(filePattern, i, at)
		target := FormatFilePattern(filePattern, i-1, at)
		lang.Assert(os.Rename(source, target) == nil, "failed to rename rollover file %q to %q", source, target)
	}

	newest := FormatFilePattern(filePattern, this.max, at)
	this.archive(file, newest)
	this.cleanup()
}

func (this *DefaultRolloverStrategy) archive(source, target string) {
	dir := filepath.Dir(target)
	if dir != "." {
		lang.Assert(os.MkdirAll(dir, 0755) == nil, "failed to create rollover directory %q", dir)
	}

	switch strings.ToLower(filepath.Ext(target)) {
	case ".gz":
		this.gzip(source, target)
	case ".zip":
		this.zip(source, target)
	default:
		lang.Assert(os.Rename(source, target) == nil, "failed to rename log file %q to %q", source, target)
	}
}

func (this *DefaultRolloverStrategy) gzip(source, target string) {
	input := optional.OfCommaErr(os.Open(source)).OrElsePanic("failed to open rollover source file %q", source)
	defer input.Close()
	output := optional.OfCommaErr(os.Create(target)).OrElsePanic("failed to create compressed rollover file %q", target)
	defer output.Close()

	writer := gzip.NewWriter(output)
	defer writer.Close()

	_, e := io.Copy(writer, input)
	lang.Assert(e == nil, "failed to compress rollover file %q to %q", source, target)

	lang.Assert(input.Close() == nil, "failed to close rollover source file %q", source)
	lang.Assert(writer.Close() == nil, "failed to finalize compressed rollover file %q", target)
	lang.Assert(output.Close() == nil, "failed to close compressed rollover file %q", target)

	lang.Assert(os.Remove(source) == nil, "failed to remove rollover source file %q", source)
}

func (this *DefaultRolloverStrategy) zip(source, target string) {
	input := optional.OfCommaErr(os.Open(source)).OrElsePanic("failed to open rollover source file %q", source)
	defer input.Close()
	output := optional.OfCommaErr(os.Create(target)).OrElsePanic("failed to create compressed rollover file %q", target)
	defer output.Close()

	writer := zip.NewWriter(output)
	defer writer.Close()
	entry := optional.OfCommaErr(writer.Create(strings.TrimSuffix(filepath.Base(target), filepath.Ext(target)))).OrElsePanic("failed to create ZIP entry for rollover file %q", target)

	_, e := io.Copy(entry, input)
	lang.Assert(e == nil, "failed to compress rollover file %q to %q", source, target)

	lang.Assert(input.Close() == nil, "failed to close rollover source file %q", source)
	lang.Assert(writer.Close() == nil, "failed to finalize compressed rollover file %q", target)
	lang.Assert(output.Close() == nil, "failed to close compressed rollover file %q", target)

	lang.Assert(os.Remove(source) == nil, "failed to remove rollover source file %q", source)
}

func (this *DefaultRolloverStrategy) cleanup() {
	if this.deleteAction != nil {
		this.deleteAction.Execute()
	}
}
