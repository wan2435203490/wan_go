package landlord_const

type TimeoutStrategy int

const (
	DoNothing TimeoutStrategy = iota
	Pass
	PlaySmallest
)

func (ts TimeoutStrategy) GetTimeoutStrategy() string {
	return []string{"不处理", "不出", "出最小的单牌, 如果管的了上家的话"}[ts]
}
