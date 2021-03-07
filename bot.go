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
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    DamageSmall string `yaml:"damage-small"`
    DamageLarge string `yaml:"damage-large"`
    Material string `yaml:"material"`
}

type armors struct {
    Items map[string]armor `yaml:"armors"`
}

type armor struct {
    Type string `yaml:"skill"`
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    AC int `yaml:"ac"`
    Material string `yaml:"material"`
    Effect string `yaml:"effect"`
    MC int `yaml:"mc"`
}

type artifacts struct {
    Items map[string]artifact `yaml:"artifacts"`
}

type artifact struct {
    BaseItem string `yaml:"base_item"`
    Alignment string `yaml:"alignment"`
    Intelligent bool `yaml:"intelligent"`
    Use string `yaml:"use"`
    Carried string `yaml:"carried"`
    Used string `yaml:"used"`
    Invoked string `yaml:"invoked"`
    Obtaining string `yaml:"obtaining"`
}

type tools struct {
    Items map[string]tool `yaml:"tools"`
}

type tool struct {
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    Use string `yaml:"use"`
    Magic bool `yaml:"magic"`
}

type wands struct {
    Items map[string]wand `yaml:"wands"`
}

type wand struct {
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    Type string `yaml:"type"`
    StartingCharges string `yaml:"starting_charges"`
    Effect string `yaml:"effect"`
    Broken string `yaml:"broken"`
}

type rings struct {
    Items map[string]ring `yaml:"rings"`
}

type ring struct {
    Cost int `yaml:"cost"`
    ExtrinsicGranted string `yaml:"extrinsic_granted"`
    Notes string `yaml:"notes"`
}

type potions struct {
    Items map[string]potion `yaml:"potions"`
}

type potion struct {
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    Effect string `yaml:"effect"`
}

type properties struct {
    Items map[string]property `yaml:"properties"`
}

type property struct {
    Effect string `yaml:"effect"`
    Sources []string `yaml:"sources"`
}

type comestibles struct {
    Items map[string]comestible `yaml:"comestibles"`
}

type comestible struct {
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    NutritionalValue int `yaml:"nutritional_value"`
    Time int `yaml:"time"`
    Conduct string `yaml:"conduct"`
    Effect string `yaml:"effect"`
}

type monsters struct {
    Items map[string]monster `yaml:"monsters"`
}

type monster struct {
    Difficulty int `yaml:"difficulty"`
    Attacks string `yaml:"attacks"`
    Speed int `yaml:"speed"`
    AC int `yaml:"ac"`
    MR int `yaml:"mr"`
    Weight int `yaml:"weight"`
    Alignment string `yaml:"alignment"`
    Genocidable bool `yaml:"genocidable"`
    NutritionalValue int `yaml:"nutritional_value"`
    Size string `yaml:"size"`
    Resistances string `yaml:"resistances"`
    ResistancesConveyed string `yaml:"resistances_conveyed"`
    CorpseSafe bool `yaml:"corpse_safe"`
    Elbereth bool `yaml:"elbereth"`
    Extra string `yaml:"extra"`
}

// Get the stats message for an individual weapon
func (a *armors) getArmorMessage(name string) string {
    var output string 
    if val, ok := a.Items[name]; ok {
        // new output for luxi
        output = fmt.Sprintf("A %s has AC %d, MC %d, weight %d, costs %dzm and is made of %s. %s",
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
    var output string 
    if val, ok := w.Items[name]; ok {
        // new output for luxi
        output = fmt.Sprintf("A %s does %s/%s damage. "+ 
        "It is made of %s, weighs "+ 
        "%d, and is valued at "+
        "%dzm. Works your skill with %s.", 
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

func (a *artifacts) getArtifactMessage(name string) string {
    var output string 
    if val, ok := a.Items[name]; ok {
        
        var intelligent string
        if val.Intelligent {
            intelligent = " and intelligent"
        }
        var carried string
        if val.Carried != "" {
            carried = fmt.Sprintf(" While carried it bestows %s.", val.Carried)
        }
        var used string
        if val.Used != "" {
            used = fmt.Sprintf(" While %s it bestows %s.", val.Use, val.Used)
        }
        var invoked string
        if val.Invoked != "" {
            invoked = fmt.Sprintf(" When invoked it bestows %s.", val.Invoked)
        }
        var obtaining string
        if val.Obtaining != "" {
            obtaining = fmt.Sprintf(" %s", val.Obtaining)
        }
        output = fmt.Sprintf("%s is a %s%s artifact whose base item is a %s.%s%s%s%s",
        strings.Title(strings.ReplaceAll(name,"-"," ")), 
        val.Alignment,
        intelligent,
        strings.ReplaceAll(val.BaseItem,"-"," "),
        carried,
        used,
        invoked,
        obtaining)
        
    }
    return output
}

// Get the stats message for an individual tool
func (t *tools) getToolMessage(name string) string {
    var output string 
    if val, ok := t.Items[name]; ok {
        var magic string
        if val.Magic {
            magic = "is magical"
        } else {
            magic = "is not magical"
        }
        // new output for luxi
        output = fmt.Sprintf("A %s costs %dzm, weighs %d and %s.  It %s.",
        name,
        val.Cost,
        val.Weight,
        magic,
        val.Use)
    }
    return output
}

func (w *wands) getWandMessage(name string) string {
    var output string 
    if val, ok := w.Items[name]; ok {
        // new output for luxi
        output = fmt.Sprintf("A %s costs %dzm, weighs %d and had %s starting charges.  It's pattern is %s. "+
        "%s %s", 
        strings.ReplaceAll(name,"-"," "), 
        val.Cost, 
        val.Weight, 
        val.StartingCharges, 
        val.Type,
        val.Effect,
        val.Broken)
    }
    return output
}

func (r *rings) getRingMessage(name string) string {
    var output string 
    if val, ok := r.Items[name]; ok {
        // new output for luxi
        output = fmt.Sprintf("A %s costs %dzm and grants %s. %s.",
        strings.ReplaceAll(name,"-"," "), 
        val.Cost, 
        val.ExtrinsicGranted, 
        val.Notes)
    }
    return output
}

func (p *properties) getPropertyMessage(name string) string {
    var output string 
    if val, ok := p.Items[name]; ok {
        // new output for luxi
        output = fmt.Sprintf("%s %s.  Notable sources include: ",
        strings.Title(strings.ReplaceAll(name,"-"," ")), 
        val.Effect)
        for i,source := range val.Sources {
            output = output + source
            if i != len(val.Sources)-1 {
                output = output + "; "
            } else {
                output = output + "."
            }
            fmt.Println("> " + source)
        }
    }
    return output
}

func (c *comestibles) getComestibleMessage(name string) string {
    var output string 
    if val, ok := c.Items[name]; ok {
        // new output for luxi
        var conduct string
        if val.Conduct == "vegan" {
            conduct = ", is vegan,"
        } else if val.Conduct == "vegetarian" {
            conduct = ", is vegetarian"
        }

        output = fmt.Sprintf("A %s costs %dzm, weighs %d, takes %d time to eat%s and grants %d points of nutrition. %s",
        strings.ReplaceAll(name,"-"," "), 
        val.Cost, 
        val.Weight, 
        val.Time,
        conduct,
        val.NutritionalValue,
        val.Effect)
    }
    return output
}

func (p *potions) getPotionMessage(name string) string {
    var output string 
    if val, ok := p.Items[name]; ok {
        output = fmt.Sprintf("A potion of %s costs %dzm and weighs %d. %s",
        strings.ReplaceAll(name,"-"," "), 
        val.Cost, 
        val.Weight, 
        val.Effect)
    }
    return output
}

// Get the stats message for an individual weapon
func (m *monsters) getMonsterMessage(name string) string {
    var output string 
    if val, ok := m.Items[name]; ok {
        var genocidable string
        if val.Genocidable {
            genocidable = "genocidable"
        } else {
            genocidable = "not genocidable"
        }
        var resistances string
        if val.Resistances != "" {
            resistances = " It is resistant to " + val.Resistances + "."
        }
        var resistancesConveyed string
        if val.ResistancesConveyed != "" {
            resistancesConveyed = " It might convey resistance to " + val.ResistancesConveyed + "."
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
        output = fmt.Sprintf("A %s has difficulty %d.  It attacks are %s. It " +
        "has speed %d, %d AC, %d MR, weighs %d, has nutritional value %d " +
        "and %s alignment.  It is a %s creature. It is %s.%s%s " + 
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

// Load in the weapon stats from a yaml file
func getWands(fname string) *wands {
    var w *wands
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
func getRings(fname string) *rings {
    var r *rings
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &r)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return r
}

// Load in the weapon stats from a yaml file
func getProperties(fname string) *properties {
    var p *properties
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &p)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return p
}

// Load in the weapon stats from a yaml file
func getComestibles(fname string) *comestibles {
    var c *comestibles
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return c
}

// Load in the weapon stats from a yaml file
func getPotions(fname string) *potions {
    var p *potions
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &p)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return p
}

// Load in the weapon stats from a yaml file
func getArtifacts(fname string) *artifacts {
    var a *artifacts
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

func parseOracleMessage(c *twitch.Client, message string, user string) {
    if message == "join" {
        fmt.Printf("%s requests we %s", user, message)
        c.Join(user)
    } else if message == "depart" {
        fmt.Printf("%s requests we %s", user, message)
        c.Depart(user)
    }
}

func parseBroadcasterMessage(c *twitch.Client, message string, user string) {
    if message == "oracle-depart" {
        fmt.Printf("%s requests we %s", user, message)
        c.Depart(user)
    }
}

// parse messages to see if we should reply or not
func parseMessage(c *twitch.Client, mess twitch.PrivateMessage, w *weapons, a *armors, m *monsters, t *tools, wa *wands, r *rings, props *properties, com *comestibles, pot *potions, art *artifacts) {
    message := mess.Message
    channel := mess.Channel
    user := mess.User.Name
   
    fmt.Println(message)
    words := strings.Split(message, " ")
    
    fmt.Println(words)
    if !strings.HasPrefix(message, "!") {
        fmt.Println("didn't start with !")
        return
    } else {
        message = strings.TrimPrefix(message, "!")
    }

    // Deal with requests for the oracle's attention
    if channel == "oracleofdelphibot" {
        parseOracleMessage(c, message, user)
    }

    // Deal with special requests from broadcasters

    if user == channel {
        fmt.Printf("%s is talking in their own channel", user)
        parseBroadcasterMessage(c, message, user)
    }
    // Deal with all other messages


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
    } else if _, waok := wa.Items[message]; waok {
        c.Say(channel, wa.getWandMessage(message))
        //return
    } else if _, rok := r.Items[message]; rok {
        c.Say(channel, r.getRingMessage(message))
    } else if _, pok := props.Items[message]; pok {
        c.Say(channel, props.getPropertyMessage(message))
    } else if _, cok := com.Items[message]; cok {
        c.Say(channel, com.getComestibleMessage(message))
    } else if _, pok := pot.Items[message]; pok {
        c.Say(channel, pot.getPotionMessage(message))
    } else if _, aok := art.Items[message]; aok {
        c.Say(channel, art.getArtifactMessage(message))
    } else {
        return
    }
}

func main() {
    // load the information from yaml files containing stats
    w := getWeapons("weapons.yaml")
    a := getArmor("armor.yaml")
    m := getMonsters("monsters.yaml")
    t := getTools("tools.yaml")
    wa := getWands("wands.yaml")
    r := getRings("rings.yaml")
    props := getProperties("properties.yaml")
    c := getComestibles("comestibles.yaml")
    pot := getPotions("potions.yaml")
    art := getArtifacts("artifacts.yaml")

    // find the bot's name, channel's name and oauth from OS env vars
    bot := os.Getenv("TWITCHBOT")
    channel := os.Getenv("TWITCHCHANNEL")
    oauth := os.Getenv("TWITCHOAUTH")

    client := twitch.NewClient(bot, oauth)

    client.OnPrivateMessage(func(message twitch.PrivateMessage) {
        fmt.Printf("%s: %s\n", message.User.Name, message.Message)
        if message.User.Name != bot {
            fmt.Println("message was not from the oracle")
            parseMessage(client, message, w, a, m, t, wa, r, props, c, pot, art)
        }
    })

    client.Join(channel)

    err := client.Connect()
    if err != nil {
        panic(err)
    }
}
