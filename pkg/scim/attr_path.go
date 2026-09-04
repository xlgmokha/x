package scim

// AttrPath is a parsed attrPath = [URI ":"] ATTRNAME *1subAttr (RFC 7644 3.4.2.2).
type AttrPath struct {
	URI  string
	Name string
	Sub  string
}

func (p AttrPath) String() string {
	s := p.Name
	if p.Sub != "" {
		s += "." + p.Sub
	}
	if p.URI != "" {
		s = p.URI + ":" + s
	}
	return s
}
