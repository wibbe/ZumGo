package main

type Cell struct {
	value string
}

func (c *Cell) String() string {
	return c.value
}
