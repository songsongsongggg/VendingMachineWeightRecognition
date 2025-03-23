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

	// 初始化层商品映射
	for _, stock := range stocks {
		if _, exists := wr.layerStockMap[stock.Layer]; !exists {
			wr.layerStockMap[stock.Layer] = make(map[string]int)
		}
		wr.layerStockMap[stock.Layer][stock.GoodsID] = stock.Num

		// 找到对应的商品
		for _, good := range goods {
			if good.ID == stock.GoodsID {
				wr.layerGoodsMap[stock.Layer] = append(wr.layerGoodsMap[stock.Layer], good)
				break
			}
		}
	}

	return wr
}

// Recognize 识别购物清单
func (wr *WeightRecognizer) Recognize(beginLayers, endLayers []model.Layer) RecognitionResult {
	result := RecognitionResult{
		Successful: true,
		Items:      make([]RecognitionItem, 0),
		Exceptions: make([]RecognitionException, 0),
	}

	// 按层号排序
	for i := 0; i < len(beginLayers); i++ {
		beginLayer := beginLayers[i]
		endLayer := endLayers[i]

		// 计算重量差
		weightDiff := beginLayer.Weight - endLayer.Weight

		// 如果重量差为0，说明没有购物
		if weightDiff == 0 {
			continue
		}

		// 识别该层的商品
		items := wr.recognizeLayer(beginLayer.Index, weightDiff)
		if len(items) > 0 {
			result.Items = append(result.Items, items...)
		}
	}

	return result
}

// recognizeLayer 识别单层的商品
func (wr *WeightRecognizer) recognizeLayer(layer int, weightDiff int) []RecognitionItem {
	items := make([]RecognitionItem, 0)
	layerGoods := wr.layerGoodsMap[layer]

	if len(layerGoods) == 0 {
		return items
	}

	// 如果只有一种商品，直接计算数量
	if len(layerGoods) == 1 {
		good := layerGoods[0]
		stock := wr.layerStockMap[layer][good.ID]

		// 计算可能的数量
		num := weightDiff / good.Weight

		// 检查数量是否合理
		if num > 0 && num <= stock {
			items = append(items, RecognitionItem{
				GoodsID: good.ID,
				Num:     num,
			})
		}

		return items
	}

	// 处理多商品的情况
	// 按重量从小到大排序
	for i := 0; i < len(layerGoods); i++ {
		for j := i + 1; j < len(layerGoods); j++ {
			if layerGoods[i].Weight > layerGoods[j].Weight {
				layerGoods[i], layerGoods[j] = layerGoods[j], layerGoods[i]
			}
		}
	}

	// 尝试识别每个商品
	for _, good := range layerGoods {
		stock := wr.layerStockMap[layer][good.ID]
		num := weightDiff / good.Weight

		if num > 0 && num <= stock {
			items = append(items, RecognitionItem{
				GoodsID: good.ID,
				Num:     num,
			})
			break // 找到第一个匹配的商品就返回
		}
	}

	return items
}
