package gallery

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/diagram"
	"github.com/egoist/mygo/ui"
)

var diagramEntries = []Entry{
	{ID: 1501, Name: "Mermaid", Category: "diagram", Summary: "Mermaid source drawn in theme colors: 33 diagram types.", Demo: demoDiagramMermaid},
}

// diagramSamples are a pottery studio's week in Mermaid.
var diagramSamples = []struct{ name, source string }{
	{"Flowchart", "flowchart LR\n  Wedge[Wedge clay] --> Throw[Throw]\n  Throw --> Dry{Bone dry?}\n" +
		"  Dry -- no --> Wait[Wait a day]\n  Wait --> Dry\n  Dry -- yes --> Bisque([Bisque fire])\n  style Bisque fill:#b4531a,stroke:#7c2d12,color:#fff"},
	{"Sequence", "sequenceDiagram\n  participant Maya as Maya Okafor\n  participant Kiln as Kiln 2\n  participant Shop as Web shop\n" +
		"  Maya->>Kiln: Load 24 mugs\n  Kiln-->>Maya: Cone 6 reached\n  Maya->>Shop: List 22 mugs\n  Shop-->>Maya: 9 sold by Friday"},
	{"Pie", "pie title October sales\n  \"Mugs\" : 46\n  \"Bowls\" : 28\n  \"Vases\" : 17\n  \"Planters\" : 9"},
	{"State", "stateDiagram-v2\n  [*] --> Loading\n  Loading --> Firing: door shut\n  Firing --> Cooling: cone 6\n" +
		"  Cooling --> Loading: under 50°C\n  Firing --> Fault: thermocouple lost\n  Fault --> [*]"},
}

func demoDiagramMermaid(c *ui.Context, log func(string)) {
	root := ui.Column(c).Width(500).Gap(16)
	at := ui.Local(root, "at", func() int { return 0 })
	root.Children(func() {
		names := make([]string, len(diagramSamples))
		for i, s := range diagramSamples {
			names[i] = s.name
		}
		if choice.Segments(c, "Diagram", names, at).Changed() {
			log("drew the " + names[*at])
		}
		diagram.Mermaid(c, diagramSamples[*at].source, diagram.MermaidOptions{Label: names[*at]})
	})
}
