package utils

import (
	"bytes"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"math/rand"
	"time"
)

// GenFixedLengthChineseChars 指定长度随机中文字符(包含复杂字符)
func GenFixedLengthChineseChars(length int) string {

	var buf bytes.Buffer

	for i := 0; i < length; i++ {
		buf.WriteRune(rune(RandInt(19968, 40869)))
	}
	return buf.String()
}

// GenRandomLengthChineseChars 指定范围随机中文字符
func GenRandomLengthChineseChars(start, end int) string {
	length := RandInt(start, end)
	return GenFixedLengthChineseChars(length)
}

// RandStr 随机英文小写字母
func RandStr(len int) string {
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, len)
	for i := 0; i < len; i++ {
		data[i] = byte(rand.Intn(26) + 97)
	}
	return string(data)
}

// RandInt 指定范围随机 int
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

// RandInt64 指定范围随机 int64
func RandInt64(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Int63n(max-min)
}

// ReverseString 反转字符串
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// PaddingZeroForNumberStart 数字填充0 666=>000666
func PaddingZeroForNumberStart(length int, num string) string {
	inputLen := len(num)
	if inputLen == length {
		return num
	}
	if inputLen < length {
		//前缀补充0
		for i := 0; i < length-inputLen; i++ {
			num = "0" + num
		}
		return num
	} else {
		//裁剪多余位数
		return num[0:length]
	}
}

//GeneratorNanoId 获取指定长度NanoId
//alphabet 根字符串
//size 指定生成字符串大小
func GeneratorNanoId(alphabet string, size int) string {
	id, _ := gonanoid.Generate(alphabet, size)
	return id
}
