package combat

import (
	"fmt"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/geom"
)

type skillExecutor struct {
	mgr     *ecs.World
	bus     *event.Bus
	field   *geom.Field
	archive SkillArchive
}

func newSkillExecutor(mgr *ecs.World, bus *event.Bus, field *geom.Field, archive SkillArchive) *skillExecutor {
	se := skillExecutor{
		mgr:     mgr,
		bus:     bus,
		field:   field,
		archive: archive,
	}
	se.bus.Subscribe(UsingSkill{}.Type(), se.handleUsingSkill)

	return &se
}

func (se *skillExecutor) handleUsingSkill(t event.Typer) {
	ev := t.(*UsingSkill)
	s := se.archive.Skill(ev.Skill)
	fmt.Printf("skillExecutor: Entity %d used %s on %v\n", ev.User, s.Name, ev.Selected)

	// TODO: implementation ...

	se.bus.Publish(&SkillUseConcluded{
		ev.User, ev.Skill, ev.Selected,
	})
}
