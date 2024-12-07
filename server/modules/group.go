package modules

import "sync"

type null *struct{}

type groupUsers map[string]null

type safeGroupUsers struct {
	sync.Mutex
	list groupUsers
}
type groupsMap map[string]*safeGroupUsers

type groups struct {
	sync.Mutex
	list groupsMap
}

func (s *Server) GetGroup(groupName string) groupUsers {
	s.groups.list[groupName].Lock()
	defer s.groups.list[groupName].Unlock()
	return s.groups.list[groupName].list
}

var Groups = groups{list: make(groupsMap)}

func (s *Server) DeleteFromGroup(user *User) {
	s.groups.list[user.GroupName].Lock()
	defer s.groups.list[user.GroupName].Unlock()
	_, ok := s.groups.list[user.GroupName]
	if !ok || user.GroupName == "" {
		return
	}
	delete(s.groups.list[user.GroupName].list, user.Name)
}

func (s *Server) AddUserToGroup(groupName string, user *User) {
	s.groups.Lock()
	defer s.groups.Unlock()

	_, ok := s.groups.list[groupName]
	if !ok {
		s.groups.list[groupName] = &safeGroupUsers{
			list: map[string]null{},
		}
	}
	user.GroupName = groupName
	s.groups.list[groupName].list[user.Name] = nil
}
