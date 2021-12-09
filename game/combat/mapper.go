package combat

type HUDData struct {
	//Current
	Background    string
	BackgroundX   int
	BackgroundY   int
	Portrait      string
	PortraitX     int
	PortraitY     int
	OverlayFrame  string
	OverlayFrameX int
	OverlayFrameY int

	Name string

	Health, HealthMax int
	Energy, EnergyMax int
	Action, ActionMax int
	Prep, PrepMax     int

	TurnQueue []QueuedParticipant

	Skills [7]UISkillInfoRow
}

type QueuedParticipant struct {
	Background    string
	BackgroundX   int
	BackgroundY   int
	Portrait      string
	PortraitX     int
	PortraitY     int
	OverlayFrame  string
	OverlayFrameX int
	OverlayFrameY int

	Prep, PrepMax int
}

func (qp QueuedParticipant) PrepPercent() int {
	return int(float64(qp.Prep) / float64(qp.PrepMax) * 26)
}

type UISkillInfoRow struct {
	Skills [2]UISkillInfo
}
type UISkillInfo struct {
	Texture string
	IconX   int
	IconY   int
	Handle  func()
}
