package modules

import "sync"

type null *struct{}

type groupsMap map[string]map[string]null

type SafeGroups struct {
	sync.Mutex
	List groupsMap
}

var Groups = SafeGroups{List: make(groupsMap)}

func (groups *SafeGroups) DeleteFromGroup(user *User) {
	groups.Lock()
	defer groups.Unlock()
	_, ok := groups.List[user.GroupName]
	if !ok || user.GroupName == "" {
		return
	}
	delete(groups.List[user.GroupName], user.UserName)
}

func (groups *SafeGroups) SetGroup(groupName string, user *User) {
	groups.Lock()
	defer groups.Unlock()

	_, ok := groups.List[groupName]
	if !ok {
		groups.List[groupName] = map[string]null{}
	}
	user.GroupName = groupName
	groups.List[groupName][user.UserName] = nil
}
