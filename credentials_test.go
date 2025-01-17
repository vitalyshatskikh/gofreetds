package freetds

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConnectionString(t *testing.T) {
	validConnStrings := []string{
		"Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=myPassword;Failover Partner=myMirror;Max Pool Size=200;Lock Timeout=1000",
		"Server=myServerAddress;Database=myDataBase;User_Id=myUsername;Password=myPassword;Failover_Partner=myMirror;Max_Pool_Size=200;Lock_Timeout=1000",
		"server=myServerAddress;database=myDataBase;user_id=myUsername;password=myPassword;failover_partner=myMirror;max_pool_size=200;lock_timeout=1000",
		"host=myServerAddress;database=myDataBase;user=myUsername;pwd=myPassword;mirror=myMirror;max_pool_size=200;lock_timeout=1000",
		"host=myServerAddress;database=myDataBase;user=myUsername;pwd=myPassword;mirror=myMirror;max_pool_size=200;lock_timeout=1000",
	}
	for _, connStr := range validConnStrings {
		testCredentials(t, NewCredentials(connStr))
	}
}

func testCredentials(t *testing.T, crd *credentials) {
	assert.NotNil(t, crd)
	assert.Equal(t, "myServerAddress", crd.host)
	assert.Equal(t, "myDataBase", crd.database)
	assert.Equal(t, "myUsername", crd.user)
	assert.Equal(t, "myPassword", crd.pwd)
	assert.Equal(t, "myMirror", crd.mirrorHost)
	assert.Equal(t, 200, crd.maxPoolSize)
	assert.Equal(t, 1000, crd.lockTimeout)
}

func TestParseConnectionStringCompatibilityMode(t *testing.T) {
	setDefaultStrings := map[string]string{
		"Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=myPassword;Failover Partner=myMirror;Max Pool Size=200;Lock Timeout=1000":                                "",
		"Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=myPassword;Failover Partner=myMirror;Max Pool Size=200;Lock Timeout=1000;compatibility_mode=Sybase":      "sybase",
		"Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=myPassword;Failover Partner=myMirror;Max Pool Size=200;Lock Timeout=1000;compatibility mode=sybase":      "sybase",
		"Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=myPassword;Failover Partner=myMirror;Max Pool Size=200;Lock Timeout=1000;Compatibility Mode=sybase":      "sybase",
		"Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=myPassword;Failover Partner=myMirror;Max Pool Size=200;Lock Timeout=1000;Compatibility_Mode=Sybase":      "sybase",
		"Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=myPassword;Failover Partner=myMirror;Max Pool Size=200;Lock Timeout=1000;Compatibility Mode=sybase_12_5": "sybase_12_5",
		"Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=myPassword;Failover Partner=myMirror;Max Pool Size=200;Lock Timeout=1000;Compatibility_Mode=Sybase_12_5": "sybase_12_5",
		"Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=myPassword;Failover Partner=myMirror;Max Pool Size=200;Lock Timeout=1000;Compatibility=Other":            "other",
		"Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=myPassword;Failover Partner=myMirror;Max Pool Size=200;Lock Timeout=1000;compatibility=other":            "other",
	}
	for connStr, expected := range setDefaultStrings {
		crd := NewCredentials(connStr)
		assert.NotNil(t, crd)
		assert.Equal(t, "myServerAddress", crd.host)
		assert.Equal(t, "myDataBase", crd.database)
		assert.Equal(t, "myUsername", crd.user)
		assert.Equal(t, "myPassword", crd.pwd)
		assert.Equal(t, "myMirror", crd.mirrorHost)
		assert.Equal(t, 200, crd.maxPoolSize)
		assert.Equal(t, 1000, crd.lockTimeout)
		assert.Equal(t, expected, crd.compatibility)
	}
}

// TestParseConnectionStringEqualsInValue tests parsing when for e.g. password contains "="
func TestParseConnectionStringEqualsInValue(t *testing.T) {
	validConnStrings := []string{
		"Server=myServerAddress;Database=myDataBase;User Id=myUsername;Password=my=Password;Failover Partner=myMirror;Max Pool Size=200;Lock Timeout=1000",
		"Server=myServerAddress;Database=myDataBase;User_Id=myUsername;Password=my=Password;Failover_Partner=myMirror;Max_Pool_Size=200;Lock_Timeout=1000",
		"server=myServerAddress;database=myDataBase;user_id=myUsername;password=my=Password;failover_partner=myMirror;max_pool_size=200;lock_timeout=1000",
		"host=myServerAddress;database=myDataBase;user=myUsername;pwd=my=Password;mirror=myMirror;max_pool_size=200;lock_timeout=1000",
		"host=myServerAddress;database=myDataBase;user=myUsername;pwd=my=Password;mirror=myMirror;max_pool_size=200;lock_timeout=1000",
	}
	for _, connStr := range validConnStrings {
		crd := NewCredentials(connStr)
		assert.NotNil(t, crd)
		assert.Equal(t, "myServerAddress", crd.host)
		assert.Equal(t, "myDataBase", crd.database)
		assert.Equal(t, "myUsername", crd.user)
		assert.Equal(t, "my=Password", crd.pwd)
		assert.Equal(t, "myMirror", crd.mirrorHost)
		assert.Equal(t, 200, crd.maxPoolSize)
		assert.Equal(t, 1000, crd.lockTimeout)
	}
}

func TestParseConnectionStringApplicationName(t *testing.T) {
	validConnStrings := []string{
		"host=myServerAddress;database=myDataBase;user=myUsername;pwd=myPassword;mirror=myMirror;app=myAppName",
		"host=myServerAddress;database=myDataBase;user=myUsername;pwd=myPassword;mirror=myMirror;application=myAppName",
		"host=myServerAddress;database=myDataBase;user=myUsername;pwd=myPassword;mirror=myMirror;application name=myAppName",
		"host=myServerAddress;database=myDataBase;user=myUsername;pwd=myPassword;mirror=myMirror;application_name=myAppName",
	}
	expectedAppName := "myAppName"

	for _, connStr := range validConnStrings {
		credentials := NewCredentials(connStr)
		assert.Equal(t, expectedAppName, credentials.appName)
	}
}

func TestParseConnectionStringConnTimeout(t *testing.T) {
	expectedTimeout := defaultConnTimeoutSec + 5
	validConnStrings := []string{
		fmt.Sprintf("host=myhost;database=mydb;user=myuser;pwd=mypwd;conn timeout=%d", expectedTimeout),
		fmt.Sprintf("host=myhost;database=mydb;user=myuser;pwd=mypwd;conn_timeout=%d", expectedTimeout),
		fmt.Sprintf("host=myhost;database=mydb;user=myuser;pwd=mypwd;connection timeout=%d", expectedTimeout),
		fmt.Sprintf("host=myhost;database=mydb;user=myuser;pwd=mypwd;connection_timeout=%d", expectedTimeout),
	}

	for _, connStr := range validConnStrings {
		credentials := NewCredentials(connStr)
		assert.Equal(t, expectedTimeout, credentials.connTimeout)
	}
}

func TestParseConnectionStringDefaultConnTimeout(t *testing.T) {
	expectedTimeout := defaultConnTimeoutSec
	validConnString := fmt.Sprintf("host=myhost;database=mydb;user=myuser;pwd=mypwd;conn timeout=%d", expectedTimeout)

	credentials := NewCredentials(validConnString)
	assert.Equal(t, expectedTimeout, credentials.connTimeout)
}

func TestParseConnectionStringTDSVersion(t *testing.T) {
	expectedVersion := "myVersion"
	validConnStrings := []string{
		fmt.Sprintf("host=myhost;database=mydb;user=myuser;pwd=mypwd;tds_version=%s", expectedVersion),
		fmt.Sprintf("host=myhost;database=mydb;user=myuser;pwd=mypwd;tds version=%s", expectedVersion),
	}

	for _, connStr := range validConnStrings {
		credentials := NewCredentials(connStr)
		assert.Equal(t, expectedVersion, credentials.tdsVersion)
	}
}
