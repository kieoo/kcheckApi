package model

type TestSuite struct {
	Name      string      `xml:"name,attr"`
	TErrors   int32       `xml:"errors,attr"`
	Failures  int32       `xml:"failures,attr"`
	Skip      int32       `xml:"skipped,attr"`
	Tests     int32       `xml:"tests,attr"`
	Time      float32     `xml:"time,attr"`
	Timestamp string      `xml:"timestamp,attr"`
	TestCases []*TestCase `xml:"testcase"`
}

type TestCase struct {
	ClassN    string    `xml:"classname,attr"`
	Name      string    `xml:"name,attr"`
	Status    string    `xml:"status,attr"`
	Time      float32   `xml:"time,attr"`
	SystemOut SystemOut `xml:"system-out"`
}

type SystemOut struct {
	Out string `xml:",cdata"`
}
