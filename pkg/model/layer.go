package model

// Layer 表示售货机的一层
type Layer struct {
	Index  int // 编号，从 1 开始
	Weight int // 重量传感器数值，单位 g
}
