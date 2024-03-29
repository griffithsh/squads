package item

// Equipment is a Component that stores the equipped items of a Character.
type Equipment struct {
	Weapon *Instance
	// RightHand *Instance
	// LeftHand  *Instance

	Helm   *Instance
	Amulet *Instance
	Armor  *Instance
	Ring1  *Instance
	Ring2  *Instance
	Belt   *Instance
	Gloves *Instance
	Boots  *Instance
}

// Type of this Component.
func (*Equipment) Type() string {
	return "Equipment"
}

// SumModifiers returns the sum of all modifiers present on all items in this
// Equipment.
func (equip *Equipment) SumModifiers() map[Modifier]float64 {
	result := map[Modifier]float64{}
	if equip == nil {
		return result
	}
	it := []*Instance{
		equip.Weapon,
		equip.Helm,
		equip.Amulet,
		equip.Armor,
		equip.Ring1,
		equip.Ring2,
		equip.Belt,
		equip.Gloves,
		equip.Boots,
	}
	for _, item := range it {
		if item == nil {
			continue
		}
		for mod, val := range item.Modifiers {
			// if this modifier does not exist already, then current will be
			// default-initialised (0) and adding it to the value of this item
			// will be correct.
			current, _ := result[mod]
			result[mod] = current + val
		}
	}
	return result
}

// WeaponClass returns the inferred ItemClass of the Weapon that is equipped (if
// one is equipped), otherwise it returns Unarmed.
func (equip *Equipment) WeaponClass() Class {
	if equip == nil {
		return UnarmedClass
	}

	if equip.Weapon == nil {
		return UnarmedClass
	}

	return equip.Weapon.Class
}

func (equip *Equipment) WeaponPreparation() int {
	if equip == nil || equip.Weapon == nil {
		// Unarmed hard-coded Prep value.
		return 400
	}
	return int(equip.Weapon.Modifiers[PreparationModifier])
}

func (equip *Equipment) WeaponActionPoints() int {
	if equip == nil || equip.Weapon == nil {
		// Unarmed hard-coded AP value.
		return 40
	}
	return int(equip.Weapon.Modifiers[ActionPointModifier])
}

func (equip *Equipment) WeaponBaseChanceToHit() float64 {
	if equip == nil || equip.Weapon == nil {
		// Unarmed hard-coded ChanceToHit
		return 0.99
	}
	return equip.Weapon.BaseChanceToHit
}
