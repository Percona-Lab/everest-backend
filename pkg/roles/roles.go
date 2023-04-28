package roles

type Role int

const (
	ClusterCreate Role = iota
	ClusterDelete
)

var roleRegistry = map[string]Role{
	"cluster.create": ClusterCreate,
	"cluster.delete": ClusterDelete,
}
