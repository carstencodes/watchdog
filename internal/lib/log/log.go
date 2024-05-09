package log

func CreateLog(minLevel Level, setup Setup) (Log, error) {
	levels := getLogLevel(minLevel)
	writer, err := setup.Build()
	if err != nil {
		return nil, err
	}

	return newLogShell(levels, writer), nil
}

func getLogLevel(minLevel Level) []Level {
	var levels []Level
	switch minLevel {
	case Debug:
		levels = append(levels, Debug)
		fallthrough
	case Info:
		levels = append(levels, Info)
		fallthrough
	case Warning:
		levels = append(levels, Warning)
		fallthrough
	case Error:
		levels = append(levels, Error)
		fallthrough
	case Fatal:
		levels = append(levels, Fatal)
	}
	return levels
}
