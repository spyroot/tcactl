package ui

import "github.com/jedib0t/go-pretty/v6/table"

type TableColorStyler struct {
	// default table styler
	Default table.Style
	// wide or not
	_isWide bool
	// color or not
	_isColor bool
}

func NewTableColorStyler() *TableColorStyler {
	t := TableColorStyler{}
	t._isColor = true
	t.Default = table.StyleColoredDark
	return &t
}

func (s *TableColorStyler) GetTableStyle() interface{} {
	return table.StyleColoredDark
}

func (s *TableColorStyler) IsColor() bool {
	return s._isColor
}

func (s *TableColorStyler) GetFields() []string {
	var f = []string{""}
	return f
}

func (s *TableColorStyler) IsWide() bool {
	return s._isWide
}

func (s *TableColorStyler) SetColor(c bool) {
	s._isColor = c
}

func (s *TableColorStyler) SetWide(v bool) {
	s._isWide = v
}
