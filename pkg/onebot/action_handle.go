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
			msg, _ := p.GetMessage("message")            // 获取标准动作参数
			detailType, _ := p.GetString("message_type") // 获取详细动作参数
			var groupID, userID int64
			if detailType == "group" {
				groupID, _ = p.GetInt64("group_id")
				userID = 0
			} else {
				common.Logger.Infof("user_id: %v", p)
				groupID = 0
				userID, _ = p.GetInt64("user_id")
			}
			w.WriteData(
				map[string]interface{}{
					"message_id": minecraft.GenerateIntID(msg.ExtractText()),
				})
			m.SendChan <- common.Event{
				Type:       "message",
				DetailType: detailType,
				Message:    msg.ExtractText(),
				GroupID:    groupID,
				UserID:     userID,
				Version:    11,
			}
		})

		v11Mux.HandleFunc(libobv11.ActionSendGroupMessage, func(w libobv11.ResponseWriter, r *libobv11.Request) {
			p := libobv11.NewParamGetter(w, r)
			msg, _ := p.GetMessage("message") // 获取标准动作参数
			detailType := "group"
			groupID, _ := p.GetInt64("group_id")
			w.WriteData(
				map[string]interface{}{
					"message_id": minecraft.GenerateIntID(msg.ExtractText()),
				})
			m.SendChan <- common.Event{
				Type:       "message",
				DetailType: detailType,
				Message:    msg.ExtractText(),
				GroupID:    groupID,
				Version:    11,
			}
		})

		v11Mux.HandleFunc(libobv11.ActionSendPrivateMessage, func(w libobv11.ResponseWriter, r *libobv11.Request) {
			p := libobv11.NewParamGetter(w, r)
			msg, _ := p.GetMessage("message") // 获取标准动作参数
			detailType := "group"
			groupID, _ := p.GetInt64("group_id")
			w.WriteData(
				map[string]interface{}{
					"message_id": minecraft.GenerateIntID(msg.ExtractText()),
				})
			m.SendChan <- common.Event{
				Type:       "message",
				DetailType: detailType,
				Message:    msg.ExtractText(),
				GroupID:    groupID,
				Version:    11,
			}
		})
		m.V11Bot.Handle(v11Mux)
	}

	if m.V12Bot != nil {
		v12Mux := libobv12.NewActionMux()
		// 12 action
		v12Mux.HandleFunc(libobv12.ActionSendMessage, func(w libobv12.ResponseWriter, r *libobv12.Request) {
			p := libobv12.NewParamGetter(w, r)
			rmsg, _ := p.GetMessage("message")          // 获取标准动作参数
			detailType, _ := p.GetString("detail_type") // 获取详细动作参数
			var groupID, userID string
			if detailType == "group" {
				groupID, _ = p.GetString("group_id")
				userID = ""
			} else {
				groupID = ""
				userID, _ = p.GetString("user_id")
			}
			w.WriteData(
				map[string]interface{}{
					"message_id": minecraft.GenerateIntID(rmsg.ExtractText()),
				})
			m.SendChan <- common.Event{
				Type:       "message",
				DetailType: detailType,
				Message:    rmsg.ExtractText(),
				GroupName:  groupID,
				Username:   userID,
				Version:    12,
			}
		})
		m.V12Bot.Handle(v12Mux)
	}
}
