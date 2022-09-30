package vc

import "fmt"

// stringify take key and map and return value if we expect a string.
func stringify(k string, args map[string]interface{}) (string, error) {

	var strVal string
	var ok bool

	if x, found := args[k]; found {
		if strVal, ok = x.(string); !ok {
			return "", fmt.Errorf("type mistmatch for key %s", k)
		}
	} else {
		return "", fmt.Errorf("key %s not found ", k)
	}

	return strVal, nil
}

// Find key in interface and return value as interface
func Find(obj interface{}, key string) (interface{}, bool) {

	//if the argument is not a map, ignore it
	mobj, ok := obj.(map[string]interface{})
	if !ok {
		return nil, false
	}

	for k, v := range mobj {
		if k == key {
			return v, true
		}
		if m, ok := v.(map[string]interface{}); ok {
			if res, ok := Find(m, key); ok {
				return res, true
			}
		}
		if va, ok := v.([]interface{}); ok {
			for _, a := range va {
				if res, ok := Find(a, key); ok {
					return res, true
				}
			}
		}
	}

	return nil, false
}
