package recognition

import (
	"VendingMachineWeightRecognition/pkg/exception"
	"VendingMachineWeightRecognition/pkg/model"
	"testing"
)

func TestWeightRecognizer_Recognize(t *testing.T) {
	// 创建测试数据
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
		{Index: 1, Weight: 900},  // 拿走1个商品1
		{Index: 2, Weight: 1800}, // 拿走1个商品2
	}

	// 创建识别器
	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)

	// 执行识别
	result := recognizer.Recognize(beginLayers, endLayers)

	// 验证结果
	if !result.Successful {
		t.Error("识别应该成功")
	}

	if len(result.Items) != 2 {
		t.Errorf("应该识别出2个商品，实际识别出%d个", len(result.Items))
	}

	// 验证商品1
	found := false
	for _, item := range result.Items {
		if item.GoodsID == "000001" && item.Num == 1 {
			found = true
			break
		}
	}
	if !found {
		t.Error("未正确识别商品1")
	}

	// 验证商品2
	found = false
	for _, item := range result.Items {
		if item.GoodsID == "000002" && item.Num == 1 {
			found = true
			break
		}
	}
	if !found {
		t.Error("未正确识别商品2")
	}
}

func TestWeightRecognizer_RecognizeWithSensorError(t *testing.T) {
	// 创建测试数据
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

	// 创建识别器
	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)

	// 执行识别
	result := recognizer.Recognize(beginLayers, endLayers)

	// 验证结果
	if len(result.Exceptions) != 1 {
		t.Errorf("应该检测到1个异常，实际检测到%d个", len(result.Exceptions))
	}

	if result.Exceptions[0].Exception != exception.SensorError {
		t.Error("异常类型应该是传感器异常")
	}
}

func TestWeightRecognizer_RecognizeWithForeignObject(t *testing.T) {
	// 创建测试数据
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
		{Index: 1, Weight: 1100}, // 重量增加，可能是放置了异物
	}

	// 创建识别器
	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)

	// 执行识别
	result := recognizer.Recognize(beginLayers, endLayers)

	// 验证结果
	if len(result.Exceptions) != 1 {
		t.Errorf("应该检测到1个异常，实际检测到%d个", len(result.Exceptions))
	}

	if result.Exceptions[0].Exception != exception.ForeignObjectError {
		t.Error("异常类型应该是异物异常")
	}
}

// TestWeightRecognizer_BasicRecognition 测试基本识别功能
func TestWeightRecognizer_BasicRecognition(t *testing.T) {
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
		{Index: 1, Weight: 900},  // 拿走1个商品1
		{Index: 2, Weight: 1800}, // 拿走1个商品2
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("基本识别应该成功")
	}
	if len(result.Items) != 2 {
		t.Errorf("应该识别出2个商品，实际识别出%d个", len(result.Items))
	}
}

// TestWeightRecognizer_EmptyLayer 测试空层架
func TestWeightRecognizer_EmptyLayer(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 1000},
		{Index: 2, Weight: 0}, // 空层架
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 900},
		{Index: 2, Weight: 0},
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("空层架识别应该成功")
	}
	if len(result.Items) != 1 {
		t.Errorf("应该只识别出1个商品，实际识别出%d个", len(result.Items))
	}
}

// TestWeightRecognizer_MultipleItems 测试多件商品
func TestWeightRecognizer_MultipleItems(t *testing.T) {
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
		{Index: 1, Weight: 700}, // 拿走3个商品
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("多件商品识别应该成功")
	}
	if len(result.Items) != 1 {
		t.Errorf("应该识别出1个商品，实际识别出%d个", len(result.Items))
	}
	if result.Items[0].Num != 3 {
		t.Errorf("应该识别出3个商品，实际识别出%d个", result.Items[0].Num)
	}
}

// TestWeightRecognizer_SensorTolerance 测试传感器容差
func TestWeightRecognizer_SensorTolerance(t *testing.T) {
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
		{Index: 1, Weight: 905}, // 考虑传感器容差10g
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("传感器容差测试应该成功")
	}
	if len(result.Items) != 1 {
		t.Errorf("应该识别出1个商品，实际识别出%d个", len(result.Items))
	}
}

// TestWeightRecognizer_PackageTolerance 测试包装容差
func TestWeightRecognizer_PackageTolerance(t *testing.T) {
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
		{Index: 1, Weight: 905}, // 考虑包装容差5%
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("包装容差测试应该成功")
	}
	if len(result.Items) != 1 {
		t.Errorf("应该识别出1个商品，实际识别出%d个", len(result.Items))
	}
}

// TestWeightRecognizer_MaxWeight 测试最大重量边界
func TestWeightRecognizer_MaxWeight(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 32767}, // 最大重量
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 32667}, // 拿走1个商品
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("最大重量测试应该成功")
	}
	if len(result.Items) != 1 {
		t.Errorf("应该识别出1个商品，实际识别出%d个", len(result.Items))
	}
}

// TestWeightRecognizer_OverflowWeight 测试重量溢出
func TestWeightRecognizer_OverflowWeight(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 32768}, // 超出最大重量
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 32668},
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if len(result.Exceptions) != 1 {
		t.Errorf("应该检测到1个异常，实际检测到%d个", len(result.Exceptions))
	}
	if result.Exceptions[0].Exception != exception.SensorError {
		t.Error("异常类型应该是传感器异常")
	}
}

// TestWeightRecognizer_MultipleLayers 测试多层同时购物
func TestWeightRecognizer_MultipleLayers(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
		{ID: "000002", Weight: 200},
		{ID: "000003", Weight: 300},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
		{GoodsID: "000002", Layer: 2, Num: 5},
		{GoodsID: "000003", Layer: 3, Num: 3},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 1000},
		{Index: 2, Weight: 2000},
		{Index: 3, Weight: 3000},
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 900},  // 拿走1个商品1
		{Index: 2, Weight: 1800}, // 拿走1个商品2
		{Index: 3, Weight: 2700}, // 拿走1个商品3
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("多层购物识别应该成功")
	}
	if len(result.Items) != 3 {
		t.Errorf("应该识别出3个商品，实际识别出%d个", len(result.Items))
	}
}

// TestWeightRecognizer_NoChange 测试无变化情况
func TestWeightRecognizer_NoChange(t *testing.T) {
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

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("无变化情况应该成功")
	}
	if len(result.Items) != 0 {
		t.Errorf("不应该识别出商品，实际识别出%d个", len(result.Items))
	}
}

// TestWeightRecognizer_MultipleGoodsInLayer 测试同层多商品
func TestWeightRecognizer_MultipleGoodsInLayer(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100}, // 重量差异大
		{ID: "000002", Weight: 500},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
		{GoodsID: "000002", Layer: 1, Num: 5},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 2000}, // 1000g + 1000g
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 1400}, // 拿走1个商品1和1个商品2
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("同层多商品识别应该成功")
	}
	if len(result.Items) != 2 {
		t.Errorf("应该识别出2个商品，实际识别出%d个", len(result.Items))
	}
}

// TestWeightRecognizer_AmbiguousRecognition 测试识别结果不唯一
func TestWeightRecognizer_AmbiguousRecognition(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
		{ID: "000002", Weight: 100}, // 相同重量
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
		{GoodsID: "000002", Layer: 1, Num: 10},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 2000},
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 1900}, // 拿走1个商品，但无法确定是哪个
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if len(result.Exceptions) != 1 {
		t.Errorf("应该检测到1个异常，实际检测到%d个", len(result.Exceptions))
	}
	if result.Exceptions[0].Exception != exception.RecognitionError {
		t.Error("异常类型应该是无法识别异常")
	}
}

// TestWeightRecognizer_PartialException 测试部分层异常
func TestWeightRecognizer_PartialException(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
		{ID: "000002", Weight: 200},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
		{GoodsID: "000002", Layer: 2, Num: 5},
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: -1},   // 第1层传感器异常
		{Index: 2, Weight: 2000}, // 第2层正常
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 900},
		{Index: 2, Weight: 1800}, // 拿走1个商品2
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if len(result.Exceptions) != 1 {
		t.Errorf("应该检测到1个异常，实际检测到%d个", len(result.Exceptions))
	}
	if result.Exceptions[0].Exception != exception.SensorError {
		t.Error("异常类型应该是传感器异常")
	}
	if len(result.Items) != 1 {
		t.Errorf("应该识别出1个商品，实际识别出%d个", len(result.Items))
	}
	if result.Items[0].GoodsID != "000002" {
		t.Error("应该识别出商品2")
	}
}

// TestWeightRecognizer_DuplicateItems 测试重复商品合并
func TestWeightRecognizer_DuplicateItems(t *testing.T) {
	goods := []model.Goods{
		{ID: "000001", Weight: 100},
	}

	stocks := []model.Stock{
		{GoodsID: "000001", Layer: 1, Num: 10},
		{GoodsID: "000001", Layer: 2, Num: 10}, // 同一商品在不同层
	}

	beginLayers := []model.Layer{
		{Index: 1, Weight: 1000},
		{Index: 2, Weight: 1000},
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 900}, // 拿走1个商品1
		{Index: 2, Weight: 900}, // 拿走1个商品1
	}

	recognizer := NewWeightRecognizer(10, 5.0, goods, stocks)
	result := recognizer.Recognize(beginLayers, endLayers)

	if !result.Successful {
		t.Error("重复商品合并应该成功")
	}
	if len(result.Items) != 1 {
		t.Errorf("应该合并为1个商品项，实际有%d个", len(result.Items))
	}
	if result.Items[0].Num != 2 {
		t.Errorf("应该合并数量为2，实际为%d", result.Items[0].Num)
	}
}
