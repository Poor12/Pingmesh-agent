package pingmesh_agent

type Storage struct {
	pinglist	*pinglist
}

func NewStorage() *Storage {
	return &Storage{
		pinglist: &pinglist{},
	}
}
