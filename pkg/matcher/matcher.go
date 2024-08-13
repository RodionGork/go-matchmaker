package matcher

import (
    "fmt"
    "strconv"
    "time"
)

type QueueElem struct {
    Name string
    Skill float64
    Latency float64
    Time int64
}

type Matcher struct {
    GroupSize int
    Queue []*QueueElem
}

func New(params string) *Matcher {
    sz, err := strconv.Atoi(params) // for now let params be simply groupSize
    if err != nil {
        panic("Please check the group size configuration")
    }
    if sz < 1 || sz > 42 {
        panic("Group size seems suspicious: " + strconv.Itoa(sz))
    }
    m := &Matcher{GroupSize: sz, Queue: []*QueueElem{}}
    go m.Run()
    return m
}

func (m *Matcher) Enqueue(name string, skill float64, latency float64) {
    elem := &QueueElem{Name:name, Skill: skill, Latency: latency, Time: time.Now().Unix()}
    m.Queue = append(m.Queue, elem)
    println("enqueued:", name, skill, latency, elem.Time)
}

func (m *Matcher) Run() {
    //naive initial implementation, no algorithm, no sync
    for true {
        time.Sleep(time.Second)
        if len(m.Queue) >= m.GroupSize {
            fmt.Print("Matched:");
            for i := 0; i < m.GroupSize; i++ {
                fmt.Print(" " + m.Queue[i].Name)
            }
            fmt.Println()
            m.Queue = m.Queue[m.GroupSize:]
        }
    }
}
