package interfaces

import (
	"github.com/spyroot/tcactl/lib/client/specs"
	"io"
)

type SpecsFromReader interface {
	SpecsFromReader(io io.Reader, f ...specs.SpecFormatType) (*specs.RequestSpec, error)
}
