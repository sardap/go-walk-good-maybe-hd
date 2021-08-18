package components

type IdentityComponent struct {
	Tags []string
}

func (i *IdentityComponent) HasTag(searchTag string) bool {
	for _, tag := range i.Tags {
		if tag == searchTag {
			return true
		}
	}

	return false
}
