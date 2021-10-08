package components

type IdentityComponent struct {
	Tags []int
}

func (i *IdentityComponent) HasTag(searchTag int) bool {
	for _, tag := range i.Tags {
		if tag == searchTag {
			return true
		}
	}

	return false
}
