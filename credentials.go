package freetds

import (
	"strconv"
	"strings"
)

const defaultConnTimeoutSec = 10

type credentials struct {
	user, pwd, host, database, mirrorHost, compatibility, appName, tdsVersion string
	maxPoolSize, lockTimeout, connTimeout                                     int
}

// NewCredentials fills credentials stusct from connection string
func NewCredentials(connStr string) *credentials {
	parts := strings.Split(connStr, ";")
	crd := &credentials{maxPoolSize: 100, connTimeout: defaultConnTimeoutSec}
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			key := strings.ToLower(strings.Trim(kv[0], " "))
			value := kv[1]
			switch key {
			case "server", "host":
				crd.host = value
			case "database":
				crd.database = value
			case "user id", "user_id", "user":
				crd.user = value
			case "password", "pwd":
				crd.pwd = value
			case "failover partner", "failover_partner", "mirror", "mirror_host", "mirror host":
				crd.mirrorHost = value
			case "app", "application", "application name", "application_name":
				crd.appName = value
			case "max pool size", "max_pool_size":
				if i, err := strconv.Atoi(value); err == nil {
					crd.maxPoolSize = i
				}
			case "compatibility_mode", "compatibility mode", "compatibility":
				crd.compatibility = strings.ToLower(value)
			case "lock timeout", "lock_timeout":
				if i, err := strconv.Atoi(value); err == nil {
					crd.lockTimeout = i
				}
			case "conn timeout", "conn_timeout", "connection timeout", "connection_timeout":
				if i, err := strconv.Atoi(value); err == nil {
					crd.connTimeout = i
				}
			case "tds version", "tds_version":
				crd.tdsVersion = value
			}
		}
	}
	return crd
}
