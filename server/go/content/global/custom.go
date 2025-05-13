package global

type Config struct {
	Moment Moment
}

type Limit struct {
	SecondLimitKey, MinuteLimitKey, DayLimitKey       string
	SecondLimitCount, MinuteLimitCount, DayLimitCount int64
}

type Moment struct {
	MaxContentLen int
	Limit         Limit
}
