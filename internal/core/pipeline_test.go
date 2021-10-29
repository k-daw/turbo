package core

import (
	"fmt"
	"strings"
	"testing"
	"turbo/internal/util"

	"github.com/pyr-sh/dag"
)

func TestPipelineAddTask(t *testing.T) {
	var g dag.AcyclicGraph
	g.Add("a")
	g.Add("b")
	g.Add("c")
	g.Connect(dag.BasicEdge("c", "b"))
	g.Connect(dag.BasicEdge("c", "a"))

	p := New(&g)
	topoDeps := make(util.Set)
	topoDeps.Add("build")
	deps := make(util.Set)
	deps.Add("prepare")
	p.AddTask(&Task{
		Name:     "build",
		TopoDeps: topoDeps,
		Deps:     deps,
		Run: func(cwd string) error {
			fmt.Println(cwd)
			return nil
		},
	})
	p.AddTask(&Task{
		Name:     "test",
		TopoDeps: topoDeps,
		Deps:     deps,
		Run: func(cwd string) error {
			fmt.Println(cwd)
			return nil
		},
	})
	p.AddTask(&Task{
		Name: "prepare",
	})

	if _, ok := p.Tasks["build"]; !ok {
		t.Fatal("AddTask is not adding tasks (build)")
	}

	if _, ok := p.Tasks["test"]; !ok {
		t.Fatal("AddTask is not adding tasks (test)")
	}
	errs := p.Run(nil, []string{"test"})

	for _, err := range errs {
		t.Fatalf("%v", err)
	}

	actual := strings.TrimSpace(p.TaskGraph.String())
	expected := strings.TrimSpace(leafString)
	if actual != expected {
		t.Fatalf("bad: \n\nactual---\n%s\n\n expected---\n%s", actual, expected)
	}
}

const leafString = `
a#build
  root
a#test
  root
b#build
  root
b#test
  root
c#test
  a#build
  b#build
root
`
