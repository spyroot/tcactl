package interfaces

import "github.com/spyroot/tcactl/lib/client/specs"

type SpecsFromFileReader interface {
	SpecsFromFile(fileName string, f ...specs.SpecFormatType) (*specs.RequestSpec, error)
}
