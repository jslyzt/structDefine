package trans

import (
	"fmt"
	"strings"
)

var (
	char     byte   // 临时变量
	noteInfo []byte // 注释数据

	curStruct *StructInfo // 当前结构
	curNode   *StructNode // 当前节点

	memName []byte // 变量名
	memType []byte // 变量类型
	memStep uint8  // 变量阶段

	keyNote    []byte // 注释开始标记
	keyStruct  []byte // struct开始标记
	keyOfile   []byte // 输出文件标记
	keyPackage []byte // 包标记

	bNote   bool // 开始注释
	bMember bool // 成员变量赋值
)

const (
	mstepName = 1 // 名称
	mstepType = 2 // 类型
	mstepDesc = 3 // 备注
)

///////////////////////////////////////////////////////////////////////////////////
func init() {
	keyNote = []byte("//")
	keyStruct = []byte("struct")
	keyOfile = []byte(" file:")
	keyPackage = []byte(" package:")
	clear()
}

///////////////////////////////////////////////////////////////////////////////////
func clear() {
	bNote, noteInfo, memStep = false, nil, 0
	memName = make([]byte, 0)
	memType = make([]byte, 0)
}

func compare(data *[]byte, key *[]byte, start, len, length int) bool {
	if data == nil || key == nil || start+len > length {
		return false
	}
	for index := 0; index < len; index++ {
		if (*data)[start+index] != (*key)[index] {
			return false
		}
	}
	return true
}

///////////////////////////////////////////////////////////////////////////////////
func transFile(path string) (string, string) {
	data := ReadData(path)
	if data == nil || len(data) <= 0 {
		return "", ""
	}

	infos := make(StructInfos, 0)
	ofileInfo := ""
	packageInfo := ""
	nlen := len(data)
	chcekNote := false

	for nindex := 0; nindex < nlen; {
		char = data[nindex]
		nindex++

		switch char {
		case '/':
			{
				if compare(&data, &keyOfile, nindex+1, 6, nlen) == true {
					nindex += 7
					sindex := nindex
					for sindex < nlen {
						if data[sindex] == '\n' {
							break
						}
						sindex++
					}
					if sindex < nlen {
						ofileInfo = strings.Trim(string(data[nindex:sindex-1]), " ")
						nindex = sindex
					}
					continue
				}

				if compare(&data, &keyPackage, nindex+1, 9, nlen) == true {
					nindex += 10
					sindex := nindex
					for sindex < nlen {
						if data[sindex] == '\n' {
							break
						}
						sindex++
					}
					if sindex < nlen {
						packageInfo = strings.Trim(string(data[nindex:sindex-1]), " ")
						nindex = sindex
					}
					continue
				}

				if compare(&data, &keyNote, nindex-1, 2, nlen) == true {
					noteInfo = make([]byte, 0)
					bNote = true
					nindex++
					continue
				}

				continue
			}
		case 's':
			{
				if compare(&data, &keyStruct, nindex-1, 6, nlen) == true {
					nindex += 6
					curStruct = &StructInfo{}
					if noteInfo != nil && len(noteInfo) > 0 {
						curStruct.Desc = strings.Trim(string(noteInfo), " ")
						noteInfo = nil
					}
					sindex := nindex
					bindex := 0
					for sindex < nlen {
						if data[sindex] == '\n' {
							break
						}
						if data[sindex] == '{' {
							bindex = sindex
						}
						sindex++
					}
					if bindex > 0 {
						curStruct.Name = strings.Trim(string(data[nindex:bindex-1]), " ")
					} else {
						curStruct.Name = strings.Trim(string(data[nindex:sindex-1]), " ")
					}
					nindex = sindex
					continue
				}
				goto DEFAULT
			}
		case ' ', '\t':
			{
				if bNote == true && noteInfo != nil {
					noteInfo = append(noteInfo, char)
				}

				// 成员变量切换
				if bMember == true && memStep > 0 {
					memStep++
				}
				bMember = false
				continue
			}
		case '}':
			{
				if curStruct != nil {
					infos = append(infos, *curStruct)
					curStruct = nil
				}
				clear()
				continue
			}
		case '\r', '{':
			{
				continue
			}
		case '\n':
			{
				bNote = false
				if curNode != nil {

					if len(memName) > 0 {
						curNode.Name = strings.Trim(string(memName), " ")
						memName = make([]byte, 0)
					}

					if len(memType) > 0 {
						curNode.Type = strings.Trim(string(memType), " ")
						memType = make([]byte, 0)
					}

					if noteInfo != nil && len(noteInfo) > 0 {
						curNode.Desc = strings.Trim(string(noteInfo), " ")
						noteInfo = nil
					}

					if curStruct != nil {
						curNode.Index = len(curStruct.Nodes) + 1
						curStruct.Nodes = append(curStruct.Nodes, *curNode)
					}
					memStep = 0
					bMember = false
					curNode = nil
				}
			}
		default:
			{
				goto DEFAULT
			}
		}
		goto CHECKEND

	DEFAULT:
		chcekNote = false
		if curStruct != nil {
			if curNode == nil {
				curNode = &StructNode{}
				memStep = mstepName
				bMember = false
			}
			switch memStep {
			case mstepName:
				{
					memName = append(memName, char)
					bMember = true
				}
			case mstepType:
				{
					memType = append(memType, char)
					bMember = true
				}
			case mstepDesc:
				{
					chcekNote = true
				}
			}
		} else {
			chcekNote = true
		}
		if chcekNote == true && bNote == true && noteInfo != nil {
			noteInfo = append(noteInfo, char)
		}

	CHECKEND:
	}

	if len(infos) <= 0 || len(ofileInfo) <= 0 {
		return "", ""
	}

	return formatStructs(&infos, packageInfo), ofileInfo
}

///////////////////////////////////////////////////////////////////////////////////

// Trans 转换函数
func Trans(files *[]string, opath string) {
	for _, file := range *files {
		data, fpath := transFile(file)
		if len(data) > 0 {
			ofile := fmt.Sprintf("%v/%v", opath, fpath)
			SaveData(ofile, data)
			runCmd("goreturns", "-w", ofile)
		}
	}

}
