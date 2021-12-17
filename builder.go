package n1qlbuilder

import (
	"bytes"
	"fmt"
	"strings"
)

type Builder struct {
	operation string
	vars      []string
	from      string
	set       string
	where     []*condition
	limit     int
	offset    int
	orderBy   string
	order     string
	returning []string
}

type condition struct {
	logic string
	test  string
}

func Query() *Builder {
	return new(Builder)
}

func (b *Builder) Select(vars ...string) *Builder {
	b.operation = "select"
	b.vars = vars
	return b
}
func (b *Builder) Update(vars ...string) *Builder {
	b.operation = "update"
	b.vars = vars
	return b
}

func (b *Builder) From(bucket string) *Builder {
	b.from = fmt.Sprintf("`%s`", bucket)
	return b
}

func (b *Builder) Set(varsToSet string) *Builder {
	b.set = fmt.Sprintf("%s", varsToSet)
	return b
}

func (b *Builder) Unnest(bucket string, arrayName string) *Builder {
	b.from += fmt.Sprintf(" UNNEST `%s`.`%s`", bucket, arrayName)
	return b
}

func (c *condition) String() string {
	return fmt.Sprintf("%s %s", c.logic, c.test)
}

func (b *Builder) Where(t string) *Builder {
	c := &condition{
		logic: "where",
		test:  t,
	}
	b.addLogic(c)
	return b
}

func (b *Builder) And(t string) *Builder {
	c := &condition{
		logic: "and",
		test:  t,
	}
	b.addLogic(c)
	return b
}

func (b *Builder) Or(t string) *Builder {
	c := &condition{
		logic: "or",
		test:  t,
	}
	b.addLogic(c)
	return b
}

func (b *Builder) addLogic(c *condition) {
	if b.where == nil {
		b.where = make([]*condition, 0)
	}
	b.where = append(b.where, c)
}

func (b *Builder) Limit(n int) *Builder {
	b.limit = n
	return b
}

func (b *Builder) Offset(n int) *Builder {
	b.offset = n
	return b
}

func (b *Builder) OrderBy(fieldName string, order string) *Builder {
	b.orderBy = fieldName
	b.order = order
	return b
}

func (b *Builder) Returning(vars ...string) *Builder {
	b.returning = vars
	return b
}

func (b *Builder) String() string {
	out := bytes.Buffer{}

	// main operation
	switch b.operation {
	case "select":
		out.WriteString("select ")
		out.WriteString(strings.Join(b.vars, ","))
		out.WriteString(fmt.Sprintf(" from %s", b.from))

	case "update":
		out.WriteString("update ")
		out.WriteString(strings.Join(b.vars, ","))
		out.WriteString(fmt.Sprintf(" set %s", b.set))
	}

	// where clause
	for _, c := range b.where {
		out.WriteString(" ")
		out.WriteString(c.String())
	}

	if len(b.returning) > 0 {
		out.WriteString(" returning ")
		out.WriteString(strings.Join(b.returning, ","))
	}

	// order by
	if "" != b.orderBy {
		out.WriteString(fmt.Sprintf(" order by %s %s", b.orderBy, b.order))
	}

	// limit
	if b.limit != 0 {
		out.WriteString(fmt.Sprintf(" limit %d", b.limit))
	}

	// offset
	if b.offset != 0 {
		out.WriteString(fmt.Sprintf(" offset  %d", b.offset))
	}

	return string(out.Bytes())
}

