package clean

// Strings applies function on string slice
func Strings(strs []string, f func(string) string) []string {
	if len(strs) == 0 || f == nil {
		return strs
	}

	resp := make([]string, len(strs))
	for i, v := range strs {
		resp[i] = f(v)
	}
	return resp
}
