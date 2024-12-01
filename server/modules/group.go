package modules

import "sync"

type null *struct{}

type groupUsers map[string]null

type groups map[string]groupUsers

type safeGroups struct {
	sync.Mutex
	List groups
}

func (sg *safeGroups) GetGroup(groupName string) groupUsers {
	sg.Lock()
	defer sg.Unlock()
	return sg.List[groupName]
}

var Groups = safeGroups{List: make(groups)}

func (groups *safeGroups) DeleteFromGroup(user *User) {
	groups.Lock()
	defer groups.Unlock()
	_, ok := groups.List[user.GroupName]
	if !ok || user.GroupName == "" {
		return
	}
	delete(groups.List[user.GroupName], user.UserName)
}

func (groups *safeGroups) AddUser(groupName string, user *User) {
	groups.Lock()
	defer groups.Unlock()

	_, ok := groups.List[groupName]
	if !ok {
		groups.List[groupName] = map[string]null{}
	}
	user.GroupName = groupName
	groups.List[groupName][user.UserName] = nil
}
