package scim

import "strings"

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

// parseAttrPath splits a concatenated path. A schema URI contains ':' (illegal
// in an attribute name), so it is everything before the final ':'; the rest is
// ATTRNAME optionally followed by ".sub".
func parseAttrPath(s string) AttrPath {
	var p AttrPath
	if i := strings.LastIndex(s, ":"); i >= 0 {
		p.URI, s = s[:i], s[i+1:]
	}
	if i := strings.IndexByte(s, '.'); i >= 0 {
		p.Name, p.Sub = s[:i], s[i+1:]
	} else {
		p.Name = s
	}
	return p
}
