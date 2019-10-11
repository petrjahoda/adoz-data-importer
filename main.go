package main

import (
	"time"
)

const version = "2019.4.1.10"
const deleteLogsAfter = 240 * time.Hour

func main() {
	LogDirectoryFileCheck("MAIN")
	LogInfo("MAIN", "Program version "+version+" started")
	CreateConfigIfNotExists()
	LoadSettingsFromConfigFile()
	LogDebug("MAIN", "Using ["+DatabaseType+"] on "+DatabaseIpAddress+":"+DatabasePort+" with database "+DatabaseName)
	SendMail("Program started", "Zapsi Service version "+version+" started")
	for {
		LogInfo("MAIN", "Program running")

		time.Sleep(10 * time.Second)
	}
}
