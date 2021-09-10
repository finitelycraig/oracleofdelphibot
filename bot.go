package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "os"
    "strings"
    "strconv"
    "regexp"
    "time"
    
    "gopkg.in/yaml.v2"
    "github.com/gempir/go-twitch-irc/v2"
    "github.com/schollz/closestmatch"
    "github.com/juliangruber/go-intersect"
)

var keys []string
var keyMatching *closestmatch.ClosestMatch
var bagOfMessages []string
var messageMatching *closestmatch.ClosestMatch
var bagOfEngraves []string
var engraveMatching *closestmatch.ClosestMatch
var weaponsInfo *weapons
var armorInfo *armors
var monstersInfo *monsters
var toolsInfo *tools
var wandsInfo *wands
var wandsByCost map[int][]string
var wandsByEngraveMessage map[string][]string
var ringsInfo *rings
var ringsByPrice map[int][]string
var scrollsInfo *scrolls
var scrollsByCost map[int][]string
var amuletsInfo *amulets
var propsInfo *properties
var comestiblesInfo * comestibles
var potionsInfo *potions
var potionsByCost map[int][]string
var artifactsInfo *artifacts
var appearsAs *appearances
var messagesMapping *messages
var infosInfo *infos

var allowedBroadcasters *allowedChannels

type allowedChannels struct {
    Names []string `yaml:"channels"`
}

type appearances struct {
    Items map[string]string `yaml:"appearances"`
}

type messages struct {
    Items map[string]message `yaml:"messages"`
}

type message struct {
    Meaning string `yaml:"meaning"`
    Property string `yaml:"property"`
}

type infos struct {
    Items map[string]info `yaml:"infos"`
}

type info struct {
    Message string `yaml:"message"`
    Link string `yaml:"link"`
}

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
    Engrave []string `yaml:"engrave"`
    Effect string `yaml:"effect"`
    Broken string `yaml:"broken"`
}

type scrolls struct {
    Items map[string]scroll `yaml:"scrolls"`
}

type scroll struct {
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    Ink string `yaml:"ink"`
    Appearance string `yaml:"appearance"`
    MonsterUse bool `yaml:"monster_use"`
    Effect string `yaml:"effect"`
    Notes string `yaml:"notes"`
}

type rings struct {
    Items map[string]ring `yaml:"rings"`
}

type amulets struct {
    Items map[string]amulet `yaml:"amulets"`
}

type amulet struct {
    Cost int `yaml:"cost"`
    Weight int `yaml:"weight"`
    Appearance string `yaml:"appearance"`
    Effect string `yaml:"effect"`
    Notes string `yaml:"notes"`
}

type ring struct {
    Cost int `yaml:/"cost"`
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
    CorpseSafe string `yaml:"corpse_safe"`
    Elbereth bool `yaml:"elbereth"`
    Extra string `yaml:"extra"`
}

func getArmorMessage(name string) string {
    var output string 
    if val, ok := armorInfo.Items[name]; ok {
        output = fmt.Sprintf("A %s has AC %d, MC %d, weight %d, costs %dzm and is made of %s. %s",
        strings.ReplaceAll(name,"-"," "), 
        val.AC, 
        val.MC, 
        val.Weight,
        val.Cost,
        val.Material,
        val.Effect)
    }
    return output
}

func getWeaponMessage(name string) string {
    var output string 
    if val, ok := weaponsInfo.Items[name]; ok {
        output = fmt.Sprintf("A %s does %s/%s damage. "+ 
        "It is made of %s, weighs "+ 
        "%d, and is valued at "+
        "%dzm. Works your skill with %s.", 
        name, // strings.ReplaceAll(name,"-"," "), 
        val.DamageSmall, 
        val.DamageLarge, 
        val.Material, 
        val.Weight, 
        val.Cost, 
        val.Skill)
    }
    return output
}

func getArtifactMessage(name string) string {
    var output string 
    if val, ok := artifactsInfo.Items[name]; ok {
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

func getScrollMessage(name string) string {
    var output string 
    if val, ok := scrollsInfo.Items[name]; ok {
        var monsterUse string
        if val.MonsterUse {
            monsterUse = "Monsters can use it."
        } else {
            monsterUse = "Monsters do not use it."
        }
        output = fmt.Sprintf("A %s costs %dzm, weights %d and takes %s ink to write. %s %s %s",
        strings.Title(strings.ReplaceAll(name,"-"," ")), 
        val.Cost,
        val.Weight,
        val.Ink,
        monsterUse,
        val.Effect,
        val.Notes)
    }
    return output
}

func getAmuletMessage(name string) string {
    var output string 
    if val, ok := amuletsInfo.Items[name]; ok {
        output = fmt.Sprintf("A %s costs %dzm and weights %d. %s %s",
        strings.Title(strings.ReplaceAll(name,"-"," ")), 
        val.Cost,
        val.Weight,
        val.Effect,
        val.Notes)
    }
    return output
}

func getToolMessage(name string) string {
    var output string 
    if val, ok := toolsInfo.Items[name]; ok {
        var magic string
        if val.Magic {
            magic = "is magical"
        } else {
            magic = "is not magical"
        }
        output = fmt.Sprintf("A %s costs %dzm, weighs %d and %s.  It %s.",
        strings.ReplaceAll(name,"-"," "), 
        val.Cost,
        val.Weight,
        magic,
        val.Use)
    }
    return output
}

func getWandMessage(name string) string {
    var output string 
    if val, ok := wandsInfo.Items[name]; ok {
        output = fmt.Sprintf("A %s costs %dzm, weighs %d and has %s starting charges.  Its pattern is %s. "+
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

func getRingMessage(name string) string {
    var output string 
    if val, ok := ringsInfo.Items[name]; ok {
        output = fmt.Sprintf("A %s costs %dzm and grants %s. %s",
        strings.ReplaceAll(name,"-"," "), 
        val.Cost, 
        val.ExtrinsicGranted, 
        val.Notes)
    }
    return output
}

func getPropertyMessage(name string) string {
    var output string 
    if val, ok := propsInfo.Items[name]; ok {
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
            //fmt.Println("> " + source)
        }
    }
    return output
}

func getComestibleMessage(name string) string {
    var output string 
    if val, ok := comestiblesInfo.Items[name]; ok {
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

func getPotionMessage(name string) string {
    var output string 
    if val, ok := potionsInfo.Items[name]; ok {
        output = fmt.Sprintf("A potion of %s costs %dzm and weighs %d. %s",
        strings.ReplaceAll(name,"-"," "), 
        val.Cost, 
        val.Weight, 
        val.Effect)
    }
    return output
}

func getInfoMessage(name string) string {
    var output string 
    if val, ok := infosInfo.Items[name]; ok {
        output = fmt.Sprintf("%s %s",val.Message,val.Link)
    }
    return output
}

func getMonsterMessage(name string) string {
    var output string 
    if val, ok := monstersInfo.Items[name]; ok {
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
        if val.CorpseSafe == "" {
            corpseSafe = "is safe"
        } else {
            corpseSafe = val.CorpseSafe
        }
        var elbereth string
        if val.Elbereth {
            elbereth = "respects"
        } else {
            elbereth = "does not respect"
        }

        output = fmt.Sprintf("A %s has difficulty %d.  It attacks are %s. It " +
        "has speed %d, %d AC, %d MR, weighs %d, has nutritional value %d " +
        "and %s alignment.  It is a %s creature. It is %s.%s%s " + 
        "Eating its corpse %s. It %s Elbereth.%s",
        strings.ReplaceAll(name,"-"," "), 
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

func getWeapons(fname string) *weapons {
    var w *weapons
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &w)
    if err != nil {
        log.Fatalf("Unmarshal weapons: %v", err)
    }

    for k := range w.Items {
        keys = append(keys, k)
    }

    return w
}

func getArmor(fname string) *armors {
    var a *armors
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &a)
    if err != nil {
        log.Fatalf("Unmarshal armor: %v", err)
    }
    
    for k := range a.Items {
        keys = append(keys, k)
    }

    return a
}

func getMonsters(fname string) *monsters {
    var m *monsters
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &m)
    if err != nil {
        log.Fatalf("Unmarshal monsters: %v", err)
    }
    for k := range m.Items {
        keys = append(keys, k)
    }

    return m
}

func getTools(fname string) *tools {
    var t *tools
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &t)
    if err != nil {
        log.Fatalf("Unmarshal tools: %v", err)
    }
    for k := range t.Items {
        keys = append(keys, k)
    }

    return t
}

func getWands(fname string) *wands {
    var w *wands
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &w)
    if err != nil {
        log.Fatalf("Unmarshal wands: %v", err)
    }

    wandsByCost = make(map[int][]string)
    for k,v := range w.Items {
        if _, ok := wandsByCost[v.Cost]; ok {
            wandsByCost[v.Cost] = append(wandsByCost[v.Cost], k)
        } else {
            wandsByCost[v.Cost] = []string{k}
        }
    }
    
    wandsByEngraveMessage = make(map[string][]string)
    for k,v := range w.Items {
        for _,message := range v.Engrave {
            if _, ok := wandsByEngraveMessage[message]; ok {
                wandsByEngraveMessage[message] = append(wandsByEngraveMessage[message], k)
                bagOfEngraves = append(bagOfEngraves, message)
            } else {
                wandsByEngraveMessage[message] = []string{k}
                bagOfEngraves = append(bagOfEngraves, message)
            }
        }
    }

    // add these wands to the keys slice
    for k := range w.Items {
        keys = append(keys, k)
    }

    return w
}

func getScrolls(fname string) *scrolls {
    var s *scrolls
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &s)
    if err != nil {
        log.Fatalf("Unmarshal scrolls: %v", err)
    }

    scrollsByCost = make(map[int][]string)
    for k,v := range s.Items {
        if _, ok := scrollsByCost[v.Cost]; ok {
            scrollsByCost[v.Cost] = append(scrollsByCost[v.Cost], k)
        } else {
            scrollsByCost[v.Cost] = []string{k}
        }
    }

    // add these wands to the keys slice
    for k := range s.Items {
        keys = append(keys, k)
    }

    return s
}

func getRings(fname string) *rings {
    var r *rings
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &r)
    if err != nil {
        log.Fatalf("Unmarshal rings: %v", err)
    }
    for k := range r.Items {
        keys = append(keys, k)
    }

    return r
}

func getAmulets(fname string) *amulets {
    var a *amulets
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &a)
    if err != nil {
        log.Fatalf("Unmarshal rings: %v", err)
    }
    for k := range a.Items {
        keys = append(keys, k)
    }

    return a
}

func getProperties(fname string) *properties {
    var p *properties
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &p)
    if err != nil {
        log.Fatalf("Unmarshal props: %v", err)
    }
    for k := range p.Items {
        keys = append(keys, k)
    }

    return p
}

func getComestibles(fname string) *comestibles {
    var c *comestibles
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &c)
    if err != nil {
        log.Fatalf("Unmarshal comestibles: %v", err)
    }
    for k := range c.Items {
        keys = append(keys, k)
    }

    return c
}

func getPotions(fname string) *potions {
    var p *potions
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &p)
    if err != nil {
        log.Fatalf("Unmarshal potions: %v", err)
    }
    for k := range p.Items {
        keys = append(keys, k)
    }

    return p
}

func getArtifacts(fname string) *artifacts {
    var a *artifacts
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &a)
    if err != nil {
        log.Fatalf("Unmarshal artifacts: %v", err)
    }
    for k := range a.Items {
        keys = append(keys, k)
    }

    return a
}

func getAppearances(fname string) *appearances {
    var a *appearances
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &a)
    if err != nil {
        log.Fatalf("Unmarshal appearances: %v", err)
    }
    for k := range a.Items {
        keys = append(keys, k)
    }

    return a
}

func getMessages(fname string) *messages {
    var m *messages
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &m)
    if err != nil {
        log.Fatalf("Unmarshal message: %v", err)
    }
    for k := range m.Items {
        bagOfMessages = append(bagOfMessages, k)
    }

    return m
}

func getInfos(fname string) *infos {
    var i *infos
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &i)
    if err != nil {
        log.Fatalf("Unmarshal message: %v", err)
    }

    return i
}

func getAllowedChannels(fname string) *allowedChannels {
    var a *allowedChannels
    yamlFile, err := ioutil.ReadFile(fname)
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, &a)
    if err != nil {
        log.Fatalf("Unmarshal channels: %v", err)
    }

    return a
}

func parseOracleMessage(c *twitch.Client, message string, user string) bool {
    if message == "join" {
        for _, allowedChannel := range allowedBroadcasters.Names {
            if user == allowedChannel {
                c.Join(user)
                return true
            }
        }
    } else if message == "depart" {
        c.Depart(user)
        return true
    } else if message == "update" {
        updateInfo()
        return true
    }
    return false
}

func parseBroadcasterMessage(c *twitch.Client, message string, user string) bool {
    if message == "oracle-depart" {
        c.Depart(user)
        return true
    } else if message == "oracle-update" {
        updateInfo()
        return true
    }
    return false
}

func parseWandID(c *twitch.Client, channel, message, user string) {
    message = strings.TrimPrefix(message, "wandID")
    re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

    var candidates []string
    cost,err := strconv.Atoi(re.FindString(message))
    if err != nil {
        fmt.Println("oops wand ID cost fjcked up")
    } else {
        if wands, ok := wandsByCost[cost]; ok {
            if candidates == nil {
                candidates = wands
            }
        }
    }
    
    words := strings.Fields(message)
    if len(words) > 1 {
        re := regexp.MustCompile("[^a-zA-Z0-9]+")
        message = re.ReplaceAllString(message, "")
        message = engraveMatching.Closest(message)
        if wands, ok := wandsByEngraveMessage[message]; ok {
            if candidates == nil {
                candidates = wands
            } else {
                intersection := intersect.Simple(candidates, wands)
                candidates = nil
                switch intersection := intersection.(type) {
                    case []string:
                        for _, v := range intersection {
                            if candidates == nil {
                                candidates = []string{v}
                            } else {
                                candidates = append(candidates, v)
                            }
                        }
                    case string:
                        candidates = []string{intersection}
                    case []interface{}:
                        if len(intersection) == 0 {
                            output := "There aren't any wands that are that price with that engrave message"
                            c.Say(channel, output)
                            return
                        } else {
                            fmt.Printf("there are %d candidates \n", len(intersection))
                            for i,_ := range intersection {
                                candidates = append(candidates,intersection[i].(string))
                            }
                        }
                    }
                }
            }
    }

    var output string

    if len(candidates) == 0 {
        output = "There aren't any wands that are that price with that engrave message"
    } else if len(candidates) == 1 {
        c.Say(channel, getWandMessage(candidates[0]))
        return
    } else {
        output = "That could be a "
        for i,wand := range candidates {
            if i == len(candidates) - 1 {
                output = output + "or " + wand + "."
            } else {
                output = output + wand + ", "
            }
        }
    }
    c.Say(channel, output)
}

func parseScrollID(c *twitch.Client, channel, message, user string) {
    message = strings.TrimPrefix(message, "scrollID")
    re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

    var candidates []string
    cost,err := strconv.Atoi(re.FindString(message))
    if err != nil {
        fmt.Println("oops scrollID cost fjcked up")
        c.Say(channel, "The scrollID command needs to know the cost of the scroll you're interested in. Try '!scrollID 40'")
    } else {
        if scrolls, ok := scrollsByCost[cost]; ok {
            if candidates == nil {
                candidates = scrolls
            }
        }
    }
    
    var output string

    if len(candidates) == 0 {
        output = "There aren't any scrolls that are that price."
    } else if len(candidates) == 1 {
        c.Say(channel, getScrollMessage(candidates[0]))
        return
    } else {
        output = "That could be a "
        for i,scroll := range candidates {
            if i == len(candidates) - 1 {
                output = output + "or " + scroll + "."
            } else {
                output = output + scroll + ", "
            }
        }
    }
    c.Say(channel, output)
}

func parseNethackMessage(c *twitch.Client, channel, message, user string) {
    message = strings.TrimPrefix(message, "message")

    words := strings.Fields(message)
    if len(words) != 0 {
        re := regexp.MustCompile("[^a-zA-Z0-9]+")
        message = re.ReplaceAllString(message, "")
        message = messageMatching.Closest(message)
        if m, ok := messagesMapping.Items[message]; ok {
            c.Say(channel, m.Meaning)
            if m.Property != "" {
                if _, ok := propsInfo.Items[m.Property]; ok {
                    time.Sleep(500 * time.Millisecond)
                    c.Say(channel, getPropertyMessage(m.Property))
                }
            }
        }
    }
}

func parseMessage(c *twitch.Client, m twitch.PrivateMessage) {
    message := m.Message
    channel := m.Channel
    user := m.User.Name
    fmt.Println(message) 
    //words := strings.Split(message, " ")
    
    if strings.HasPrefix(message, "!") {
        fmt.Println("oops")
        message = strings.TrimPrefix(message, "!")
        words := strings.Fields(message)
        if words[0] == "wandID" {
            fmt.Printf("%s wants to ID a wand\n", user)
            parseWandID(c, channel, message, user)
            return
        } else if words[0] == "scrollID" {
            fmt.Printf("%s wants to ID a scroll\n", user)
            parseScrollID(c, channel, message, user)
            return
        } else if words[0] == "message" {
            fmt.Printf("%s wants to ID a message\n", user)
            parseNethackMessage(c, channel, message, user)
            return
        } else if _, ok := infosInfo.Items[message]; ok {
            c.Say(channel, getInfoMessage(message))
            return
        }
    } else if strings.HasPrefix(message, "?") {
        fmt.Print("prefix is ?")
        message = strings.TrimPrefix(message, "?")
    } else {
        return
    }

    // Deal with requests for the oracle's attention
    if channel == "oracleofdelphibot" {
        if ok := parseOracleMessage(c, message, user); ok {
            return
        }
    }

    // Deal with special requests from broadcasters
    if user == channel {
        if ok := parseBroadcasterMessage(c, message, user); ok {
            return
        }
    }

    // Deal with all other messages
    if _, ok := weaponsInfo.Items[message]; ok {
        c.Say(channel, getWeaponMessage(message))
    } else if _, ok := armorInfo.Items[message]; ok {
        c.Say(channel, getArmorMessage(message))
    } else if _, ok := monstersInfo.Items[message]; ok {
        c.Say(channel, getMonsterMessage(message))
    } else if _, ok := toolsInfo.Items[message]; ok {
        c.Say(channel, getToolMessage(message))
    } else if _, ok := wandsInfo.Items[message]; ok {
        c.Say(channel, getWandMessage(message))
    } else if _, ok := scrollsInfo.Items[message]; ok {
        c.Say(channel, getScrollMessage(message))
    } else if _, ok := ringsInfo.Items[message]; ok {
        c.Say(channel, getRingMessage(message))
    } else if _, ok := amuletsInfo.Items[message]; ok {
        c.Say(channel, getAmuletMessage(message))
    } else if _, ok := propsInfo.Items[message]; ok {
        c.Say(channel, getPropertyMessage(message))
    } else if _, ok := comestiblesInfo.Items[message]; ok {
        c.Say(channel, getComestibleMessage(message))
    } else if _, ok := potionsInfo.Items[message]; ok {
        c.Say(channel, getPotionMessage(message))
    } else if _, ok := artifactsInfo.Items[message]; ok {
        c.Say(channel, getArtifactMessage(message))
    } else if actualName, ok := appearsAs.Items[message]; ok {
        m.Message = "!"+actualName
        parseMessage(c, m)
    } else if _, ok := wandsByEngraveMessage[message]; ok {
        return // this is a shit way of handling this case
    } else {
        message = keyMatching.Closest(message)
        fmt.Println(message)
        m.Message = "?"+message
        parseMessage(c, m)
    }



}
func updateInfo() {
    //reset the matching variables
    keys = nil
    keyMatching = nil

    // load the information from yaml files containing stats
    allowedBroadcasters = getAllowedChannels("allowed-channels.yaml")
    weaponsInfo = getWeapons("weapons.yaml")
    armorInfo = getArmor("armor.yaml")
    monstersInfo = getMonsters("monsters.yaml")
    toolsInfo = getTools("tools.yaml")
    wandsInfo = getWands("wands.yaml")
    scrollsInfo = getScrolls("scrolls.yaml")
    ringsInfo = getRings("rings.yaml")
    amuletsInfo = getAmulets("amulets.yaml")
    propsInfo = getProperties("properties.yaml")
    comestiblesInfo = getComestibles("comestibles.yaml")
    potionsInfo = getPotions("potions.yaml")
    artifactsInfo = getArtifacts("artifacts.yaml")
    appearsAs = getAppearances("appearances.yaml")
    messagesMapping = getMessages("messages.yaml")
    infosInfo = getInfos("info.yaml")

    bagSizes := []int{2, 3, 4}
    keyMatching = closestmatch.New(keys, bagSizes)
    messageMatching = closestmatch.New(bagOfMessages, bagSizes)
    engraveMatching = closestmatch.New(bagOfEngraves, bagSizes)
    fmt.Println(keyMatching.AccuracyMutatingWords())

}

func main() {
    allowedBroadcasters = getAllowedChannels("allowed-channels.yaml")
    // load the information from yaml files containing stats
    updateInfo()
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

    // find the bot's name, channel's name and oauth from OS env vars
    bot := os.Getenv("TWITCHBOT")
    channel := os.Getenv("TWITCHCHANNEL")
    oauth := os.Getenv("TWITCHOAUTH")

    client := twitch.NewClient(bot, oauth)

    client.OnPrivateMessage(func(message twitch.PrivateMessage) {
        if message.User.Name != bot {
            parseMessage(client, message)
        }
    })

    // If the broadcaster leaves the stream then the bot should too
    client.OnUserPartMessage(func(message twitch.UserPartMessage) {
        if message.User == message.Channel {
            client.Depart(message.Channel)
        }
    })

    client.Join(channel)

    err := client.Connect()
    if err != nil {
        panic(err)
    }
}
