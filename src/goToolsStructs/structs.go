package goToolsStructs

type Setting struct {
	TradeTargets     []TradeTargetSetting
	SubscribeTargets []SubscribeSetting
	CoinTradeRule    map[string]map[string]uint8 //交易额的最大小数精度
	MinTradeAmount   map[string]string// 交易额货币 相应的最小消耗的货币数量，不能等于
}

type TradeTargetSetting struct {
	Enable               bool
	Chip                 string
	Target               string
	Symbol               string
	Interval             uint8
	SmallPeriod          uint8
	BigPeriod            uint8
	Indicator            string
	EnableShort          bool
	EnableCover          bool
	Priority             uint8
	StopLosses           []stopLossSetting
	MaxInvestRate        float32 // 该币种的 价值不能超过总币种的占比
	TradeRecordTableName string
}

type SubscribeSetting struct {
	Enable     bool
	Chip       string
	Target     string
	SymbolName string
	Interval   uint8
	Indicators []IndicatorSetting
	DateStart  string
	TableName  string
}

type stopLossSetting struct {
	ChangeRate float32
	SELLRate   float32
}

type IndicatorSetting struct {
	Enable          bool
	Name            string
	PeriodStart     uint8
	PeriodEnd       uint8
	PeriodStep      uint8
	IntervalStart   uint8
	IntervalEnd     uint8
	IntervalStep    uint8
	TableNamePrefix string
}

