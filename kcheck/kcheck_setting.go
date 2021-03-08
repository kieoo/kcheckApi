package kcheck

// Checkers are all the checkers being ready to use.
var Checkers map[string]Checker

func init() {
	// The checkers which are ready to use should be initialized here.
	Checkers = make(map[string]Checker)
	//Checkers["RunningOnDifferentNodes"] = &RunningOnDifferentNodes{}
	// Deploy Checker
	Checkers["DWithGracefulTermination"] = &DWithGracefulTermination{}
	Checkers["DWithLivenessCheck"] = &DWithLivenessCheck{}
	Checkers["DWithResourceRequestAndLimit"] = &DWithResourceRequestAndLimit{}
	Checkers["DWithReadiness"] = &DWithReadiness{}
	Checkers["WithRollingUpdate"] = &WithRollingUpdate{}

	// Stateful Checker
	Checkers["SWithGracefulTermination"] = &SWithGracefulTermination{}
	Checkers["SWithLivenessCheck"] = &SWithLivenessCheck{}
	Checkers["SWithResourceRequestAndLimit"] = &SWithResourceRequestAndLimit{}
	Checkers["SWithReadiness"] = &SWithReadiness{}

}

var spotCheckSet = &Rule{
	Name:     "spot",
	Checkers: []Checker{
		//&RunningOnDifferentNodes{},
	},
}

var deploymentCheckSet = &Rule{
	Name: "deployment",
	Checkers: []Checker{
		&DWithLivenessCheck{},
		&DWithResourceRequestAndLimit{},
		&DWithReadiness{},
	},
}

var statefulsetCheckSet = &Rule{
	Name: "statefulSet",
	Checkers: []Checker{
		&SWithLivenessCheck{},
		&SWithResourceRequestAndLimit{},
		&SWithReadiness{},
	},
}

// initialize the default rule set
var ruleSet = []*Rule{
	spotCheckSet,
	deploymentCheckSet,
	statefulsetCheckSet,
}
