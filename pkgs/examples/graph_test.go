package examples

import (
	"gitee.com/aurora-engine/aurora/pkgs/graph"
	"testing"
)

func TestMatrixGraph_Point(t *testing.T) {
	m := graph.MatrixGraph[int]{}

	m.Point("0", "v0", 0)
	m.Point("1", "v1", 0)
	m.Point("2", "v2", 0)
	m.Point("3", "v3", 0)
	m.Point("4", "v4", 0)
	m.Point("5", "v5", 0)
	m.Point("6", "v6", 0)

	m.Drawing("0", "1", 1)
	m.Drawing("0", "2", 1)
	m.Drawing("0", "3", 1)

	m.Drawing("1", "2", 1)
	m.Drawing("1", "3", 1)
	m.Drawing("1", "4", 1)
	m.Drawing("0", "0", 1)

	m.Drawing("2", "1", 1)
	m.Drawing("2", "6", 1)
	m.Drawing("2", "0", 1)

	m.Drawing("3", "1", 1)
	m.Drawing("3", "5", 1)
	m.Drawing("3", "0", 1)

	m.Drawing("5", "6", 1)
	m.Drawing("5", "3", 1)

	m.Drawing("6", "2", 1)
	m.Drawing("6", "5", 1)

	m.Print()

	m.DFS("0")

}
