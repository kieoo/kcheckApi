package model

import "encoding/xml"

type TestSuite struct {
	XMLName   xml.Name    `xml:"testsuite"`
	TErrors   int32       `xml:"errors,attr"`
	Failures  int32       `xml:"failures,attr"`
	Name      string      `xml:"name,attr"`
	Skip      int32       `xml:"skipped,attr"`
	Tests     int32       `xml:"tests,attr"`
	Time      float32     `xml:"time,attr"`
	Timestamp string      `xml:"timestamp,attr"`
	TestCases []*TestCase `xml:"testcase"`
}

type TestCase struct {
	ClassN    string  `xml:"classname,attr"`
	Name      string  `xml:"name,attr"`
	Status    string  `xml:"status,attr"`
	Time      float32 `xml:"time,attr"`
	Failure   *Failure
	SystemOut SystemOut `xml:"system-out"`
}

type Failure struct {
	XMLName xml.Name `xml:"failure,omitempty"`
	Message string   `xml:"message,attr"`
	Out     string   `xml:",cdata"`
}

type SystemOut struct {
	Out string `xml:",cdata"`
}
