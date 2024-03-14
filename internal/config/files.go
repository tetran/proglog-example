package config

import (
	"os"
	"path/filepath"
)

var (
	CAFile         = configFile("ca.pem")
	ServerCertFile = configFile("server.pem")
	ServerKeyFile  = configFile("server-key.pem")
	ClientCertFile = configFile("client.pem")
	ClientKeyFile  = configFile("client-key.pem")

	// OtherCAFile         = configFile("other-ca.pem")
	// OtherServerCertFile = configFile("other-server.pem")
	// OtherServerKeyFile  = configFile("other-server-key.pem")
	// OtherClientCertFile = configFile("other-client.pem")
	// OtherClientKeyFile  = configFile("other-client-key.pem")
)

func configFile(filename string) string {
	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, filename)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".proglog-example", filename)
}
