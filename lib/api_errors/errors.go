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

type TenantNotFound struct {
	errMsg string
}

func NewTenantNotFound(errMsg string) *TenantNotFound {
	return &TenantNotFound{errMsg: errMsg}
}

func (e *TenantNotFound) Error() string {
	return " tenant '" + e.errMsg + "' not found."
}

type VimNotFound struct {
	errMsg string
}

func NewVimNotFound(errMsg string) *VimNotFound {
	return &VimNotFound{errMsg: errMsg}
}

func (e *VimNotFound) Error() string {
	return " vim '" + e.errMsg + "' not found."
}

type FileNotFound struct {
	errMsg string
}

func NewFileNotFound(errMsg string) *FileNotFound {
	return &FileNotFound{errMsg: errMsg}
}

func (e *FileNotFound) Error() string {
	return " file '" + e.errMsg + "' not found."
}

type InvalidArgument struct {
	errMsg string
}

func NewInvalidArgument(errMsg string) *InvalidArgument {
	return &InvalidArgument{errMsg: errMsg}
}

func (e *InvalidArgument) Error() string {
	return " invalid '" + e.errMsg + "' argument."
}

// InvalidVimFormat error must returned if client supplied incorrect format for vim ID
type InvalidVimFormat struct {
	errMsg string
}

func (m *InvalidVimFormat) Error() string {
	return "vim id format " + m.errMsg + " invalid. Example vmware_FB40D3DE2967483FBF9033B451DC7571"
}

// InvalidTaskId error must returned if client supplied incorrect task id
type InvalidTaskId struct {
	errMsg string
}

func NewInvalidTaskId(errMsg string) *InvalidTaskId {
	return &InvalidTaskId{errMsg: errMsg}
}

func (m *InvalidTaskId) Error() string {
	return "invalid task id " + m.errMsg + ". Example 9411f70f-d24d-4842-ab56-b7214d39d1b1"
}
