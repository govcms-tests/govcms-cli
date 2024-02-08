package config

import "sort"

var DockerImages = []string{
	"adminer",
	"govcmsextras/dnsmasq",
	"mariadb:lts",
	"mailhog/mailhog",
	"nginxproxy/nginx-proxy",
}

func init() {
	sort.Strings(DockerImages)
}
