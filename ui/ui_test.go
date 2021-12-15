package ui

import (
	"image"
	"strings"
	"testing"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
)

func TestUpdate(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			uic := NewUI(strings.NewReader(tc.xml))
			uic.Data = tc.data
			mgr := ecs.NewWorld()
			mgr.AddComponent(mgr.NewEntity(), uic)
			bus := &event.Bus{}
			sys := NewUISystem(mgr, bus)
			uiScale := 2.0 // FIXME: should this be tied to zoom? Or separate?
			bus.Publish(&game.WindowSizeChanged{
				OldW: 0,
				OldH: 0,
				NewW: int(800 * uiScale),
				NewH: int(600 * uiScale),
			})
			err := sys.Update()

			// Checking
			if err != nil {
				t.Fatalf("UISystem.Update: %v", err)
			}
			if len(uic.RenderInstructions()) != len(tc.wantRenders) {
				t.Errorf("want %d renderInstructions, got %d", len(tc.wantRenders), len(uic.RenderInstructions()))
			}
			if len(uic.interactives) != len(tc.wantInteractives) {
				t.Errorf("want %d interactives, got %d", len(tc.wantInteractives), len(uic.interactives))
			}
			for _, want := range tc.wantRenders {
				found := false
				for _, got := range uic.renderinstructions {
					if want == got {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("did not find renderInstruction %v in %v", want, uic.renderinstructions)
				}
			}
			for _, want := range tc.wantInteractives {
				found := false
				for _, got := range uic.interactives {
					if want.Bounds == got.Bounds {
						found = true
						break
					}
				}
				if !found {
					clickers := []image.Rectangle{}
					for _, interact := range uic.interactives {
						clickers = append(clickers, interact.Bounds)
					}
					t.Errorf("did not find interactive %v in %v", want.Bounds, clickers)
				}
			}
		})
	}
}

type updateTest struct {
	name string
	xml  string
	data interface{}

	wantRenders      []RenderInstruction
	wantInteractives []InteractiveRegion
}

var tests = []updateTest{
	{
		"hello-world",
		`<UI><Text value="Hello, world!" /></UI>`,
		nil,
		[]RenderInstruction{
			TextRenderInstruction{
				"Hello, world!",
				TextSizeNormal,
				image.Rect(0, 0, 800, 600),
				TextLayoutLeft,
			},
		},
		nil,
	},
	{
		"clicky",
		`<UI><Button label="Click me!" width="100" onclick="Handle" /></UI>`,
		struct{ Handle func() }{Handle: func() {}},
		[]RenderInstruction{
			ButtonRenderInstruction{
				false,
				image.Rect(0, 0, 100, 15),
				"Click me!",
			},
		},
		[]InteractiveRegion{
			{
				image.Rect(0, 0, 100, 15),
				nil,
			},
		},
	},
	{
		"menu",
		`<UI align="center" valign="middle">
				<Panel width="300" height="200">
					<Image texture="image.png" x="12" y="48" width="32" height="32"/>
				</Panel>
			</UI>`,
		struct{ Handle func() }{Handle: func() {}},
		[]RenderInstruction{
			PanelRenderInstruction{
				image.Rect(250, 200, 550, 400),
			},
			ImageRenderInstruction{
				"image.png",
				image.Rect(12, 48, 44, 80),
				250, 200,
			},
		},
		nil,
	},
	{
		"padding",
		`<UI>
				<Padding>
					<Padding left="30">
						<Padding all="15">
							<Button label="Ok" width="100"/>
						</Padding>
						<Image texture="image.png" x="0" y="0" width="100" height="100"/>
					</Padding>
				</Padding>
			</UI>`,
		struct{ Handle func() }{Handle: func() {}},
		[]RenderInstruction{
			ButtonRenderInstruction{
				false,
				image.Rect(45, 15, 145, 30),
				"Ok",
			},
			ImageRenderInstruction{
				"image.png",
				image.Rect(0, 0, 100, 100),
				30, 45,
			},
		},
		[]InteractiveRegion{
			{
				image.Rect(45, 15, 145, 30),
				nil,
			},
		},
	},
	{
		"columns",
		`<UI>
			<Column twelfths="3" />
			<Column twelfths="3" />
			<Column twelfths="3">
				<Image texture="image.png" x="0" y="0" width="50" height="50" />
			</Column>
			<Column twelfths="3" />
		</UI>`,
		nil,
		[]RenderInstruction{
			ImageRenderInstruction{
				"image.png",
				image.Rect(0, 0, 50, 50),
				400, 0,
			},
		},
		nil,
	},
	// IGNORE: known bug
	// {
	// 	"range-columns",
	// 	`<UI>
	// 		<Range over="Counts">
	// 			<Column twelfths="3">
	// 				<Text value="A{{.}}" />
	// 			</Column>
	// 		</Range>
	// 		<Column twelfths="2" />
	// 		<Range over="Counts">
	// 			<Column twelfths="2">
	// 				<Text value="B{{.}}" />
	// 			</Column>
	// 		</Range>
	// 	</UI>`,
	// 	struct{ Counts [2]int }{[2]int{11, 13}},
	// 	[]RenderInstruction{
	// 		TextRenderInstruction{
	// 			"A11",
	// 			TextSizeNormal,
	// 			image.Rect(0, 0, 200, 600),
	// 			TextLayoutLeft,
	// 		},
	// 		TextRenderInstruction{
	// 			"A13",
	// 			TextSizeNormal,
	// 			image.Rect(200, 0, 400, 600),
	// 			TextLayoutLeft,
	// 		},
	// 		TextRenderInstruction{
	// 			"B11",
	// 			TextSizeNormal,
	// 			image.Rect(533, 0, 667, 600),
	// 			TextLayoutLeft,
	// 		},
	// 		TextRenderInstruction{
	// 			"B13",
	// 			TextSizeNormal,
	// 			image.Rect(667, 0, 800, 600),
	// 			TextLayoutLeft,
	// 		},
	// 	},
	// 	nil,
	// },
	{
		"conditional",
		`<UI>
			<Range over="Names">
				<Padding all="10">
					<If expr=".">
						<Text value="{{.}}" />
					</If>
				</Padding>
			</Range>
		</UI>`,
		struct{ Names []string }{[]string{"", "real", ""}},
		[]RenderInstruction{
			TextRenderInstruction{
				"real",
				TextSizeNormal,
				image.Rect(10, 30, 790, 42),
				TextLayoutLeft,
			},
		},
		nil,
	},
}
