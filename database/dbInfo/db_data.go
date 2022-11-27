package dbInfo

import "log"

type Data struct {
	Time    string
	Tester  string
	Name    string
	Board   string
	Model   string
	Version string
	LogPath string
	Result  string
}

func LogPrint(logString string) {
	log.Printf("\t    %v\n", logString)
}
