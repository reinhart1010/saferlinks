package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3/log"
)

func getHostnameInitials(hostname string) string {
	splits := strings.SplitN(hostname, ".", 4)
	initials := ""
	for _, split := range splits {
		initials += string(split[0])
	}
	return initials
}

func resolveHostnameAlias(hostname string) string {
	initials := getHostnameInitials(hostname)

	rawFile, err := os.Open(fmt.Sprintf("./data/aliases/%s.csv", initials))
	if err != nil {
		log.Error(err)
		return hostname
	}

	aliasFile := csv.NewReader(rawFile)
	if aliasFile == nil {
		return hostname
	}

	for {
		record, err := aliasFile.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error(err)
			break
		}

		if record[0] == hostname {
			hostname = record[1]
			break
		}
	}
	return hostname
}
