package dbbr

import (
	"testing"
)

func TestConvertToLines(t *testing.T) {
	expectedResults := []Line{
		Line{0, "check_string: hello"},
		Line{1, "	+ world"},
		Line{2, "		+ goodbye"},
		Line{3, "			+ :False:"},
		Line{4, "			- :True:"},
		Line{5, "		- john"},
		Line{6, ""},
	}
	strLines := []string{
		"check_string: hello",
		"	+ world",
		"		+ goodbye",
		"			+ :False:",
		"			- :True:",
		"		- john",
		"",
	}
	lines := ConvertToLines(strLines)
	compareLineSlices(lines, expectedResults, t)
}

func TestReadLinesFromFile(t *testing.T) {
	demo1File := "demo/demo1.dbbr"
	lines, err := ReadLinesFromFile(demo1File)
	if err != nil {
		t.Fatalf("error opening demo file: %v", err)
	}
	expectedResults := []Line{
		Line{1, "# Demo 1:"},
		Line{2, "# A simple single rule, to check if string contains:"},
		Line{3, "#	\"hello\" and ( ( \"world\" and not \"goodbye\" ) or \"john\" )"},
		Line{4, "# The parser's logic:"},
		Line{5, "# 	input.contains(word)"},
		Line{6, ""},
		Line{7, "check_string: hello"},
		Line{8, "	+ world"},
		Line{9, "		+ goodbye"},
		Line{10, "			+ :False:"},
		Line{11, "			- :True:"},
		Line{12, "		- john"},
		Line{13, ""},
	}
	compareLineSlices(lines, expectedResults, t)
}

func TestSectionLines(t *testing.T) {
	t.Run("testOneSection", testOneSection)
	t.Run("testMultipleSections", testMultipleSections)
	t.Run("testMultipleSectionsSpaced", testMultipleSectionsSpaced)
}

func testOneSection(t *testing.T) {
	expectedResults := [][]Line{[]Line{
		Line{0, "check_string: hello"},
		Line{1, "	+ world"},
		Line{2, "		+ goodbye"},
		Line{3, "			+ :False:"},
		Line{4, "			- :True:"},
		Line{5, "		- john"},
	}}
	lines := []Line{
		Line{0, "check_string: hello"},
		Line{1, "	+ world"},
		Line{2, "		+ goodbye"},
		Line{3, "			+ :False:"},
		Line{4, "			- :True:"},
		Line{5, "		- john"},
		Line{6, ""},
	}
	sections := SectionLines(lines)
	compareSectionSlices(sections, expectedResults, t)
}

func testMultipleSections(t *testing.T) {
	expectedResults := [][]Line{
		[]Line{
			Line{0, "check_string: hello"},
			Line{1, "	+ world"},
			Line{2, "		+ goodbye"},
			Line{3, "			+ :False:"},
			Line{4, "			- :True:"},
			Line{5, "		- john"},
		},
		[]Line{
			Line{7, "check_integer: i < 53"},
			Line{8, "	+ i > 24"},
			Line{9, "		+ :False:"},
			Line{10, "	- i < 100"},
			Line{11, "		- i > 203"},
		},
		[]Line{
			Line{13, "check_rules: true"},
			Line{14, "	+ :check_string:"},
			Line{15, "		- :check_integer:"},
			Line{16, "	- :False:"},
		},
	}
	lines := []Line{
		Line{0, "check_string: hello"},
		Line{1, "	+ world"},
		Line{2, "		+ goodbye"},
		Line{3, "			+ :False:"},
		Line{4, "			- :True:"},
		Line{5, "		- john"},
		Line{6, ""},
		Line{7, "check_integer: i < 53"},
		Line{8, "	+ i > 24"},
		Line{9, "		+ :False:"},
		Line{10, "	- i < 100"},
		Line{11, "		- i > 203"},
		Line{12, ""},
		Line{13, "check_rules: true"},
		Line{14, "	+ :check_string:"},
		Line{15, "		- :check_integer:"},
		Line{16, "	- :False:"},
	}
	sections := SectionLines(lines)
	compareSectionSlices(sections, expectedResults, t)
}

func testMultipleSectionsSpaced(t *testing.T) {
	expectedResults := [][]Line{
		[]Line{
			Line{3, "check_string: hello"},
			Line{4, "	+ world"},
			Line{5, "		+ goodbye"},
			Line{6, "			+ :False:"},
			Line{7, "			- :True:"},
			Line{8, "		- john"},
		},
		[]Line{
			Line{13, "check_integer: i < 53"},
			Line{14, "	+ i > 24"},
			Line{15, "		+ :False:"},
			Line{16, "	- i < 100"},
			Line{17, "		- i > 203"},
		},
		[]Line{
			Line{21, "check_rules: true"},
			Line{22, "	+ :check_string:"},
			Line{23, "		- :check_integer:"},
			Line{24, "	- :False:"},
		},
	}
	lines := []Line{
		Line{0, ""},
		Line{1, ""},
		Line{2, ""},
		Line{3, "check_string: hello"},
		Line{4, "	+ world"},
		Line{5, "		+ goodbye"},
		Line{6, "			+ :False:"},
		Line{7, "			- :True:"},
		Line{8, "		- john"},
		Line{9, ""},
		Line{10, ""},
		Line{11, ""},
		Line{12, ""},
		Line{13, "check_integer: i < 53"},
		Line{14, "	+ i > 24"},
		Line{15, "		+ :False:"},
		Line{16, "	- i < 100"},
		Line{17, "		- i > 203"},
		Line{18, ""},
		Line{19, ""},
		Line{20, ""},
		Line{21, "check_rules: true"},
		Line{22, "	+ :check_string:"},
		Line{23, "		- :check_integer:"},
		Line{24, "	- :False:"},
		Line{25, ""},
		Line{26, ""},
	}
	sections := SectionLines(lines)
	compareSectionSlices(sections, expectedResults, t)

}

func compareLineSlices(lines, expectedResults []Line, t *testing.T) {
	if len(lines) != len(expectedResults) {
		t.Fatalf(
			"len mismatch. want: %v have: %v",
			len(expectedResults),
			len(lines),
		)
	}
	for i, have := range lines {
		want := expectedResults[i]
		if have.Number != want.Number {
			t.Fatalf(
				"line.Number mismatch between:\nwant: '%v'\nhave: '%v'",
				want,
				have,
			)
		}
		if have.Value != want.Value {
			t.Fatalf(
				"line.Value mismatch between:\nwant: '%v'\nhave: '%v'",
				want,
				have,
			)
		}
	}
}

func compareSectionSlices(sections, expectedResults [][]Line, t *testing.T) {
	if len(sections) != len(expectedResults) {
		t.Fatalf(
			"len mismatch. want: %v have: %v",
			len(expectedResults),
			len(sections),
		)
	}
	for i, section := range sections {
		t.Log("testing section number:", i)
		compareLineSlices(section, expectedResults[i], t)
	}
}
