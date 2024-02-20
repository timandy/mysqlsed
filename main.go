package main

import (
	"fmt"
	"os"

	"mysqlsed/sed"
)

func main() {
	if len(os.Args) < 3 {
		helpMsg := `mysqlsed 替换 SQL 中的特殊文本, 将更新语句写入到 out.sql
usage: mysqlsed <InputFile> [Source>>>Target] [Source>>>Target]`
		fmt.Println(helpMsg)
		return
	}
	filePath := os.Args[1]
	if len(filePath) == 0 {
		fmt.Println("<InputFile> can not be empty")
		return
	}

	//解析替换词
	sourceTargets := os.Args[2:]
	var replaceTokens []sed.ReplaceToken
	for _, st := range sourceTargets {
		token := sed.ResolveReplaceToken(st)
		replaceTokens = append(replaceTokens, token)
	}

	//执行替换
	if err := sed.SedFile(filePath, "out.sql", replaceTokens); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("完成!!!")
}
