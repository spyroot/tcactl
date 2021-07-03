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
