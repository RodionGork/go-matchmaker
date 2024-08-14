package matcher

import (
    "fmt"
    "math"
)

type minMaxAvg struct {
    sum float64
    cnt int
    min float64
    max float64
    avg float64
}

func (m *minMaxAvg) add(value float64) {
    m.sum += value
    if m.cnt == 0 {
        m.min = math.MaxFloat64
        m.max = -math.MaxFloat64
    }
    m.min = math.Min(m.min, value)
    m.max = math.Max(m.max, value)
    m.cnt += 1
    m.avg = m.sum / float64(m.cnt)
}

func (m minMaxAvg) String() string {
    // configurable format could be better
    return fmt.Sprintf("avg=%.3f, min=%.3f, max=%.3f", m.avg, m.min, m.max)
}
