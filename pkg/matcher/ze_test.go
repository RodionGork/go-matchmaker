package matcher

import (
    "fmt"
    "testing"
)

func TestMinMaxAvg(t *testing.T) {
    x := minMaxAvg{}
    for _, v := range []float64{5, 8, 13, 1, 2, 3} {
        x.add(v)
    }
    if x.min != 1 {
        t.Error("min calc failed")
    }
    if x.max != 13 {
        t.Error("max calc failed")
    }
    if x.avg != float64(32) / float64(6) {
        t.Error("avg calc failed")
    }
}

func TestLua(t *testing.T) {
    groups := []*QueueElem {
        &QueueElem{Skill:10, Latency:8, Time: 100},
        &QueueElem{Skill:11, Latency:4, Time: 101},
        &QueueElem{Skill:13, Latency:2, Time: 102},
        &QueueElem{Skill:12, Latency:3, Time: 103},
        &QueueElem{Skill:17, Latency:1, Time: 104},
    }
    cnt, idx := groupThemWithLua(2, groups, 105, simpleLuaMatcher)
    if cnt != 2 {
        t.Error("Group count should be 2 but is", cnt)
    }
    if fmt.Sprintf("%v", idx) != "[1 1 2 2 -1]" {
        t.Errorf("wrong group indices returned: %v", idx)
    }
}
