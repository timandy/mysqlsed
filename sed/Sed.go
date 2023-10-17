package sed

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"mysqlsed/sqlutil"
	"os"
	"strings"
)

const (
	insertPrefix = "INSERT INTO"
	bufferSize   = 100 * 1024 * 1024 //100MB
)

// SedFile 处理sql文件, 将 insert 语句中的指定字符串替换为目标字符串, 然后转换为更新语句by id, 写入到目标文件
func SedFile(sourcePath string, targetPath string, replaceTokens []ReplaceToken) error {
	//打开源文件
	file, err := os.OpenFile(sourcePath, os.O_RDONLY, 0)
	if err != nil {
		msg := fmt.Sprintf("Error opening file: %v", err)
		return errors.New(msg)
	}
	defer file.Close()
	//删除目标文件
	_ = os.Remove(targetPath)
	//创建目标文件
	outFile, err2 := os.Create(targetPath)
	if err2 != nil {
		msg := fmt.Sprintf("Error opening file: %v", err2)
		return errors.New(msg)
	}
	defer outFile.Close()
	//创建目标文件写入器
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()
	//逐行处理
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, bufferSize), bufferSize)
	for scanner.Scan() {
		line := scanner.Text()
		if err3 := ReplaceLine(line, replaceTokens, outFile); err3 != nil {
			msg := fmt.Sprintf("Error proc sql: %v", err3.Error())
			return errors.New(msg)
		}
	}
	if err4 := scanner.Err(); err4 != nil {
		msg := fmt.Sprintf("Error reading file: %v", err4)
		return errors.New(msg)
	}
	return nil
}

// ReplaceLine 替换一行; 将 insert 语句中的指定字符串替换为目标字符串, 然后转换为更新语句by id, 写入到目标文件
func ReplaceLine(line string, replaceTokens []ReplaceToken, writer io.Writer) error {
	//判断是否插入语句
	if !strings.HasPrefix(line, insertPrefix) {
		return nil
	}
	//替换关键字
	result := line
	for _, token := range replaceTokens {
		if !strings.Contains(result, token.Source) {
			continue
		}
		result = strings.ReplaceAll(result, token.Source, token.Target)
	}
	if result == line {
		return nil
	}
	//解析插入语句
	tableName, columns, values, err := sqlutil.ResolveInsertStatement(result)
	if err != nil {
		return err
	}
	//与构建更新语句
	updateExprs, idValue, err := sqlutil.PrebuildUpdateStatement(line, columns, values)
	if err != nil {
		return err
	}
	//构建更新语句
	updateStatement := sqlutil.BuildUpdateStatement(tableName, updateExprs, idValue)
	updateSql := sqlutil.SerializeToSQL(updateStatement)
	//写入文件
	if _, err := writer.Write(updateSql); err != nil {
		msg := fmt.Sprintf("写入文件失败 %v", err)
		return errors.New(msg)
	}
	//写入文件
	if _, err := writer.Write([]byte(";\n")); err != nil {
		msg := fmt.Sprintf("写入文件失败 %v", err)
		return errors.New(msg)
	}
	return nil
}
