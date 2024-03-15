package minecraft

import (
	"MCOnebot/pkg/common"
	"fmt"
	"github.com/Tnze/go-mc/chat"
	"regexp"
)

type PlayerMessage struct {
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
	text := ""
	if len(msg.With) > 0 {
		for _, v := range msg.With {
			text += v.ClearString()
		}
	} else {
		text = msg.ClearString()
	}

	params, err := FormatParams(text, regexps)
	if err != nil {
		common.Logger.Warnf("Failed to format player message: %s", err)
		return PlayerMessage{}, err
	}
	return PlayerMessage{
		Title:    params["title"],
		Username: params["player"],
		Message:  params["message"],
	}, nil
}
