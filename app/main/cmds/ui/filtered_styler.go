package ui

import "github.com/jedib0t/go-pretty/v6/table"

type FilteredOutputStyler struct {
	Fields  []string
	_isWide bool
}

func NewFilteredOutputStyler(fields []string) *FilteredOutputStyler {
	t := FilteredOutputStyler{}
	t.Fields = fields
	return &t
}

func (s *FilteredOutputStyler) GetTableStyle() interface{} {
	return table.StyleDefault
}

func (s *FilteredOutputStyler) GetFields() []string {
	return s.Fields
}

func (s *FilteredOutputStyler) IsColor() bool {
	return true
}

func (s *FilteredOutputStyler) SetColor(c bool) {
}

func (s *FilteredOutputStyler) IsWide() bool {
	return s._isWide
}

func (s *FilteredOutputStyler) SetWide(v bool) {
	s._isWide = v
}
