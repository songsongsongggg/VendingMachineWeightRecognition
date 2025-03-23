package main

import (
	"VendingMachineWeightRecognition/pkg/model"
	"VendingMachineWeightRecognition/pkg/recognition"
	"fmt"
	"log"
)

func main() {
	log.Println("程序启动...")

	// 初始化商品数据
	goods := []model.Goods{
		{ID: "1", Weight: 500},
		{ID: "2", Weight: 500},
		{ID: "3", Weight: 550},
	}

	// 初始化库存数据
	stocks := []model.Stock{
		{Layer: 1, GoodsID: "1", Num: 5},
		{Layer: 1, GoodsID: "2", Num: 5},
		{Layer: 2, GoodsID: "3", Num: 5},
	}

	// 创建重量识别器
	recognizer := recognition.NewWeightRecognizer(
		10,   // 传感器容差
		0.05, // 包装容差
		goods,
		stocks,
	)

	// 模拟层重量变化
	beginLayers := []model.Layer{
		{Index: 1, Weight: 5000},
		{Index: 2, Weight: 5500},
	}

	endLayers := []model.Layer{
		{Index: 1, Weight: 4500},
		{Index: 2, Weight: 4950},
	}

	// 识别结果
	result := recognizer.Recognize(beginLayers, endLayers)

	// 输出识别结果
	fmt.Println("识别结果:")
	for _, item := range result.Items {
		fmt.Printf("商品ID: %s, 数量: %d\n", item.GoodsID, item.Num)
	}

	log.Println("程序运行完成")
}
