package roles

type Registry struct {
	roleSqlName map[Role]string
}

func New() *Registry {
	r := &Registry{}

	r.roleSqlName = make(map[Role]string)
	for k, v := range roleRegistry {
		r.roleSqlName[v] = k
	}

	return r
}

func (r *Registry) FindSqlName(role Role) string {
	name, ok := r.roleSqlName[role]
	if !ok {
		panic("Role not found in sql name registry")
	}
	return name
}
