package matcher

import (
    _ "embed"
    "fmt"
    "time"

    "github.com/Shopify/go-lua"

    "github.com/rodiongork/go-matchmaker/pkg/utils"
)

//go:embed simple_matcher.lua
var simpleLuaMatcher string

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

func (m *Matcher) makeGroupsAndReduceQueue(groupCount int, groupIndices []int, ts float64) []*Group {
    groups := make([]*Group, groupCount)
    for i := 0; i < groupCount; i++ {
        m.groupCounter++
        groups[i] = &Group {
            Id: m.groupCounter,
            Members: make([]*QueueElem, 0, m.groupSize),
            Time: ts,
        }
    }
    newQueue := make([]*QueueElem, 0, len(m.queue) - groupCount * m.groupSize)
    for i := 0; i < len(m.queue); i++ {
        g := groupIndices[i] - 1
        if g >= 0 {
            groups[g].Members = append(groups[g].Members, m.queue[i])
        } else {
            newQueue = append(newQueue, m.queue[i])
        }
    }
    m.queue = newQueue
    return groups
}

// return count of groups created from N queue elements
// and array of N indices, telling which element goes to which group (-1 means not groupped)
func (m *Matcher) GroupThem(queue []*QueueElem, ts float64) (groupCount int, indices []int) {
    groupCount = len(queue) / m.groupSize
    indices = make([]int, len(queue))
    for i := 0; i < groupCount * m.groupSize; i++ {
        indices[i] = i / m.groupSize
    }
    for i := groupCount * m.groupSize; i < len(indices); i++ {
        indices[i] = -1
    }
    return
}

func groupThemWithLua(groupSize int, queue []*QueueElem, ts float64) (groupCount int, indices []int) {
    st := lua.NewState()
    lua.OpenLibraries(st)
    st.PushInteger(groupSize)
    st.SetGlobal("group_size")
    st.CreateTable(len(queue), 0)
    for i, user := range queue {
        st.PushInteger(i + 1)
        st.CreateTable(4, 0)
        st.PushInteger(1)
        st.PushNumber(user.Skill)
        st.SetTable(-3)
        st.PushInteger(2)
        st.PushNumber(user.Latency)
        st.SetTable(-3)
        st.PushInteger(3)
        st.PushNumber(ts - user.Time)
        st.SetTable(-3)
        st.SetTable(-3)
    }
    st.SetGlobal("users")
    err := lua.DoString(st, simpleLuaMatcher)
    if err != nil {
        fmt.Println("Error on Matching:", err.Error())
    }
    grpIdx := make([]int, len(queue))
    st.Global("users")
    for i, _ := range queue {
        st.PushInteger(i + 1)
        st.Table(-2)
        if st.TypeOf(-1) != lua.TypeTable {
            v, _ := st.ToNumber(-1)
            fmt.Println("Error on retrieve match results from 'groups' - table broken", v);
        }
        st.PushInteger(4)
        st.Table(-2)
        v, ok := st.ToInteger(-1)
        if !ok {
            fmt.Println("Error on retrieve match results from 'users' - not int");
        }
        grpIdx[i] = v
        st.Pop(2)
    }
    st.Pop(1)
    st.Global("group_count")
    grpCnt, ok := st.ToInteger(-1)
    if !ok {
        fmt.Println("Error on retrieve match results from 'group_count'");
    }
    return grpCnt, grpIdx
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
        ts := utils.UnixTimeAsFloat()
        groupCount, groupIndices := groupThemWithLua(m.groupSize, m.queue, ts)
        groups := m.makeGroupsAndReduceQueue(groupCount, groupIndices, ts)
        if m.debug {
            fmt.Println("Groups created:", len(groups), ", users still waiting:", len(m.queue))
        }
        for _, g := range groups {
            processGroup(g)
        }
    }
}

