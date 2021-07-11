package interfaces

import "github.com/spyroot/tcactl/lib/client/specs"

type SpecsFromStringReader interface {
	SpecsFromString(s string, f ...specs.SpecFormatType) (*specs.RequestSpec, error)
}

type SpecCreator func(s string) (*specs.RequestSpec, error)
