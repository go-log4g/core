package impl

type FileAppenderDelegate interface {
	BeforeAppend(data []byte)
	Write(data []byte)
	AfterAppend(data []byte)
}
