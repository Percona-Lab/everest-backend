package roles

func Get(role string) (Role, bool) {
	ro, ok := roleRegistry[role]
	return ro, ok
}
