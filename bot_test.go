package main

import (
    "testing"
	"fmt"
)

func TestInjest(t *testing.T) {
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

	query := "bec de corbin"

	match := keyMatching.Closest(query)
	if _, ok := weaponsInfo.Items[match]; ok {
		fmt.Println(getWeaponMessage(match))
	}
	
    query = "mattock"

	match = keyMatching.Closest(query)
	if _, ok := weaponsInfo.Items[match]; ok {
		fmt.Println(getWeaponMessage(match))
	}
    query = "mathock"

	match = keyMatching.Closest(query)
	if _, ok := weaponsInfo.Items[match]; ok {
		fmt.Println(getWeaponMessage(match))
	}
    query = "Long-sword"

	match = keyMatching.Closest(query)
	if _, ok := weaponsInfo.Items[match]; ok {
		fmt.Println(getWeaponMessage(match))
	}
    query = "cock"

	match = keyMatching.Closest(query)
	if _, ok := monstersInfo.Items[match]; ok {
		fmt.Println(getMonsterMessage(match))
	}
    query = "scroll of light"

	match = keyMatching.Closest(query)
	if _, ok := scrollsInfo.Items[match]; ok {
		fmt.Println(getScrollMessage(match))
	}
    
    query = "wand of cancellation"

	match = keyMatching.Closest(query)
	if _, ok := wandsInfo.Items[match]; ok {
		fmt.Println(getWandMessage(match))
	}

    cost := 500

    if wands, ok := wandsByCost[cost]; ok {
        for _,wand := range wands {
            fmt.Println(wand)
        }
    }

    message := "the wand fights your attempt to write"
	match = keyMatching.Closest(message)
	if _, ok := wandsByEngraveMessage[match]; ok {
		fmt.Println(wandsByEngraveMessage[match])
	}
}

// TestMessageLengths makes sure that no single message goes over the 500
// character limit of twitch chat
func TestMessageLengths(t *testing.T) {
    limit := 500
    for k := range weaponsInfo.Items {
        messageLength := len(getWeaponMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
    for k := range armorInfo.Items {
        messageLength := len(getArmorMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
    for k := range monstersInfo.Items {
        messageLength := len(getMonsterMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
    for k := range toolsInfo.Items {
        messageLength := len(getToolMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
    for k := range amuletsInfo.Items {
        messageLength := len(getAmuletMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
    for k := range wandsInfo.Items {
        messageLength := len(getWandMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
    for k := range scrollsInfo.Items {
        messageLength := len(getScrollMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
    for k := range ringsInfo.Items {
        messageLength := len(getRingMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
    for k := range propsInfo.Items {
        messageLength := len(getPropertyMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
    for k := range comestiblesInfo.Items {
        messageLength := len(getComestibleMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
    for k := range potionsInfo.Items {
        messageLength := len(getPotionMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
    for k := range artifactsInfo.Items {
        messageLength := len(getArtifactMessage(k))
        if messageLength > limit {
            t.Errorf("Querying %s results in an oversided output of %d characters\n", k, messageLength)
        }
    }
}
