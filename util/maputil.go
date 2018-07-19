package util

func ContainsKey(m map[string]interface{}, key string) (bool, interface{}) {
	if len(m) == 0 || m == nil {
		return false, nil
	}
	for k, v := range m {
		if k == key {
			return true, v
		}
	}
	return false, nil
}
