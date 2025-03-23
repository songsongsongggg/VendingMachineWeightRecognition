package exception

// ExceptionEnum 异常类型枚举
type ExceptionEnum int

const (
	SensorError ExceptionEnum = iota
	ForeignObjectError
	RecognitionError
)
