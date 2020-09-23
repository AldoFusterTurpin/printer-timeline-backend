package maputil

import ()

func JoinMaps(source map[string]string, destination map[string]string) map[string]string {
	if source == nil && destination == nil {
		return make(map[string]string)
	}

	if source == nil && destination != nil {
		return destination
	}

	if destination == nil && source != nil {
		return source
	}

	for k, v := range source {
		destination[k] = v
	}
	return destination
}
