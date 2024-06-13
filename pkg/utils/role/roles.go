package role

import "encoding/json"

type Roles []Role

func NewRoles(s []string) *Roles {
	var roles Roles

	for _, role := range s {
		role, err := New(role)
		if err != nil {
			continue
		}
		roles = append(roles, role)
	}

	return &roles
}

func NewRolesFromJSONString(s string) (*Roles, error) {
	var data []string
	err := json.Unmarshal([]byte(s), &data)

	if err != nil {
		return nil, err
	}

	return NewRoles(data), nil
}

func (r *Roles) Can(rl Role) bool {
	for _, role := range *r {
		if role == rl {
			return true
		}
	}
	return false
}

func (r *Roles) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func (r *Roles) Grant(rl Role) {
	*r = append(*r, rl)
}

func (r *Roles) Revoke(rl Role) {

}