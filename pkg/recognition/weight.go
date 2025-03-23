package recognition

import (
	"VendingMachineWeightRecognition/pkg/model"
)

// WeightRecognizer 重量识别器
type WeightRecognizer struct {
	goods         []model.Goods
	stocks        []model.Stock
	layerGoodsMap map[int][]model.Goods
	layerStockMap map[int]map[string]int
}

// NewWeightRecognizer 创建新的重量识别器
func NewWeightRecognizer(goods []model.Goods, stocks []model.Stock) *WeightRecognizer {
	wr := &WeightRecognizer{
		goods:         goods,
		stocks:        stocks,
		layerGoodsMap: make(map[int][]model.Goods),
		layerStockMap: make(map[int]map[string]int),
	}
	return wr
}

// Recognize 识别购物清单
func (wr *WeightRecognizer) Recognize(beginLayers, endLayers []model.Layer) RecognitionResult {
	return RecognitionResult{
		Successful: true,
		Items:      make([]RecognitionItem, 0),
		Exceptions: make([]RecognitionException, 0),
	}
}
