package recognition

import (
	"VendingMachineWeightRecognition/pkg/exception"
	"VendingMachineWeightRecognition/pkg/model"
)

// WeightRecognizer 重量识别器
type WeightRecognizer struct {
	sensorTolerance int // 传感器容差
	goods           []model.Goods
	stocks          []model.Stock
	layerGoodsMap   map[int][]model.Goods
	layerStockMap   map[int]map[string]int
}

// NewWeightRecognizer 创建新的重量识别器
func NewWeightRecognizer(sensorTolerance int, goods []model.Goods, stocks []model.Stock) *WeightRecognizer {
	wr := &WeightRecognizer{
		sensorTolerance: sensorTolerance,
		goods:           goods,
		stocks:          stocks,
		layerGoodsMap:   make(map[int][]model.Goods),
		layerStockMap:   make(map[int]map[string]int),
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

		// 检查传感器异常
		if beginLayer.Weight < 0 || beginLayer.Weight > 32767 ||
			endLayer.Weight < 0 || endLayer.Weight > 32767 {
			result.Exceptions = append(result.Exceptions, RecognitionException{
				Layer:       beginLayer.Index,
				Exception:   exception.SensorError,
				BeginWeight: beginLayer.Weight,
				EndWeight:   endLayer.Weight,
			})
			continue
		}

		// 检查异物异常
		if endLayer.Weight > beginLayer.Weight {
			result.Exceptions = append(result.Exceptions, RecognitionException{
				Layer:       beginLayer.Index,
				Exception:   exception.ForeignObjectError,
				BeginWeight: beginLayer.Weight,
				EndWeight:   endLayer.Weight,
			})
			continue
		}

		// 计算重量差
		weightDiff := beginLayer.Weight - endLayer.Weight

		// 考虑传感器容差，判断是否无购物
		if weightDiff <= wr.sensorTolerance && weightDiff >= -wr.sensorTolerance {
			continue // 无购物
		}

		// 识别该层的商品
		items := wr.recognizeLayer(beginLayer.Index, weightDiff)
		if len(items) == 0 {
			result.Exceptions = append(result.Exceptions, RecognitionException{
				Layer:       beginLayer.Index,
				Exception:   exception.RecognitionError,
				BeginWeight: beginLayer.Weight,
				EndWeight:   endLayer.Weight,
			})
			continue
		}

		result.Items = append(result.Items, items...)
	}

	return result
}

// recognizeLayer 识别单层的商品
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

		// 计算可能的数量范围
		minNum := (weightDiff - wr.sensorTolerance) / good.Weight
		maxNum := (weightDiff + wr.sensorTolerance) / good.Weight

		// 限制在库存范围内
		if minNum < 0 {
			minNum = 0
		}
		if maxNum > stock {
			maxNum = stock
		}

		// 如果范围合理，取中间值
		if minNum <= maxNum && minNum > 0 {
			num := minNum // 优先选择最小数量
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

	// 检查是否有相同重量的商品
	hasSameWeight := false
	for i := 1; i < len(layerGoods); i++ {
		if layerGoods[i].Weight == layerGoods[i-1].Weight {
			hasSameWeight = true
			break
		}
	}

	// 如果有相同重量的商品，返回识别异常
	if hasSameWeight {
		return nil
	}

	// 尝试所有可能的组合
	bestItems := make([]RecognitionItem, 0)
	bestDiff := weightDiff + wr.sensorTolerance + 1

	// 生成所有可能的组合
	maxCombinations := 1 << len(layerGoods) // 2^n
	for i := 1; i < maxCombinations; i++ {
		currentItems := make([]RecognitionItem, 0)
		totalWeight := 0
		valid := true

		// 检查每个商品是否在当前组合中
		for j := 0; j < len(layerGoods); j++ {
			if (i & (1 << j)) != 0 {
				good := layerGoods[j]
				stock := wr.layerStockMap[layer][good.ID]

				// 检查库存
				if stock <= 0 {
					valid = false
					break
				}

				// 计算可能的数量范围
				minNum := (weightDiff - wr.sensorTolerance) / good.Weight
				maxNum := (weightDiff + wr.sensorTolerance) / good.Weight

				// 限制在库存范围内
				if minNum < 0 {
					minNum = 0
				}
				if maxNum > stock {
					maxNum = stock
				}

				// 如果范围合理，取中间值
				if minNum <= maxNum && minNum > 0 {
					num := minNum // 优先选择最小数量
					currentItems = append(currentItems, RecognitionItem{
						GoodsID: good.ID,
						Num:     num,
					})
					totalWeight += good.Weight * num
				} else {
					valid = false
					break
				}
			}
		}

		if !valid {
			continue
		}

		// 计算差异
		diff := abs(totalWeight - weightDiff)

		// 如果在容差范围内
		if diff <= wr.sensorTolerance {
			// 如果找到更好的组合（商品数量更多或差异更小）
			if len(currentItems) > len(bestItems) || (len(currentItems) == len(bestItems) && diff < bestDiff) {
				bestItems = currentItems
				bestDiff = diff
			}
		} else if len(bestItems) == 0 && diff < bestDiff {
			// 如果还没有找到在容差范围内的组合，记录最接近的
			bestItems = currentItems
			bestDiff = diff
		}
	}

	// 如果最小差异超过容差范围的两倍，返回空
	if bestDiff > wr.sensorTolerance*2 {
		return nil
	}

	return bestItems
}

// abs 返回整数的绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
