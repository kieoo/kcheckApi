package kcheck

// Checkers are all the checkers being ready to use.
var Checkers map[string]Checker

func init() {
	// The checkers which are ready to use should be initialized here.
	Checkers = make(map[string]Checker)
	//Checkers["RunningOnDifferentNodes"] = &RunningOnDifferentNodes{}
	// Deploy Checker
	Checkers["DWithGracefulTermination"] = &DWithGracefulTermination{}
	Checkers["DWithLiveness"] = &DWithLivenessCheck{}
	Checkers["DWithResourceRequestAndLimit"] = &DWithResourceRequestAndLimit{}
	Checkers["DWithReadiness"] = &DWithReadiness{}
	Checkers["WithRollingUpdate"] = &WithRollingUpdate{}
	Checkers["DWithEmptyDirSizeLimit"] = &DWithEmptyDirSizeLimit{}
	Checkers["DWithLivenessReadinessDelayCheck"] = &DWithLivenessReadinessDelayCheck{}
	Checkers["DWithTerminationGrace"] = &DWithTerminationGrace{}

	// Stateful Checker
	Checkers["SWithGracefulTermination"] = &SWithGracefulTermination{}
	Checkers["SWithLiveness"] = &SWithLivenessCheck{}
	Checkers["SWithResourceRequestAndLimit"] = &SWithResourceRequestAndLimit{}
	Checkers["SWithReadiness"] = &SWithReadiness{}
	Checkers["SWithEmptyDirSizeLimit"] = &SWithEmptyDirSizeLimit{}
	Checkers["SWithLivenessReadinessDelayCheck"] = &SWithLivenessReadinessDelayCheck{}
	Checkers["SWithTerminationGrace"] = &SWithTerminationGrace{}

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
		&WithRollingUpdate{},
		&DWithEmptyDirSizeLimit{},
	},
}

var statefulsetCheckSet = &Rule{
	Name: "statefulSet",
	Checkers: []Checker{
		&SWithLivenessCheck{},
		&SWithResourceRequestAndLimit{},
		&SWithReadiness{},
		&SWithEmptyDirSizeLimit{},
	},
}

// initialize the default rule set
var ruleSet = []*Rule{
	spotCheckSet,
	deploymentCheckSet,
	statefulsetCheckSet,
}
