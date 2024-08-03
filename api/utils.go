package api

import "strings"

func Hosts(commaHosts string) []string {
	hosts := []string{}

	hosts = append(hosts, strings.Split(commaHosts, ",")...)

	return hosts
}
