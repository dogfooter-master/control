package service

import "strings"

type NameObject struct {
	FirstName  string
	MiddleName string
	LastName   string
	FullName   string
}

func (n *NameObject) GetFullName() (fullName string) {
	nameRuleList := GetConfigRuleName()
	for _, e := range nameRuleList {
		if e == "first" {
			fullName += n.FirstName
		} else if e == "middle" {
			fullName += n.MiddleName
		} else if e == "last" {
			fullName += n.LastName
		}
		fullName += " "
	}

	fullName = strings.TrimSpace(fullName)

	return
}