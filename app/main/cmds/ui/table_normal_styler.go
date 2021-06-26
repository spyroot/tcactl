package ui

import "github.com/jedib0t/go-pretty/v6/table"

type TableNormalStyler struct {
	Default table.Style
	_isWide bool
}

func NewNormalStyler() *TableColorStyler {
	t := TableColorStyler{}
	t.Default = table.StyleColoredDark
	return &t
}

func (s *TableNormalStyler) GetTableStyle() interface{} {
	return table.StyleDefault
}

func (s *TableNormalStyler) IsColor() bool {
	return true
}

func (s *TableNormalStyler) GetFields() []string {
	var f = []string{""}
	return f
}

func (s *TableNormalStyler) IsWide() bool {
	return s._isWide
}

func (s *TableNormalStyler) SetWide(v bool) {
	s._isWide = v
}

func (s *TableNormalStyler) SetColor(c bool) {
}
