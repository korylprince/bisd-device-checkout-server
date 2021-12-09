package main

import (
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
	auth "github.com/korylprince/go-ad-auth/v3"
)

//Config represents options given in the environment
type Config struct {
	SessionExpiration int //in minutes; default: 60

	LDAPServer   string //required
	LDAPPort     int    //default: 389
	LDAPBaseDN   string //required
	LDAPGroup    string //optional
	LDAPSecurity string //default: none
	ldapSecurity auth.SecurityType

	SQLDriver    string //required
	InventoryDSN string //required
	SkywardDSN   string //required

	ListenAddr string //addr format used for net.Dial; required
	Prefix     string //url prefix to mount api to without trailing slash
}

var config = &Config{}

func checkEmpty(val, name string) {
	if val == "" {
		log.Fatalf("INVENTORY_%s must be configured\n", name)
	}
}

func init() {
	err := envconfig.Process("INVENTORY", config)
	if err != nil {
		log.Fatalln("Error reading configuration from environment:", err)
	}

	if config.SessionExpiration == 0 {
		config.SessionExpiration = 60
	}
	checkEmpty(config.LDAPServer, "LDAPSERVER")

	if config.LDAPPort == 0 {
		config.LDAPPort = 389
	}

	checkEmpty(config.LDAPBaseDN, "LDAPBASEDN")

	switch strings.ToLower(config.LDAPSecurity) {
	case "", "none":
		config.ldapSecurity = auth.SecurityNone
	case "tls":
		config.ldapSecurity = auth.SecurityTLS
	case "starttls":
		config.ldapSecurity = auth.SecurityStartTLS
	default:
		log.Fatalln("Invalid INVENTORY_LDAPSECURITY:", config.LDAPSecurity)
	}

	checkEmpty(config.SQLDriver, "SQLDRIVER")
	checkEmpty(config.InventoryDSN, "INVENTORYDSN")
	checkEmpty(config.SkywardDSN, "SKYWARDDSN")

	if config.SQLDriver == "mysql" && !strings.Contains(config.InventoryDSN, "?parseTime=true") {
		log.Fatalln("mysql DSN must contain \"?parseTime=true\"")
	}

	checkEmpty(config.ListenAddr, "LISTENADDR")
}
