package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// sendResponse 发送响应
func sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error() + "\n"
		// 写入错误信息
		// 格式：-5 ERROR
		tmp := fmt.Sprintf("-%d ", len(errString)) + errString
		_, err = conn.Write([]byte(tmp))
		return err
	}
	valueLen := fmt.Sprintf("%d ", len(value))
	_, err = conn.Write(append([]byte(valueLen), value...))
	return err
}

// readLen 以空格为分隔符读取字符串长度并转换为整数
func readLen(r *bufio.Reader) (int, error) {
	// 读入字符串直到遇到空格
	tmp, e := r.ReadString(' ')
	if e != nil {
		return 0, e
	}
	// 去掉空格并转换为整数
	length, e := strconv.Atoi(strings.TrimSpace(tmp))

	if e != nil {
		return 0, e
	}
	return length, nil
}
