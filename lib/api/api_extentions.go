package api

import (
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/pkg/errors"
)

// ExtensionQuery - query for all extension api
func (a *TcaApi) ExtensionQuery() (*response.Extensions, error) {

	if a == nil {
		return nil, errors.NilError
	}

	return a.rest.ExtensionQuery()
}
