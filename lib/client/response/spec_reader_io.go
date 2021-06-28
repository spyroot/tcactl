package response

import "io"

type InstanceCreator interface {
	NewInstance(r io.Reader) (interface{}, error)
}
