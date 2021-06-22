package ui

type PrinterStyle interface {
	GetTableStyle() interface{}
	IsColor() bool
	GetFields() []string
	IsWide() bool
	SetWide(bool)
	SetColor(term bool)
}
