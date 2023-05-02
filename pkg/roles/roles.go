package roles

type Role int

const (
	ClusterAdmin Role = iota
	ClusterCreate
	ClusterDelete
)

var roleRegistry = map[string]Role{
	"cluster.admin":  ClusterAdmin,
	"cluster.create": ClusterCreate,
	"cluster.delete": ClusterDelete,
}
