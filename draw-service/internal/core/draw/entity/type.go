package entity

type LotteryType string

const (
	LotteryTypeLotto LotteryType = "lotto"
)

type DrawID = int32

type Draw struct {
	id DrawID
}
