package onebot

import (
	"MCOnebot/pkg/common"
	libobv11 "MCOnebot/pkg/libonebotv11"
	libobv12 "MCOnebot/pkg/libonebotv12"
	"MCOnebot/pkg/minecraft"
)

func (m *Manager) HandleActionMux() {
	if m.V11Bot != nil {
		v11Mux := libobv11.NewActionMux()

		// 11 action
		v11Mux.HandleFunc(libobv11.ActionSendMessage, func(w libobv11.ResponseWriter, r *libobv11.Request) {
			p := libobv11.NewParamGetter(w, r)
			rmsg, _ := p.GetString("raw_message")          // 获取标准动作参数
			detailedType, _ := p.GetString("message_type") // 获取详细动作参数
			var groupID, userID string
			if detailedType == "group" {
				groupID, _ = p.GetString("group_id")
				userID = ""
			} else {
				groupID = ""
				userID, _ = p.GetString("user_id")
			}
			w.WriteData(
				map[string]interface{}{
					"message_id": minecraft.GenerateUserID(rmsg),
				})
			m.EventChan <- common.Event{
				Type:    "message",
				SubType: detailedType,
				Message: rmsg,
				GroupID: groupID,
				UserID:  userID,
				Version: 11,
			}
		})
		m.V11Bot.Handle(v11Mux)
	}

	if m.V12Bot != nil {
		v12Mux := libobv12.NewActionMux()
		// 12 action
		v12Mux.HandleFunc(libobv12.ActionSendMessage, func(w libobv12.ResponseWriter, r *libobv12.Request) {
			p := libobv12.NewParamGetter(w, r)
			rmsg, _ := p.GetString("message")               // 获取标准动作参数
			detailedType, _ := p.GetString("detailed_type") // 获取详细动作参数
			var groupID, userID string
			if detailedType == "group" {
				groupID, _ = p.GetString("group_id")
				userID = ""
			} else {
				groupID = ""
				userID, _ = p.GetString("user_id")
			}
			w.WriteData(
				map[string]interface{}{
					"message_id": minecraft.GenerateUserID(rmsg),
				})
			m.EventChan <- common.Event{
				Type:    "message",
				SubType: detailedType,
				Message: rmsg,
				GroupID: groupID,
				UserID:  userID,
				Version: 12,
			}
		})
		m.V12Bot.Handle(v12Mux)
	}
}
