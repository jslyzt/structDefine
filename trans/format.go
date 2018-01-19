package trans

import "fmt"

//////////////////////////////////////////////////////////////////////////
func (node *StructNode) goNode() string {
	return fmt.Sprintf("\t%v %v `%v %v` // %v \n", toUpper(node.Name), node.Type, node.toBson(), node.toJSON(), node.Desc)
}

func (node *StructNode) toBson() string {
	return fmt.Sprintf("bson:\"%v\"", toBson(node.Name, node.Type))
}

func (node *StructNode) toJSON() string {
	return fmt.Sprintf("json:\"%v\"", toJSON(node.Name))
}

func (node *StructNode) goMemberFunc(form string) string {
	name := toUpper(node.Name)
	return fmt.Sprintf(`
// Set%v %v
func (form *%v) Set%v(val %v) {
	form.%v = val
	form.MarkDirty(%v)
}`, name, node.Desc, form, name, node.Type, name, node.Index)
}

//////////////////////////////////////////////////////////////////////////

func (form *StructInfo) goStruct() string {
	return fmt.Sprintf(`
// %v %v
type %v struct {
%v
	%v
}`, form.Name, form.Desc, form.Name, form.goNodes(), "BTMark mark.BitMark `bson:\"-\" json:\"-\"` // 标记")
}

func (form *StructInfo) goNodes() string {
	ostr := ""
	for _, node := range form.Nodes {
		ostr = ostr + node.goNode()
	}
	return ostr
}

func (form *StructInfo) goFuncs() string {
	return fmt.Sprintf(`
// member opt funcs
%v

// find funcs
%v

// update funcs
%v

// mark funcs
%v`, form.goMemberFuncs(), form.goFindFuncs(), form.goUpdateFuncs(), form.goMarkFuncs())
}

func (form *StructInfo) goMemberFuncs() string {
	ostr := ""
	for _, node := range form.Nodes {
		ostr = ostr + node.goMemberFunc(form.Name)
	}
	return ostr
}

func (form *StructInfo) goFindFuncs() string {
	return fmt.Sprintf(`
// Find mongo 查找
func (form *%v) Find(conn *mgo.Collection) *mgo.Query {
	condition := form.GetCondition()
	if condition != nil {
		return conn.Find(bson.M(*condition))
	}
	return nil
}`, form.Name)
}

func (form *StructInfo) goUpdateFuncs() string {
	return fmt.Sprintf(`
// UpdateID mongo 更新
func (form *%v) UpdateID(conn *mgo.Collection) error {
	condition := form.GetCondition()
	if condition != nil {
		return conn.UpdateId(form.ID, bson.M(*condition))
	}
	return nil
}`, form.Name)
}

func (form *StructInfo) goMarkFuncs() string {
	return fmt.Sprintf(`
// setBTMark 设置标记
func (form *%v) setBTMark(id uint32, set bool) {
	form.BTMark.Set(id, set)
}`, form.Name) +
		fmt.Sprintf(`
// MarkDirty 标记
func (form *%v) MarkDirty(id uint32) {
	form.BTMark.Set(id, true)
}`, form.Name) +
		fmt.Sprintf(`
// ClearDirty 去除标记
func (form *%v) ClearDirty(id uint32) {
	form.BTMark.Set(id, false)
}`, form.Name) +
		fmt.Sprintf(`
// ClearMark 去除所有标记
func (form *%v) ClearMark() {
	form.BTMark.Clear()
}`, form.Name) +
		fmt.Sprintf(`
// GetCondition 获取条件map
func (form *%v) GetCondition() *map[string]interface{} {
	condition := make(map[string]interface{})
	%v
	return &condition
}`, form.Name, form.checkCondition())
}

func (form *StructInfo) checkCondition() string {
	ostr := ""
	for _, node := range form.Nodes {
		ostr = ostr + fmt.Sprintf(`
	if form.BTMark.Is(%v) == true {
		condition["%v"] = form.%v
	}`, node.Index, toBson(node.Name, node.Type), toUpper(node.Name))
	}
	return ostr
}

//////////////////////////////////////////////////////////////////////////
func formatImport() string {
	return `
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/jslyzt/structDefine/mark"
	`
}

func formatStruct(infos *StructInfos) string {
	ostr := ""
	for _, form := range *infos {
		ostr = ostr + fmt.Sprintf(`
//////////////////////////////////////////////////////////////////////////

%v

%v`, form.goStruct(), form.goFuncs())
	}
	return ostr
}

//////////////////////////////////////////////////////////////////////////

func formatStructs(infos *StructInfos, pkg string) string {
	return fmt.Sprintf(`
package %v

import (
	%v
)

%v`, pkg, formatImport(), formatStruct(infos))
}
