package model

// Goods 表示商品信息
type Goods struct {
	ID     string // 6 位的商品编号，每个商品唯一
	Weight int    // 商品单件重量，单位 g
}
