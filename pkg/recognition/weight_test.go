package recognition

import (
	"VendingMachineWeightRecognition/pkg/model"
	"testing"
)

func TestWeightRecognizer_Recognize(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
		{ID: "000002", Weight: 200},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
		{GoodsID: "000002", Layer: 2, Num: 5},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 1000},
		{Index: 2, Weight: 2000},
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 900},
		{Index: 2, Weight: 1800},
	}

	recognizer := NewWeightRecognizer(goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("识别应该成功")
	}
}
