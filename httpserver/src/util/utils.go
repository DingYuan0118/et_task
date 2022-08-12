package util

import (
	"math/rand"
	"mime/multipart"
	"net/http"
)

// 使用文件前 512 字节判断文件类型
func GetFileContentType(out multipart.File) (string, error) {

    // 只需要前 512 个字节就可以了
    buffer := make([]byte, 512)

    _, err := out.Read(buffer)
    if err != nil {
        return "", err
    }

    contentType := http.DetectContentType(buffer)
    return contentType, nil
}

func Contains(slices []string, value string) bool {
    for _, s := range slices {
        if s == value {
            return true
        }
    }
    return false
}

// 生成 [start, end] 中的随机数
func GenerateRandomIdNum(start int, end int) int {
    num := rand.Intn((end - start)) + start
    return num
}