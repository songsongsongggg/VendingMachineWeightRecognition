package recognition

import (
	"VendingMachineWeightRecognition/pkg/exception"
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
	if len(result.Items) != 2 {
		t.Errorf("应该识别出2个商品，实际识别出%d个", len(result.Items))
	}
}

func TestWeightRecognizer_BasicRecognition(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 1000},
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 900},
	}

	recognizer := NewWeightRecognizer(goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("基本识别应该成功")
	}
	if len(result.Items) != 1 {
		t.Errorf("应该识别出1个商品，实际识别出%d个", len(result.Items))
	}
	if result.Items[0].Num != 1 {
		t.Errorf("商品数量应该是1，实际是%d", result.Items[0].Num)
	}
}

func TestWeightRecognizer_EmptyLayer(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 1000},
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 1000}, // 重量无变化
	}

	recognizer := NewWeightRecognizer(goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("空层识别应该成功")
	}
	if len(result.Items) != 0 {
		t.Errorf("不应该识别出商品，实际识别出%d个", len(result.Items))
	}
}

func TestWeightRecognizer_SensorError(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: -1}, // 传感器异常
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 900},
	}

	recognizer := NewWeightRecognizer(goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if len(result.Exceptions) != 1 {
		t.Errorf("应该检测到1个异常，实际检测到%d个", len(result.Exceptions))
	}
	if result.Exceptions[0].Exception != exception.SensorError {
		t.Error("异常类型应该是传感器异常")
	}
}

func TestWeightRecognizer_ForeignObjectError(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 1000},
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 1100}, // 重量增加，异物异常
	}

	recognizer := NewWeightRecognizer(goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if len(result.Exceptions) != 1 {
		t.Errorf("应该检测到1个异常，实际检测到%d个", len(result.Exceptions))
	}
	if result.Exceptions[0].Exception != exception.ForeignObjectError {
		t.Error("异常类型应该是异物异常")
	}
}

func TestWeightRecognizer_RecognitionError(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 1000},
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 950}, // 无法识别的重量差
	}

	recognizer := NewWeightRecognizer(goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if len(result.Exceptions) != 1 {
		t.Errorf("应该检测到1个异常，实际检测到%d个", len(result.Exceptions))
	}
	if result.Exceptions[0].Exception != exception.RecognitionError {
		t.Error("异常类型应该是无法识别异常")
	}
}
