package main

import (
	"flag"
	"fmt"

	"github.com/jslyzt/structDefine/trans"
)

var (
	input  *string // 输入目录
	ifile  *string // 输入文件
	output *string // 输出目录
)

func main() {
	input = flag.String("input", "", "input file dir")
	ifile = flag.String("ifile", "", "input file path")
	output = flag.String("output", "", "output file dir")

	flag.Parse()

	if (input == nil || len(*input) <= 0) && (ifile == nil || len(*ifile) <= 0) {
		fmt.Println("please input input file path or file dir, use -h to see args")
		return
	}

	if output == nil || len(*output) <= 0 {
		fmt.Println("please input output file path or file dir, use -h to see args")
		return
	}

	files := make([]string, 0)
	if input != nil && len(*input) > 0 {
		trans.GetDirFiles(*input, &files)
	}

	if ifile != nil && len(*ifile) > 0 {
		files = append(files, *ifile)
	}

	if len(files) <= 0 {
		fmt.Println("has no input files find")
		return
	}

	trans.Trans(&files, *output)
}
