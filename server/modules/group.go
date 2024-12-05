package modules

import "sync"

type null *struct{}

type groupUsers map[string]null

type Groups map[string]groupUsers

type SafeGroups struct {
	mu   sync.Mutex
	List Groups
}

func (sg *SafeGroups) GetGroup(groupName string) groupUsers {
	sg.mu.Lock()
	defer sg.mu.Unlock()
	return sg.List[groupName]
}

func (groups *SafeGroups) DeleteFromGroup(user *User) {
	groups.mu.Lock()
	defer groups.mu.Unlock()
	_, ok := groups.List[user.GroupName]
	if !ok || user.GroupName == "" {
		return
	}
	delete(groups.List[user.GroupName], user.UserName)
}

func (groups *SafeGroups) AddUser(groupName string, user *User) {
	groups.mu.Lock()
	defer groups.mu.Unlock()

	_, ok := groups.List[groupName]
	if !ok {
		groups.List[groupName] = map[string]null{}
	}
	user.GroupName = groupName
	groups.List[groupName][user.UserName] = nil
}
