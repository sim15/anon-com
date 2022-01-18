package algebra

import "math/big"

type DDHGroup struct {
	*Group
	ExpGroup *Group
}

// new group over a specified field with generator g
func NewDDHGroup(group *Group, expGroup *Group) *DDHGroup {

	q := group.Field.Pminus1()
	q.Div(q, big.NewInt(2))
	if q.Cmp(expGroup.Field.P) != 0 {
		panic("DDH group requires the based field to be a safe prime")
	}

	return &DDHGroup{group, expGroup}
}
