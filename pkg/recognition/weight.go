package recognition

import (
	"VendingMachineWeightRecognition/pkg/exception"
	"VendingMachineWeightRecognition/pkg/model"
	"math"
	"sort"
)

// WeightRecognizer 重量识别器
type WeightRecognizer struct {
	sensorTolerance  int     // 传感器容差
	packageTolerance float64 // 包装容差
	goods            []model.Goods
	stocks           []model.Stock
	layerGoodsMap    map[int][]model.Goods  // 层号到商品的映射
	layerStockMap    map[int]map[string]int // 层号到商品库存的映射
}

// NewWeightRecognizer 创建新的重量识别器
func NewWeightRecognizer(sensorTolerance int, packageTolerance float64, goods []model.Goods, stocks []model.Stock) *WeightRecognizer {
	wr := &WeightRecognizer{
		sensorTolerance:  sensorTolerance,
		packageTolerance: packageTolerance,
		goods:            goods,
		stocks:           stocks,
		layerGoodsMap:    make(map[int][]model.Goods),
		layerStockMap:    make(map[int]map[string]int),
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
	sort.Slice(beginLayers, func(i, j int) bool {
		return beginLayers[i].Index < beginLayers[j].Index
	})
	sort.Slice(endLayers, func(i, j int) bool {
		return endLayers[i].Index < endLayers[j].Index
	})

	// 处理每一层
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

		// 合并相同商品
		result.Items = wr.mergeItems(result.Items, items)
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

		// 考虑包装容差
		minWeight := int(float64(good.Weight) * (1 - wr.packageTolerance/100))
		maxWeight := int(float64(good.Weight) * (1 + wr.packageTolerance/100))

		// 计算可能的数量范围
		// 使用浮点数计算以提高精度
		minNumFloat := float64(weightDiff-wr.sensorTolerance) / float64(maxWeight)
		maxNumFloat := float64(weightDiff+wr.sensorTolerance) / float64(minWeight)

		// 向上/向下取整
		minNum := int(math.Ceil(minNumFloat))
		maxNum := int(math.Floor(maxNumFloat))

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
	// 按重量从小到大排序，便于组合
	sort.Slice(layerGoods, func(i, j int) bool {
		return layerGoods[i].Weight < layerGoods[j].Weight
	})

	// 检查是否有相同重量的商品
	hasSameWeight := false
	for i := 1; i < len(layerGoods); i++ {
		if layerGoods[i].Weight == layerGoods[i-1].Weight {
			hasSameWeight = true
			break
		}
	}

	// 如果有相同重量的商品，直接返回识别异常
	if hasSameWeight {
		return nil
	}

	// 尝试所有可能的组合
	bestItems := wr.findBestCombination(layerGoods, layer, weightDiff)
	if len(bestItems) > 0 {
		items = append(items, bestItems...)
	}

	return items
}

// findBestCombination 查找最佳组合
func (wr *WeightRecognizer) findBestCombination(goods []model.Goods, layer int, targetWeight int) []RecognitionItem {
	// 按重量从小到大排序
	sort.Slice(goods, func(i, j int) bool {
		return goods[i].Weight < goods[j].Weight
	})

	// 检查是否有相同重量的商品
	for i := 1; i < len(goods); i++ {
		if goods[i].Weight == goods[i-1].Weight {
			return nil // 有相同重量的商品，返回识别异常
		}
	}

	// 尝试所有可能的组合
	bestItems := make([]RecognitionItem, 0)
	bestDiff := targetWeight + wr.sensorTolerance + 1

	// 生成所有可能的组合
	maxCombinations := 1 << len(goods) // 2^n
	for i := 1; i < maxCombinations; i++ {
		items := make([]RecognitionItem, 0)
		totalWeight := 0
		valid := true

		// 检查每个商品是否在当前组合中
		for j := 0; j < len(goods); j++ {
			if (i & (1 << j)) != 0 {
				good := goods[j]
				stock := wr.layerStockMap[layer][good.ID]

				// 检查库存
				if stock <= 0 {
					valid = false
					break
				}

				// 考虑包装容差
				minWeight := int(float64(good.Weight) * (1 - wr.packageTolerance/100))
				maxWeight := int(float64(good.Weight) * (1 + wr.packageTolerance/100))

				// 使用平均重量
				avgWeight := (minWeight + maxWeight) / 2
				totalWeight += avgWeight

				items = append(items, RecognitionItem{
					GoodsID: good.ID,
					Num:     1,
				})
			}
		}

		if !valid {
			continue
		}

		// 计算差异
		diff := abs(totalWeight - targetWeight)

		// 如果在容差范围内
		if diff <= wr.sensorTolerance {
			// 如果找到更好的组合（商品数量更多或差异更小）
			if len(items) > len(bestItems) || (len(items) == len(bestItems) && diff < bestDiff) {
				bestItems = items
				bestDiff = diff
			}
		} else if len(bestItems) == 0 && diff < bestDiff {
			// 如果还没有找到在容差范围内的组合，记录最接近的
			bestItems = items
			bestDiff = diff
		}
	}

	// 如果最小差异超过容差范围的两倍，返回空
	if bestDiff > wr.sensorTolerance*2 {
		return nil
	}

	return bestItems
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// abs 返回整数的绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// mergeItems 合并相同商品的项
func (wr *WeightRecognizer) mergeItems(items1, items2 []RecognitionItem) []RecognitionItem {
	result := make([]RecognitionItem, 0)
	itemMap := make(map[string]int)

	// 合并所有项
	for _, item := range items1 {
		itemMap[item.GoodsID] += item.Num
	}
	for _, item := range items2 {
		itemMap[item.GoodsID] += item.Num
	}

	// 转换回切片
	for goodsID, num := range itemMap {
		if num > 0 {
			result = append(result, RecognitionItem{
				GoodsID: goodsID,
				Num:     num,
			})
		}
	}

	return result
}
