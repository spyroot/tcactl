package errors

import "github.com/pkg/errors"

var NilError = errors.New("uninitialized class")
var RestNil = errors.New("uninitialized class")
var SpecNil = errors.New("nil spec argument")
