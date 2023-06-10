package workerman_go

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	ConstSignBy    = iota
	ConstTimeStamp //timestamp
)

type generateSignToken struct {
	Code  string
	Year  int
	Month int
	Day   int
	Hour  int
	By    string
}

// GenerateSignTime 签名函数,参数 data 为 key-value 键值对，不包含 sign 字段
func GenerateSignTime(data map[string]string, secretKey string, funcTime func() time.Duration) string {
	expireTime := time.Now().Unix() + int64(funcTime().Seconds())
	data[strconv.Itoa(ConstTimeStamp)] = fmt.Sprintf("%d", expireTime)
	// 将 key-value 键值对排序
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 按照 key-value 的顺序组装字符串
	dataStr := ""
	for i, key := range keys {
		if i != 0 {
			dataStr += "&"
		}
		dataStr += fmt.Sprintf("%s=%s", key, data[key])
	}

	// 在字符串末尾追加 secretKey，生成签名字符串
	//signStr := dataStr + secretKey
	sign := []byte(GenerateSignString(data, secretKey))

	// 将签名结果转换为字符串，并进行 hex 编码
	//signStr = hex.EncodeToString(sign[:])

	// 将签名字符串追加到 key-value 键值对的末尾
	dataStr += fmt.Sprintf("&sign=%s", sign)

	// 对字符串进行 hex 编码
	hexDataStr := hex.EncodeToString([]byte(dataStr))

	return hexDataStr
}

// ParseSignCode 参数 secretKey 为密钥，跟生成激活码序列号方法的 secretKey 参数一致
func ParseSignCode(code string, secretKey string) (map[string]string, error) {
	// 对激活码序列号进行 hex 解码
	data, err := hex.DecodeString(code)
	if err != nil {
		return nil, err
	}

	// 将解码后的数据转换为字符串
	dataStr := string(data)

	// 将字符串按照 & 分割成多个 key-value 键值对
	parts := strings.Split(dataStr, "&")

	// 分离签名和其它 key-value 键值对

	signPart := ""
	// 将其它 key-value 键值对再次按照 = 分割成多个键和值
	dataMap := make(map[string]string)

	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			return nil, errors.New("Invalid data")
		}
		if kv[0] != "sign" {
			dataMap[kv[0]] = kv[1]
		} else {
			signPart = kv[1]
		}
	}
	timestamp, err := strconv.ParseInt(dataMap[strconv.Itoa(ConstTimeStamp)], 10, 64)
	if err != nil {
		// 处理错误
		return dataMap, errors.New("时间转换错误")
	}

	if timestamp < time.Now().Unix() {
		return dataMap, errors.New("Time has expired")
	}

	// 生成签名字符串
	signDataStr := GenerateSignString(dataMap, secretKey)
	// 验证签名是否正确
	if signDataStr != signPart {
		return nil, errors.New("Invalid sign")
	}
	return dataMap, nil
}

// GenerateSignString 生成带有签名hex字符串
func GenerateSignString(data map[string]string, secretKey string) string {
	// 将 key-value 键值对排序
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 按照 key-value 的顺序组装字符串
	dataStr := ""
	for i, key := range keys {
		if i != 0 {
			dataStr += "&"
		}
		dataStr += fmt.Sprintf("%s=%s", key, data[key])
	}

	// 在字符串末尾追加 secretKey，生成签名字符串
	signStr := dataStr + secretKey
	sign := md5.Sum([]byte(signStr))

	// 将签名结果转换为字符串，并进行 hex 编码
	signDataStr := hex.EncodeToString(sign[:])

	return signDataStr
}
