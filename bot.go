package main

import (
        "fmt"
        "log"
        "io/ioutil"
        "os"
        "strings"

        "gopkg.in/yaml.v2"
        "github.com/gempir/go-twitch-irc/v2"
)

type weapons struct {
    Items map[string]weapon `yaml:"weapons"`
}

type weapon struct {
   Skill string `yaml:"skill"`
   Cost string `yaml:"cost"`
   Weight string `yaml:"weight"`
   DamageSmall string `yaml:"damage-small"`
   DamageLarge string `yaml:"damage-large"`
   Material string `yaml:"material"`
}

type armors struct {
    Items map[string]armor `yaml:"armors"`
}

type armor struct {
   Type string `yaml:"skill"`
   Cost string `yaml:"cost"`
   Weight string `yaml:"weight"`
   AC string `yaml:"ac"`
   Material string `yaml:"material"`
   Effect string `yaml:"effect"`
   MC string `yaml:"mc"`
}

type tools struct {
    Items map[string]tool `yaml:"tools"`
}

type tool struct {
   Cost string `yaml:"cost"`
   Weight string `yaml:"weight"`
   Use string `yaml:"use"`
   Magic bool `yaml:"magic"`
}

type monsters struct {
    Items map[string]monster `yaml:"monsters"`
}

type monster struct {
   Difficulty string `yaml:"difficulty"`
   Attacks string `yaml:"attacks"`
   Speed string `yaml:"speed"`
   AC string `yaml:"ac"`
   MR string `yaml:"mr"`
   Weight string `yaml:"weight"`
   Alignment string `yaml:"alignment"`
   Genocidable bool `yaml:"genocidable"`
   NutritionalValue string `yaml:"nutritional_value"`
   Size string `yaml:"size"`
   Resistances string `yaml:"resistances"`
   ResistancesConveyed string `yaml:"resistances_conveyed"`
   CorpseSafe bool `yaml:"corpse_safe"`
   Elbereth bool `yaml:"elbereth"`
   Extra string `yaml:"extra"`
}

// Get the stats message for an individual weapon
func (a *armors) getArmorMessage(name string) string {
    output := "Adventurer, I know of no weapon called " + name 
    if val, ok := a.Items[name]; ok {
        // new output for luxi
        output = fmt.Sprintf("A %s has AC %s, MC %s, weight %s, costs %s and is made of %s. %s",
        name, 
        val.AC, 
        val.MC, 
        val.Weight,
        val.Cost,
        val.Material,
        val.Effect)
    }
    return output
}

// Get the stats message for an individual weapon
func (w *weapons) getWeaponMessage(name string) string {
    output := "Adventurer, I know of no weapon called " + name 
    if val, ok := w.Items[name]; ok {
        // new output for luxi
        output = fmt.Sprintf("A %s does %s/%s damage. "+ 
        "It is made of %s, weighs "+ 
        "%s, and is valued at "+
        "%s. Works your skill with %s.", 
        name, 
        val.DamageSmall, 
        val.DamageLarge, 
        val.Material, 
        val.Weight, 
        val.Cost, 
        val.Skill)
    }
    return output
}

// Get the stats message for an individual tool
func (t *tools) getToolMessage(name string) string {
    output := "Adventurer, I know of no weapon called " + name 
    if val, ok := t.Items[name]; ok {
        var magic string
        if val.Magic {
            magic = "is magical"
        } else {
            magic = "is not magical"
        }
        // new output for luxi
        output = fmt.Sprintf("A %s costs %s zm, weighs %s and %s.  It %s ", 
        name, 
        val.Cost, 
        val.Weight, 
        magic, 
        val.Use)
    }
    return output
}

// Get the stats message for an individual weapon
func (m *monsters) getMonsterMessage(name string) string {
    output := "Adventurer, I know of no weapon called " + name 
    if val, ok := m.Items[name]; ok {
        var genocidable string
        if val.Genocidable {
            genocidable = "genocidable"
        } else {
            genocidable = "not genocidable"
        }
        var resistances string
        if val.Resistances != "" {
            resistances = "It is resistant to " + val.Resistances + "."
        }
        var resistancesConveyed string
        if val.ResistancesConveyed != "" {
            resistancesConveyed = "It might convey resistance to " + val.ResistancesConveyed + "."
        }
        var corpseSafe string
        if val.CorpseSafe {
            corpseSafe = "safe"
        } else {
            corpseSafe = "not safe"
        }
        var elbereth string
        if val.Elbereth {
            elbereth = "respects"
        } else {
            elbereth = "does not respect"
        }

        // new output for luxi
        output = fmt.Sprintf("A %s has difficulty %s.  It attacks for %s, " +
        "has speed %s, %s AC, %s MR, weighs %s, has nutritional value %s " +
        "and %s alignment.  It is a %s creature.  It is %s.%s%s " + 
        "Its corpse is %s to eat. It %s Elbereth.%s",
        name, 
        val.Difficulty, 
        val.Attacks,
        val.Speed,
        val.AC,
        val.MR,
        val.Weight,
        val.NutritionalValue,
        val.Alignment,
        val.Size,
        genocidable,
        resistances,
        resistancesConveyed,
        corpseSafe,
        elbereth,
        val.Extra)
    }
    return output
}

// Load in the weapon stats from a yaml file
func getWeapons(fname string) *weapons {
    var w *weapons
	yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &w)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

	return w
}

// Load in the weapon stats from a yaml file
func getArmor(fname string) *armors {
    var a *armors
	yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &a)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

	return a
}

// Load in the weapon stats from a yaml file
func getMonsters(fname string) *monsters {
    var m *monsters
	yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &m)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

	return m
}

// Load in the weapon stats from a yaml file
func getTools(fname string) *tools {
    var t *tools
	yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &t)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

	return t
}

// parse messages to see if we should reply or not
func parseMessage(c *twitch.Client, channel, message string, w *weapons, a *armors, m *monsters, t *tools) {
    fmt.Println(message)
    words := strings.Split(message, " ")
    fmt.Println(words)
    if !strings.HasPrefix(message, "!") {
        fmt.Println("didn't start with !")
        return
    } else {
        message = strings.TrimPrefix(message, "!")
    }

    if _, wok := w.Items[message]; wok {
        c.Say(channel, w.getWeaponMessage(message))
    } else if _, aok := a.Items[message]; aok {
        c.Say(channel, a.getArmorMessage(message))
        //return
    } else if _, mok := m.Items[message]; mok {
        c.Say(channel, m.getMonsterMessage(message))
        //return
    } else if _, tok := t.Items[message]; tok {
        c.Say(channel, t.getToolMessage(message))
        //return
    } else {
        return 
    }
}

func main() {
    fmt.Println("")
        // load the information from yaml files containing stats
		w := getWeapons("weapons.yaml")
		a := getArmor("armor.yaml")
		m := getMonsters("monsters.yaml")
		t := getTools("tools.yaml")

        // find the bot's name, channel's name and oauth from OS env vars
        bot := os.Getenv("TWITCHBOT")
        channel := os.Getenv("TWITCHCHANNEL")
        oauth := os.Getenv("TWITCHOAUTH")

        client := twitch.NewClient(bot, oauth)

        client.OnPrivateMessage(func(message twitch.PrivateMessage) {
            fmt.Printf("%s: %s\n", message.User.Name, message.Message)
            if message.User.Name != bot {
                fmt.Println("message was not from the oracle")
                parseMessage(client, channel, message.Message, w, a, m, t)
            }
        })

        client.Join(channel)

        err := client.Connect()
        if err != nil {
            panic(err)
        }
}
