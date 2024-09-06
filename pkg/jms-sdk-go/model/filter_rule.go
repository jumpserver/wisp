package model

import (
	"regexp"
	"regexp/syntax"
	"sort"
)

var _ sort.Interface = CommandACLs{}

type CommandACLs []CommandACL

func (f CommandACLs) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f CommandACLs) Len() int {
	return len(f)
}

/*
	core 优先级的值越小，优先级越高，因此按此排序，第一个是优先级最高的
	优先级相同则 Action Deny 的优先级更高
*/

func (f CommandACLs) Less(i, j int) bool {
	switch {
	case f[i].Priority == f[j].Priority:
		return actionPriorityMap[f[i].Action] < actionPriorityMap[f[j].Action]
	default:
		return f[i].Priority < f[j].Priority
	}
}

type CommandACL struct {
	ID            string         `json:"id"`
	Action        CommandAction  `json:"action"`
	CommandGroups []CommandGroup `json:"command_groups"`
	IsActive      bool           `json:"is_active"`
	Name          string         `json:"name"`
	Priority      int            `json:"priority"`
	Reviewers     []interface{}  `json:"reviewers"`
}

func (cf *CommandGroup) CompilePattern() *regexp.Regexp {
	if cf.compiled {
		return cf.pattern
	}
	syntaxFlag := syntax.Perl
	if cf.IgnoreCase {
		syntaxFlag = syntax.Perl | syntax.FoldCase
	}
	syntaxReg, err := syntax.Parse(cf.Pattern, syntaxFlag)
	if err != nil {
		return nil
	}
	pattern, err := regexp.Compile(syntaxReg.String())
	if err == nil {
		cf.pattern = pattern
		cf.compiled = true
	}
	return pattern
}

func (sf *CommandACL) Match(cmd string) (CommandGroup, CommandAction, string) {
	for i := range sf.CommandGroups {
		item := sf.CommandGroups[i]
		pattern := item.CompilePattern()
		if pattern == nil {
			continue
		}
		found := pattern.FindString(cmd)
		if found == "" {
			continue
		}
		return item, sf.Action, found
	}
	return CommandGroup{}, ActionUnknown, ""
}

type CommandGroup struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	IgnoreCase bool   `json:"ignore_case"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Pattern    string `json:"pattern"`

	pattern  *regexp.Regexp
	compiled bool
}

type CommandAction string

const (
	ActionReject        = "reject"
	ActionAccept        = "accept"
	ActionReview        = "review"
	ActionWarning       = "warning"
	ActionNotifyAndWarn = "notify_and_warn"
	ActionUnknown       = "Unknown"
)

var (
	actionPriorityMap = map[CommandAction]int{
		ActionReject:        0,
		ActionReview:        1,
		ActionNotifyAndWarn: 2,
		ActionWarning:       3,
		ActionAccept:        4,
		ActionUnknown:       5,
	}
)
