// wlparser
package wlparser

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	EXPORT_RE      = `(\d{2}\.\d{2}\.\d{4}),\w*\s*,(\w*)\s*,(\d+)`
	IMPORT_DATE_RE = `(\d{2}\.\d{2}\.\d{2,4})`
	IMPORT_HOUR_RE = `(\d{2}\:\d{2})\s+(\d{2}\:\d{2})\s+(\w+)`
)

var (
	ParserTypes = map[string]bool{"import": true, "export": true}
	months      = [...]time.Month{time.January,
		time.February,
		time.March,
		time.April,
		time.May,
		time.June,
		time.July,
		time.August,
		time.September,
		time.October,
		time.November,
		time.December}
)

// Basic types
type ParserTypeError string

func (e ParserTypeError) Error() string {
	return "No such parser for type: " + string(e)
}

//
// Interface definition
//
type RECompiler interface {
	compileExpr() error
}

type LineParser interface {
	RECompiler
	GetDateTimeProj(string, string) (string, int, bool)
}

//
// types definition
//
type ExportLineParser struct {
	projName string
	re       *regexp.Regexp
}

//
// Compile regular expressions for structure
//
func (t *ExportLineParser) compileExpr() (err error) {
	t.re, err = regexp.Compile(EXPORT_RE)
	return
}

//
// Parese line using re and try to retrieve date, minutes and flag
//
func (t *ExportLineParser) GetDateTimeProj(line, prevdate string) (date string, minutes int, found bool) {
	if res := t.re.FindStringSubmatch(line); len(res) == 4 {
		if found = (t.projName == "" || res[2] == t.projName); found {
			date = res[1]
			// no convrsion error expected, regexp should provide only digits
			minutes, _ = strconv.Atoi(res[3])
		}
	}
	return
}

type ImportLineParser struct {
	projName         string
	dateExp, hourExp *regexp.Regexp
}

//
// Parse line to get date and minutes for specyfic project
//
func (t *ImportLineParser) GetDateTimeProj(line, prevdate string) (date string, minutes int, found bool) {
	var res []string
	if res = t.dateExp.FindStringSubmatch(line); len(res) == 0 {
		res = t.hourExp.FindStringSubmatch(line)
	}
	switch len(res) {
	case 2:
		if date = res[1]; len(date) == 8 {
			date = date[:6] + "20" + date[6:]
		}
		found = true
		break
	case 4:
		if found = (t.projName == "" || t.projName == res[3]); found {
			hs := HourToDate(res[1])
			he := HourToDate(res[2])
			minutes = int(he.Sub(hs).Minutes())
			date = prevdate
		}
		break
	}
	return
}

//
// Compile regular expresion for this structure
//
func (t *ImportLineParser) compileExpr() (err error) {
	if t.dateExp, err = regexp.Compile(IMPORT_DATE_RE); err != nil {
		return
	}
	if t.hourExp, err = regexp.Compile(IMPORT_HOUR_RE); err != nil {
		return
	}
	return
}

//
// Factory method to create parser 
//
func NewParser(parserType, projectName string) (parser LineParser, err error) {
	// check if projecr exists
	if _, ok := ParserTypes[parserType]; !ok {
		err = ParserTypeError(parserType)
		return
	}
	switch parserType {
	case "export":
		parser = &ExportLineParser{projName: projectName}
		break
	case "import":
		parser = &ImportLineParser{projName: projectName}
		break
	}
	err = parser.compileExpr()
	return
}

//
// Results of parse
//
type ParseResult struct {
	res         map[string]int
	orderedkeys []string
}

//
// Add element to orderedkeys
//
func (t *ParseResult) addKey(key string) {
	t.orderedkeys = append(t.orderedkeys, key)
}

//
// Getter for map res
//
func (t *ParseResult) get(key string) (val int, ok bool) {
	val, ok = t.res[key]
	return
}

//
// Setter for map res
//
func (t *ParseResult) set(key string, val int) {
	t.res[key] = val
}

//
// Print prarse result on stdout
//
func (t *ParseResult) Print() {
	t.Fprint(os.Stdout)
}

//
// Print results on file
//
func (t *ParseResult) Fprint(file *os.File) {
	for _, k := range t.orderedkeys {
		v, _ := t.get(k)
		fmt.Fprintf(file, fmt.Sprintf("%-10s %10.2f\n", k, float32(v)/60))
	}
}

//
// Print gruped results on stdout
//
func (t *ParseResult) PrintWeeks() {
	t.FprintWeeks(os.Stdout)
}

//
// Print gruped results to file
//
func (t *ParseResult) FprintWeeks(file *os.File) {
	total := 0
	sum := 0
	for _, key := range t.orderedkeys {
		data := ParseDate(key)
		weekday := data.Weekday()
		if weekday == time.Monday {
			fmt.Fprintf(file, fmt.Sprintf("%21s\n%21.2f\n\n",
				"---",
				MinutesToHours(sum)))
			sum = 0
		}
		if weekday != time.Saturday && weekday != time.Sunday {
			value := t.res[key]
			fmt.Fprintf(file, fmt.Sprintf("%-10s %10.2f\n",
				key, MinutesToHours(value)))
			sum += value
			total += value
		}
	}
	fmt.Fprintf(file, fmt.Sprintf("%21s\n%21.2f\n\n",
		"---",
		MinutesToHours(sum)))

	fmt.Fprintf(file, fmt.Sprintf("%21s\n%-11s%10.2f\n\n",
		"---", "TOTAL:",
		MinutesToHours(total)))
}

//create new parseResult
func newParseResult() (result *ParseResult) {
	result = &ParseResult{make(map[string]int), []string{}}
	return
}

//
// Parse date in format dd.mm.yyyy into  time.Time object
//
func ParseDate(str string) (date time.Time) {
	var dtarray [3]int
	spl := strings.Split(str, ".")
	for id, el := range spl {
		val, _ := strconv.Atoi(el)
		dtarray[id] = val
	}
	month := IntToMonth(dtarray[1])
	date = time.Date(dtarray[2], month, dtarray[0],
		0, 0, 0, 0, time.UTC)
	return
}

//
//Convert int to time.Month
//
func IntToMonth(nbr int) (month time.Month) {
	for _, month = range months {
		if nbr == int(month) {
			return
		}
	}
	return
}

//
// Convert minutes to hours. Where 1h 30m is 1.5
//
func MinutesToHours(minutes int) (hours float32) {
	hours = float32(minutes) / 60
	return
}

//
// Convert hour string to date
//
func HourToDate(str string) (date time.Time) {
	spl := strings.Split(str, ":")
	h, _ := strconv.Atoi(spl[0])
	m, _ := strconv.Atoi(spl[1])
	return time.Date(2000, 1, 1, h, m, 0, 0, time.UTC)
}

//
// Main parse function. Collect results from LineParser into ParseResult
//
func Parse(parser LineParser, lines []string) (results *ParseResult) {
	var (
		prevdate string
	)
	results = newParseResult()
	for _, line := range lines {
		if date, min, found := parser.GetDateTimeProj(line, prevdate); found {
			prevdate = date
			if v, ok := results.get(date); ok {
				min += v
			} else {
				results.addKey(date)
			}
			results.set(date, min)
		}
	}
	return
}
