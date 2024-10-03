package main

import (
	"errors"
	"github.com/Dartmouth-OpenAV/microservice-framework/framework"
)

func setFrameworkGlobals() {
	framework.MicroserviceName = "OpenAV Global Cache Serial Extron Microservice"

	// Define device specific globals
	framework.DefaultSocketPort = 4999   // Global Cache's default socket port
	// var defaultLineDelimiter = 0xA // 0xA==\n==linefeed, 0xD==\r==carriagereturn

	// globals that change modes in the microservice framework:
	framework.CheckFunctionAppendBehavior = "Don't add another instance"
	framework.RegisterMainGetFunc(doDeviceSpecificGet)
	framework.RegisterMainSetFunc(doDeviceSpecificSet)

}

// Every microservice using this golang microservice framework needs to provide this function to invoke functions to do sets.
// socketKey is the network connection for the framework to use to communicate with the device.
// setting is the first parameter in the URI.
// arg1 are the second and third parameters in the URI.
//   Example PUT URIs that will result in this function being invoked:
// 	 ":address/:setting/"
//   ":address/:setting/:arg1"
//   ":address/:setting/:arg1/:arg2"
func doDeviceSpecificSet(socketKey string, setting string, arg1 string, arg2 string, arg3 string) (string, error) {
	function := "doDeviceSpecificSet"

	// Add a case statement for each set function your microservice implements.  These calls can use 0, 1, or 2 arguments.
	switch setting {
		case "videomute":
			return setVideoMute(socketKey, arg1, arg2, "video")
		case "videosyncmute":
			return setVideoMute(socketKey, arg1, arg2, "videosync")
		case "rawcommand":
			return setRawCommand(socketKey, arg1)
	}

	// If we get here, we didn't recognize the setting.  Send an error back to the config writer who had a bad URL.
	errMsg := function + " - unrecognized setting in URI: " + setting
	framework.AddToErrors(socketKey, errMsg)
	err := errors.New(errMsg)
	return setting, err
}

// Every microservice using this golang microservice framework needs to provide this function to invoke functions to do gets.
// socketKey is the network connection for the framework to use to communicate with the device.
// setting is the first parameter in the URI.
// arg1 are the second and third parameters in the URI.
//   Example GET URIs that will result in this function being invoked:
// 	 ":address/:setting/"
//   ":address/:setting/:arg1"
//   ":address/:setting/:arg1/:arg2"
// Every microservice using this golang microservice framework needs to provide this function to invoke functions to do gets.
func doDeviceSpecificGet(socketKey string, setting string, arg1 string, arg2 string) (string, error) {
	function := "doDeviceSpecificGet"

	switch setting {
	case "videomute":
		return getVideoMute(socketKey, arg1, "video")
	case "videosyncmute":
		return getVideoMute(socketKey, arg1, "videosync")
	//case "rawcommand":
		//return getRawCommand(socketKey, arg1)
	}

	// If we get here, we didn't recognize the setting.  Send an error back to the config writer who had a bad URL.
	errMsg := function + " - unrecognized setting in URI: " + setting
	framework.AddToErrors(socketKey, errMsg)
	err := errors.New(errMsg)
	return setting, err
}

func main() {
	setFrameworkGlobals()
	framework.Startup()
 } 
