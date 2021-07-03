package api_errors

type TemplateNotFound struct {
	errMsg string
}

func NewTemplateNotFound(errMsg string) *TemplateNotFound {
	return &TemplateNotFound{errMsg: errMsg}
}

func (e *TemplateNotFound) Error() string {
	return " cluster template '" + e.errMsg + "' not found."
}

type TemplateInvalidType struct {
	errMsg string
}

func NewTemplateInvalidType(errMsg string) *TemplateInvalidType {
	return &TemplateInvalidType{errMsg: errMsg}
}

func (e *TemplateInvalidType) Error() string {
	return e.errMsg
}

type InvalidSpec struct {
	errMsg string
}

func NewInvalidSpec(errMsg string) *InvalidSpec {
	return &InvalidSpec{errMsg: errMsg}
}

func (e *InvalidSpec) Error() string {
	return e.errMsg
}

type DatastoreNotFound struct {
	errMsg string
}

func NewDatastoreNotFound(errMsg string) *DatastoreNotFound {
	return &DatastoreNotFound{errMsg: errMsg}
}

func (e *DatastoreNotFound) Error() string {
	return e.errMsg
}

// CatalogNotFound error raised if tenant cloud not found
type CatalogNotFound struct {
	errMsg string
}

//
func (m *CatalogNotFound) Error() string {
	return "Catalog entity '" + m.errMsg + "' not found"
}

func NewCatalogNotFound(errMsg string) *CatalogNotFound {
	return &CatalogNotFound{errMsg: errMsg}
}

type ExtensionsNotFound struct {
	errMsg string
}

func (m *ExtensionsNotFound) Error() string {
	return "extension '" + m.errMsg + "' not found"
}

func NewExtensionsNotFound(errMsg string) *ExtensionsNotFound {
	return &ExtensionsNotFound{errMsg: errMsg}
}
