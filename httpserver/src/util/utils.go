package util

import (
	"mime/multipart"
	"net/http"
)

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