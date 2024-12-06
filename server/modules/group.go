package modules

import "sync"

type null *struct{}

type groupUsers map[string]null

type groupsMap map[string]groupUsers

type groups struct {
	sync.Mutex
	list groupsMap
}

func (s *Server) GetGroup(groupName string) groupUsers {
	s.groups.Lock()
	defer s.groups.Unlock()
	return s.groups.list[groupName]
}

var Groups = groups{list: make(groupsMap)}

func (s *Server) DeleteFromGroup(user *User) {
	s.groups.Lock()
	defer s.groups.Unlock()
	_, ok := s.groups.list[user.GroupName]
	if !ok || user.GroupName == "" {
		return
	}
	delete(s.groups.list[user.GroupName], user.Name)
}

func (s *Server) AddUserToGroup(groupName string, user *User) {
	s.groups.Lock()
	defer s.groups.Unlock()

	_, ok := s.groups.list[groupName]
	if !ok {
		s.groups.list[groupName] = map[string]null{}
	}
	user.GroupName = groupName
	s.groups.list[groupName][user.Name] = nil
}
