package request

const (
	// SpecTypeProviderReg provider registration spec
	SpecTypeProviderReg SpecType = "provider"
)

type SpecKind string

type Spec interface {
	GetKind() *SpecKind
}
