package matcher

import (
    "fmt"
    "time"

    "github.com/rodiongork/go-matchmaker/pkg/utils"
)

type QueueElem struct {
    Name string
    Skill float64
    Latency float64
    Time float64
}

type Group struct {
    Id int
    Members []*QueueElem
    Time float64
}

type Matcher struct {
    groupSize int
    groupCounter int
    purgatory chan *QueueElem
    queue []*QueueElem
    debug bool
}

func New() *Matcher {
    sz := utils.IntFromEnv("GROUP_SIZE", -1)
    if sz < 1 {
        panic("Please specify GROUP_SIZE, value should be positive")
    }
    m := &Matcher{
        groupSize: sz,
        groupCounter: 0,
        queue: []*QueueElem{},
        purgatory: make(chan *QueueElem, 100),
        debug: utils.IntFromEnv("DEBUG_MATCHER", 0) != 0,
    }
    go m.Run(utils.IntFromEnv("MATCHER_PERIOD", 1))
    return m
}

func (m *Matcher) Enqueue(name string, skill float64, latency float64) {
    elem := &QueueElem{
        Name:name,
        Skill: skill,
        Latency: latency,
        Time: float64(time.Now().UnixNano()) / 1e9}
    m.purgatory <- elem
    if (m.debug) {
        fmt.Printf("enqueued: %s, skill=%.1f, latency=%.1f, time=%.2f\n",
            name, skill, latency, elem.Time)
    }
}

func (m *Matcher) purgatoryToQueue() {
    for {
        select {
            case elem := <-m.purgatory:
                m.queue = append(m.queue, elem)
            default:
                return
        }
    }
}

func (m *Matcher) makeGroups() []*Group {
    groups := make([]*Group, 0)
    for len(m.queue) >= m.groupSize {
        m.groupCounter++
        g := &Group{
            Id: m.groupCounter,
            Time: utils.UnixTimeAsFloat(),
            Members: m.queue[0:m.groupSize], // we may want to clone here
        }
        m.queue = m.queue[m.groupSize:]
        groups = append(groups, g)
    }
    return groups
}

func processGroup(group *Group) {
    fmt.Printf("Group #%d:", group.Id)
    skill := minMaxAvg{}
    latency := minMaxAvg{}
    waiting := minMaxAvg{}
    for _, user := range group.Members {
        fmt.Printf(" %s", user.Name)
        skill.add(user.Skill)
        latency.add(user.Latency)
        waiting.add(group.Time - user.Time)
    }
    fmt.Println()
    fmt.Println("\tSkills:", skill)
    fmt.Println("\tLatencies:", latency)
    fmt.Println("\tWaiting times:", waiting)
}

func (m *Matcher) Run(period int) {
    for true {
        time.Sleep(time.Second * time.Duration(period))
        m.purgatoryToQueue()
        groups := m.makeGroups()
        if m.debug {
            fmt.Println("Groups created:", len(groups), ", users still waiting:", len(m.queue))
        }
        for _, g := range groups {
            processGroup(g)
        }
    }
}
