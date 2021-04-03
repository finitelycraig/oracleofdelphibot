package main

import (
    "testing"
	"fmt"
)

func TestIngest(t *testing.T) {
	//weaponsInfo = getWeapons("weapons.yaml")
	//armorInfo = getArmor("armor.yaml")
	//monstersInfo = getMonsters("monsters.yaml")
	//toolsInfo = getTools("tools.yaml")
	//wandsInfo = getWands("wands.yaml")
	//ringsInfo = getRings("rings.yaml")
	//propsInfo = getProperties("properties.yaml")
	//comestiblesInfo = getComestibles("comestibles.yaml")
	//potionsInfo = getPotions("potions.yaml")
	//artifactsInfo = getArtifacts("artifacts.yaml")
	//appearsAs = getAppearances("appearances.yaml")
	updateInfo()

	weapon := "bec de corbin"

	match := keyMatching.Closest(weapon)
	if _, ok := weaponsInfo.Items[match]; ok {
		fmt.Println(getWeaponMessage(match))
	}
	
    weapon = "mattock"

	match = keyMatching.Closest(weapon)
	if _, ok := weaponsInfo.Items[match]; ok {
		fmt.Println(getWeaponMessage(match))
	}
    //got := -1
    //t.Errorf("Abs(-1) = %d; want 1", got)
}
