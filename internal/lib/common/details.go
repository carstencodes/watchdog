package common

var version string

const applicationName string = "watchdog"

const author string = "Carsten Igel"

const yearStart string = "2022"

var yearNow string

type ApplicationDetails interface {
	Name() string
	Version() string
	Author() string
	Copyright() string
}

type applicationDetails struct {
	name    string
	version string
	years   string
	author  string
}

func (a applicationDetails) Name() string {
	return a.name
}

func (a applicationDetails) Version() string {
	return a.version
}

func (a applicationDetails) Author() string {
	return a.author
}

func (a applicationDetails) Copyright() string {
	return "(C) " + a.years + " " + a.author
}

func ApplicationInfo() ApplicationDetails {
	return applicationDetails{
		name:    applicationName,
		version: version,
		years:   yearStart + "-" + yearNow,
		author:  author,
	}
}
