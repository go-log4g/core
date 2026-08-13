package impl

import "path/filepath"

type FilePatternConverter struct {
	AbstractPatternConverter
}

func NewFilePatternConverter(formatting FormattingInfo) *FilePatternConverter {
	return &FilePatternConverter{
		AbstractPatternConverter: NewAbstractPatternConverter(formatting),
	}
}

func (this *FilePatternConverter) Append(result []byte, event *LogEvent) []byte {
	return this.AppendFormatted(result, filepath.Base(event.CallerContext.Caller.File))
}
