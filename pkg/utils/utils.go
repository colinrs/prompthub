package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"runtime/debug"
)

// 生成随机字符串，长度由用户指定
func generateRandomString(length int) (string, error) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		b, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[b.Int64()]
	}
	return string(ret), nil
}

// GenerateRandomString 生成随机字符串，长度由用户指定
func GenerateRandomString(length int) (string, error) {
	return generateRandomString(length)
}

// GenerateRandomNumber 生成随机数字，长度由用户指定
func GenerateRandomNumber(length int) (string, error) {
	const numbers = "0123456789"
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		b, err := rand.Int(rand.Reader, big.NewInt(int64(len(numbers))))
		if err != nil {
			return "", err
		}
		ret[i] = numbers[b.Int64()]
	}
	return string(ret), nil
}

// GenerateRandomNumberString 生成随机数字，长度由用户指定
func GenerateRandomNumberString(length int) (string, error) {
	return GenerateRandomNumber(length)
}

func PageToOffsetLimit(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	return offset, limit
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func Stack() string {
	return string(debug.Stack())
}
