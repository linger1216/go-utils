package jelly

const (
	TimestampString1       = "2019-10-13T19:52:25+08:00"
	TimestampString1Json   = `"2019-10-13T19:52:25+08:00"`
	TimestampString1Number = "1570967545"
	Timestamp1             = int64(1570967545)
)

var exampleFormats = []string{
	`"2019-10-13T19:52:25+08:00"`,
	`"2019-10-13 19:52:25+0800"`,
	`"2019-10-13 19:52:25 08:00"`,
	`"2019-10-13 19:52:25"`,
	`"2019/10/13 19:52:25 CST"`,
	`"2019/10/13 19:52:25"`,
	`"2019/10/13 11:52:25+0000"`, // UTC 时间，等同于北京时间 2019/10/13 19:52:25+0800
	`"2019/10/13 11:52:25+00:00"`,
	`"2019-10-13 11:52:25+0000"`,
	`"2019-10-13 11:52:25+00:00"`,
	`"2019-10-13 18:52:25 07:00"`, // test escaped '+'
	`"1570967545"`,
	`"1570967545000"`,
	`1570967545`,
}
