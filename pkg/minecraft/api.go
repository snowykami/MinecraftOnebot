package minecraft

type Player struct {
	GroupID  string
	UserID   string //默认存uuid，v11自动转换
	Nickname string
	Title    string
	Role     string
}

// GetServerPlayerList 获取服务器玩家列表
func (m *Manager) GetServerPlayerList(groupID string, noCache bool) (map[string]string, error) {
	return nil, nil
}
