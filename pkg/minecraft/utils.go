package minecraft

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/Tnze/go-mc/chat"
	"regexp"
)

type PlayerMessage struct {
	Type     string
	Title    string
	Username string
	Message  string
}

func FormatParams(text string, regexps []*regexp.Regexp) (map[string]string, error) {
	params := make(map[string]string)
	for _, re := range regexps {
		match := re.FindStringSubmatch(text)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				if i <= len(match) {
					params[name] = match[i]
				} else {
					continue
				}
			}
		}
		if params["player"] != "" {
			return params, nil
		}
	}
	return nil, fmt.Errorf("no matching template")
}

func FormatPlayerMessage(msg chat.Message, regexps []*regexp.Regexp) (PlayerMessage, error) {
	textClearWith := ""
	textWith := ""
	if len(msg.With) > 0 {
		for _, v := range msg.With {
			textClearWith += v.ClearString()
			textWith += v.String()
		}
	}

	text := msg.String()
	textClear := msg.ClearString()

	texts := []string{textWith, textClearWith, text, textClear}

	for _, text := range texts {
		params, err := FormatParams(text, regexps)
		if err != nil {
			continue
		}
		return PlayerMessage{
			Title:    params["title"],
			Username: params["player"],
			Message:  params["message"],
			Type:     params["type"],
		}, nil
	}
	return PlayerMessage{}, fmt.Errorf("no matching template")
}

func InArray(element string, array []string) bool {
	for _, v := range array {
		if v == element {
			return true
		}
	}
	return false
}

func GenerateUserID(name string) int64 {
	hasher := sha256.New()
	hasher.Write([]byte(name))
	hash := hasher.Sum(nil)
	// 将前8个字节转换为int64
	intValue := binary.BigEndian.Uint64(hash[:8])
	return int64(intValue)
}

func GetRawByList(hashString string, list []string) interface{} {
	// 从列表中遍历计算hash值，若匹配则返回
	for _, v := range list {
		if GenerateUserID(v) == GenerateUserID(hashString) {
			return v
		}
	}
	return nil
}
