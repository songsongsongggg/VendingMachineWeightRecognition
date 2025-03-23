package model

// Stock 表示库存信息
type Stock struct {
	GoodsID string // 库存对应的商品
	Layer   int    // 库存对应的层架
	Num     int    // 库存数量
}
