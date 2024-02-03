package dgo

type node struct {
	identifier string
	label      string
}

type edge struct {
	src, dst int
	label    string
}

type Diagram struct {
	nodes []node
	msgs  []edge
	rsps  []edge
}

func (d *Diagram) addNode(n node) {
	d.nodes = append(d.nodes, n)
}

func (d *Diagram) addMsg(e edge) {
	d.msgs = append(d.msgs, e)
}

func (d *Diagram) addRsp(e edge) {
	d.rsps = append(d.rsps, e)
}

func (d *Diagram) find(id string) int {
	for i, n := range d.nodes {
		if n.identifier == id {
			return i
		}
	}
	return -1
}

func (d Diagram) Generate() string {
	// TODO: generate SVG
	return ""
}
