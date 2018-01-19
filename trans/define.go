package trans

// StructNode 结构节点
type StructNode struct {
	Name  string
	Type  string
	Index int
	Desc  string
}

// StructInfo 结构信息
type StructInfo struct {
	Name  string
	Desc  string
	Nodes []StructNode
}

// StructInfos 信息列表
type StructInfos []StructInfo
