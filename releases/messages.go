package releases

import "github.com/charmbracelet/bubbles/table"

type deleteMsg struct {
	err error
}
type listMsg struct {
	content []table.Row
	err     error
}
type historyMsg struct {
	content []table.Row
	err     error
}

type rollbackMsg struct {
	err error
}

type upgradeMsg struct {
	err error
}

type notesMsg struct {
	err     error
	content string
}

type metadataMsg struct {
	err     error
	content string
}

type hooksMsg struct {
	err     error
	content string
}

type valuesMsg struct {
	err     error
	content string
}

type templatesMsg struct {
	err     error
	content string
}
