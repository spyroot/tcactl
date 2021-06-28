package response

type SpecsFromStringReader interface {
	InstanceSpecsFromString(s string) (interface{}, error)
}

type SpecCreator func(s string) (interface{}, error)
