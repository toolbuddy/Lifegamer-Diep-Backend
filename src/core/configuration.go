package core

import (
	"encoding/json"
	"os"
)

/**
 * ServerConfiguration:
 * The struct to present the server configuration.
 *
 * @property {string} Host 								- the host string of the server
 * @property {string} Port								- the port string of the server
 * @property {int} MaxRoom								- the max number of the room
 */
type ServerConfiguration struct {
	Host string
	Port string
	MaxRoom int
	MaxRoomMember int
}


/**
 * Configuration:
 * The struct to present all configuration of the project.
 *
 * @property {ServerConfiguration} Server - the configuration of the server
 */
type Configuration struct {
	Server ServerConfiguration
}


/**
 * <*Configuration>.loadFromFile:
 * The function in Configuration to load the json file.
 *
 * @return {error}
 */
func (c *Configuration) loadFromFile() error {
	file, _ := os.Open("src/config/main.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&c)
	if err != nil {
		return err
	}
	c.Server.MaxRoom = 50
	return nil
}
