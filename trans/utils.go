package trans

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// GetDirFiles 获取文件夹文件
func GetDirFiles(path string, files *[]string) {
	flist, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("read dir %v, error: %v\n", path, err)
		return
	}
	for _, dir := range flist {

		fname := dir.Name()
		if dir.IsDir() == true {
			GetDirFiles(fmt.Sprintf("%v/%v", path, fname), files)
			continue
		}

		lname := len(fname)
		if lname > 4 && fname[lname-4:] == ".def" {
			*files = append(*files, fmt.Sprintf("%v/%v", path, fname))
		}
	}
}

// SaveData 保存文件
func SaveData(path, data string) {
	ioutil.WriteFile(path, []byte(data), os.ModePerm)
}

// SaveBytes 保存字符数组
func SaveBytes(path string, data []byte) {
	ioutil.WriteFile(path, data, os.ModePerm)
}

// ReadData 读取文件
func ReadData(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("read file %v, error: %v\n", path, err)
		return nil
	}
	return data
}

// 首字母大写
func toUpper(str string) string {
	if len(str) <= 0 {
		return str
	}
	ostr := strings.ToUpper(string(str[0])) + string(str[1:])
	if ostr == "Id" {
		return "ID"
	}
	return ostr
}

// json风格
func toJSON(str string) string {
	lstr := len(str)
	if lstr <= 0 {
		return ""
	}
	if str == "id" || str == "ID" || str == "Id" || str == "iD" {
		return "id"
	}
	ostr := make([]byte, 0, lstr*2)
	nstr := strings.ToLower(str)
	for index := 0; index < lstr; index++ {
		if index > 0 && nstr[index] != str[index] {
			ostr = append(ostr, '_')
		}
		ostr = append(ostr, nstr[index])
	}
	return string(ostr)
}

// bson风格
func toBson(str, stp string) string {
	if stp == "bson.ObjectId" {
		if str == "id" || str == "ID" || str == "Id" || str == "iD" {
			return "_id"
		}
	}
	return toJSON(str)
}

// runCmd 执行命令
func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
