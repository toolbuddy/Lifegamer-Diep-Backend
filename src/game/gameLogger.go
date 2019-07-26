package game

import (
	"github.com/f26401004/Lifegamer-Diep-backend/src/util"
	"github.com/sirupsen/logrus"
	"os"
	"time"
	"strconv"
)

/**
 * GameLogger:
 * The struct of game logger.
 *
 * @property {*logrus.Logger}					 													- the game logger instance
 */
type GameLogger struct {
	instance *logrus.Logger
}

/**
 * <game>.NewLogger:
 * The function to new a game logger instance.
 *
 * @return {*logrus.Logger}
 */
func NewLogger (roomName string) *GameLogger {
	var baseLogger = logrus.New()
	// set output format to json
	baseLogger.SetFormatter(&logrus.JSONFormatter{})
	// default record the method
	baseLogger.SetReportCaller(true)
	// set the output stream
	file, err := os.OpenFile("./logs/" + roomName + "-" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".log", os.O_WRONLY | os.O_CREATE, 0755)
	if err != nil {
		return nil
	}
	baseLogger.SetOutput(file)
	var gameLogger = &GameLogger {
		instance: baseLogger,
	}
	// write the first log
	baseLogger.WithFields(logrus.Fields {
		"room-name": roomName,
	}).Info("The room being opened")
	return gameLogger
}

/**
 * <*GameLogger>.InvalidArg:
 * The function to record invalid argument message.
 *
 * @params {string} argumentName																- The argument name
 * 
 * @return {nil}
 */
func (l *GameLogger) InvalidArg (argumentName string) {
	l.instance.WithFields(logrus.Fields {
		"argmuent": argumentName,
	}).Error("Invalid argument")
}

/**
 * <*GameLogger>.InvalidArgValue:
 * The function to record invalid argument value message.
 *
 * @params {string} argumentName																- The argument name
 * @params {string} argumentValue																- The argument value
 * 
 * @return {nil}
 */
func (l *GameLogger) InvalidArgValue (argumentName, argumentValue string) {
	l.instance.WithFields(logrus.Fields {
		"argument": argumentName,
		"value": argumentValue,
	}).Error("Invalid argument value")
}

/**
 * <*GameLogger>.establishConnection:
 * The function to record establish connection success.
 *
 * @params {string} ip																					- The request ip
 * @params {string} playerName																	- The player name
 * @params {string} roomName																		- The room name that player joined
 * @params {int} roomMember																			- The member number of the room
 * 
 * @return {nil}
 */
func (l *GameLogger) establishConnection (ip, playerName, roomName string, roomMember int) {
	l.instance.WithFields(logrus.Fields {
		"request-ip": ip,
		"player-name": playerName,
		"room-name": roomName,
		"room-member": roomMember,
	}).Info("Establish connection success")
}

/**
 * <*GameLogger>.looseConnection:
 * The function to record loose connection.
 *
 * @params {string} playerName																	- The player name
 * @params {string} roomName																		- The room name that player joined
 * @params {int} roomMember																			- The member number of the room
 * 
 * @return {nil}
 */
func (l *GameLogger) looseConnection (playerName, roomName string, roomMember int) {
	l.instance.WithFields(logrus.Fields {
		"player-name": playerName,
		"room-name": roomName,
		"room-member": roomMember,
	}).Warn("Loose connection")
}

/**
 * <*GameLogger>.closeConnection:
 * The function to record closed connection.
 *
 * @params {string} playerName																	- The player name
 * @params {string} roomName																		- The room name that player joined
 * @params {int} roomMember																			- The member number of the room
 * 
 * @return {nil}
 */
func (l *GameLogger) closeConnection (playerName, roomName string, roomMember int) {
	l.instance.WithFields(logrus.Fields {
		"player-name": playerName,
		"room-name": roomName,
		"room-member": roomMember,
	}).Info("Close connection")
}

/**
 * <*GameLogger>.updatePlayerStatus:
 * The function to record update player session message.
 *
 * @params {string} playerName																	- The player name
 * @params {MoveDirection} playerMoving													- The move direction of the player
 * @params {float64} playerRotation															- The rotation of the player
 * 
 * @return {nil}
 */
func (l *GameLogger) updatePlayerStatus (playerName string, playerMoving util.MoveDirection, playerRotation float64) {
	l.instance.WithFields(logrus.Fields {
		"player-name": playerName,
		"player-moving": playerMoving,
		"player-rotation": playerRotation,
	}).Info("Update player status")
}

/**
 * <*GameLogger>.shootBullet:
 * The function to record shoot bullet message.
 *
 * @params {string} playerName																	- The player name
 * @params {int} number																					- The bullet number
 * @params {float64} angle																			- The angle pf the current shoot angle
 * 
 * @return {nil}
 */
func (l *GameLogger) shootBullet (playerName string, number int, angle float64) {
	l.instance.WithFields(logrus.Fields {
		"player-name": playerName,
		"angle": angle,
		"number": number,
	}).Info("Shoot bullet")
}

/**
 * <*GameLogger>.playerEvaluation:
 * The function to record shoot bullet message.
 *
 * @params {string} playerName																	- The player name
 * @params {string} attr																				- The target attribute to evaluate
 * @params {int} from																						- the origin attribute level
 * @params {int} to																							- the evaluated attribute level
 * 
 * @return {nil}
 */
 func (l *GameLogger) playerEvaluation (playerName, attr string, from, to int) {
	l.instance.WithFields(logrus.Fields {
		"player-name": playerName,
		"attribute": attr,
		"origin": from,
		"evaluated": to,
	}).Info("Shoot bullet")
}

/**
 * <*GameLogger>.deadMessage:
 * The function to record dead message.
 *
 * @params {string} target																			- The name of the killed
 * @params {string} killedBy																		- The name of the killer
 * 
 * @return {nil}
 */
func (l *GameLogger) deadMessage (target, killedBy string) {
	l.instance.WithFields(logrus.Fields {
		"target-name": target,
		"killed-by": killedBy,
	}).Info("Player dead")
}
