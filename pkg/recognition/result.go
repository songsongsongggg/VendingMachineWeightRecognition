package recognition

import "VendingMachineWeightRecognition/pkg/exception"

// RecognitionItem 识别结果项
type RecognitionItem struct {
	GoodsID string
	Num     int
}

// RecognitionException 识别异常
type RecognitionException struct {
	Layer       int
	Exception   exception.ExceptionEnum
	BeginWeight int
	EndWeight   int
}

// RecognitionResult 识别结果
type RecognitionResult struct {
	Successful bool
	Items      []RecognitionItem
	Exceptions []RecognitionException
}
