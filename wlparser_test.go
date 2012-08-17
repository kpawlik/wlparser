package wlparser

import (
	"testing"
)

func TestExportParser(t *testing.T) {
	line1 := "06.06.2012,KPA,MMMM ,540,0,1"
	line2 := "07.06.2012,KPA,MMMM ,560,0,1"
	parser, _ := NewParser("export", "MMMM")
	parser.compileExpr()
	exDate := "06.06.2012"
	t.Log(parser)
	date, min, found := parser.GetDateTimeProj(line1, exDate)
	if !found {
		t.Error("1 - Shoul be found")
	}
	if date != exDate {
		t.Error("1 - Date not correct")
	}
	if min != 540 {
		t.Error("1 - Minutes not correct")
	}

	date, min, found = parser.GetDateTimeProj(line2, exDate)
	if !found {
		t.Error("2 - Shoul be found")
	}
	if date != "07.06.2012" {
		t.Error("2 - Date not correct")
	}
	if min != 560 {
		t.Error("3 - Minutes not correct")
	}
	line3 := "07.06.2012,"
	date, min, found = parser.GetDateTimeProj(line3, exDate)
	if found {
		t.Error("3 - Shouldn't be found")
	}
	if date != "" {
		t.Error("3 - Date found")
	}
	if min != 0 {
		t.Error("3 - Minutes found")
	}
}

func TestImportParser(t *testing.T) {
	line1 := "06.06.2012"
	line2 := "08:00 10:00 MMMM W"
	line3 := "asdfsdf asdfadfadf  sdfas"
	//line3 := "08:00 10:00 MMMM W"
	//line4 := "08:00 10:00 XXX W"
	parser, _ := NewParser("import", "MMMM")
	parser.compileExpr()
	exDate := "06.06.2012"
	date, min, found := parser.GetDateTimeProj(line1, exDate)

	if !found {
		t.Error("3 - Shoul be found")
	}
	if date != exDate {
		t.Error("3 - Date not correct")
	}
	if min != 0 {
		t.Error("3 - Minutes not correct")
	}

	date, min, found = parser.GetDateTimeProj(line2, exDate)
	if !found {
		t.Error("4 - Shoul be found")
	}
	if date != exDate {
		t.Error("4 - Date not correct")
	}
	if min != 120 {
		t.Error("4 - Minutes not correct")
	}
	date, min, found = parser.GetDateTimeProj(line3, exDate)
	if found {
		t.Error("5 - Should'nt be found")
	}
	if date != "" {
		t.Error("5 - Date not correct")
	}
	if min != 0 {
		t.Error("5 - Minutes not correct")
	}
}

func TestNewParserFunction(t *testing.T) {
	var (
		p   LineParser
		err error
	)
	p, err = NewParser("export", "")
	if p == nil && err != nil {
		t.Error("Parser should create without error")
	}
	p, err = NewParser("some", "")
	if p != nil && err == nil {
		t.Error("Parser should return error")
	}
}
func TestParseF(t *testing.T) {
	lines := []string{"PROJECT MM Mercator Services (29.05.2012 - 27.06.2012)",
		"Created 28.5.2012",
		"",
		"Date,User,Project,Minutes,MINUTES,DAYS",
		"",
		"29.05.2012,KPA,MMMERCSRV ,300,300,0",
		"30.05.2012,KPA,MMMERCSRV ,570,0,1",
		"31.05.2012,KPA,MMMERCSRV ,615,0,1",
		"01.06.2012,KPA,MMMERCSRV ,225,225,0",
		"02.06.2012,KPA,MMMERCSRV ,0,0,0",
		"03.06.2012,KPA,MMMERCSRV ,0,0,0"}
	resultLength := 6
	parser, _ := NewParser("export", "")

	pr := Parse(parser, lines)
	v, ok := pr.get("29.05.2012")
	if !ok {
		t.Error("parse result test value not found")
	}
	if v != 300 {
		t.Error("parse result test wrong value")
	}
	if len(pr.orderedkeys) != resultLength {
		t.Error("parse result test wrong length")
	}
	if len(pr.res) != resultLength {
		t.Error("parse result test wrong length")
	}
}
func BenchmarkExportParser(b *testing.B) {
	b.StopTimer()
	line1 := "06.06.2012,KPA,MMMM ,540,0,1"
	parser, _ := NewParser("export", "")
	parser.compileExpr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		parser.GetDateTimeProj(line1, "01.01.2012")
	}

}

func BenchmarkImportParser(b *testing.B) {
	b.StopTimer()
	line1 := "08:00 10:00 MMMM W"
	parser, _ := NewParser("import", "")
	parser.compileExpr()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		parser.GetDateTimeProj(line1, "01.01.2012")
	}
}
