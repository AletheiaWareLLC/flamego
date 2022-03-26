package intermediate

import "fmt"

type Label struct {
	Statement
	name string
}

var id int

func NewLabel(c string) *Label {
	i := id
	id++
	return NewLabelWithName(fmt.Sprintf("#_label%d", i), c)
}

func NewLabelWithName(n, c string) *Label {
	return &Label{
		Statement: Statement{
			comment: c,
		},
		name: n,
	}
}

func (l *Label) String() string {
	return l.name + l.Statement.String()
}
