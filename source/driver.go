package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/Dartmouth-OpenAV/microservice-framework/framework"
)

func convertAndSend(socketKey string, cmdString string) (bool){
	sent := framework.WriteLineToSocket(socketKey, cmdString)

	if !sent {
		errMsg := "Error sending command"
		framework.AddToErrors(socketKey, errMsg)
	}

	return sent
}

func readAndConvert(socketKey string) (string, error){
	resp := framework.ReadLineFromSocket(socketKey)

	if resp == "" {
		errMsg := "response was blank"
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	return resp, nil
}

func setVideoMute(socketKey string, output string, state string, muteType string) (string, error){
	function := "setVideoMute"

	value := "notok"
	err := error(nil)
	maxRetries := 2

	for maxRetries > 0 {
		value, err = setVideoMuteDo(socketKey, output, state, muteType)
		if value != "ok" { // Something went wrong - perhaps try again
			framework.Log(function + " - j4do2md retrying videomute operation")
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func setVideoMuteDo(socketKey string, output string, state string, muteType string) (string, error){
	function := "setVideoMuteDo"

	var muteCmd string

	framework.Log("State: " + state)
	framework.Log("MuteType: " + muteType)

	if state == `"false"`{
		muteCmd = "0"
	} else if state == `"true"`{
		if muteType == "video" {
			muteCmd = "1"
		} else if muteType == "videosync"{
			muteCmd = "2"
		}
	} else if state == `"toggle"`{
		//Getting the current state from Extron device.
		//Not using the get function since it needs an immediate response.
		cmdString := "B"
		sent := convertAndSend(socketKey, cmdString)
		if !sent {
			errMsg := fmt.Sprintf(function + "FSD342 - error sending command")
			framework.AddToErrors(socketKey, errMsg)
			return errMsg, errors.New(errMsg)
		}

		resp, err := readAndConvert(socketKey)

		if err != nil{
			errMsg := fmt.Sprintf(function + "l2ei4n - Error reading from socket")
			framework.AddToErrors(socketKey, errMsg)
			return errMsg, errors.New(errMsg)
		}

		framework.Log("Resp: " + resp)
		resp2 := strings.Split(resp, " ")
		framework.Log("Resp2: " + fmt.Sprint(resp2))

		outputInt, err := strconv.Atoi(output)

		if err != nil{
			errMsg := fmt.Sprintf(function + "gjwe4h - error converting from string to int")
			framework.AddToErrors(socketKey, errMsg)
			return errMsg, errors.New(errMsg)
		}

		currentSetting := resp2[outputInt - 1]
		framework.Log("Current Setting: " + currentSetting)

		if currentSetting == "0"{
			if muteType == "video" {
				muteCmd = "1"
			} else if muteType == "videosync"{
				muteCmd = "2"
			}
		} else if currentSetting == "1" || currentSetting == "2"{
			muteCmd = "0"
		}
	}else {
		errMsg := fmt.Sprintf(function + " - unrecognized state value: " + state)
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	output = strings.Trim(output, "\"")

	cmdString := output + "*" + muteCmd + "B"
	framework.Log("Sending videomute cmdString: " + cmdString)

	sent := convertAndSend(socketKey, cmdString)

	if !sent {
		errMsg := fmt.Sprintf(function+" - mgk4kd error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp, err := readAndConvert(socketKey)
	
	if err != nil{
		errMsg := fmt.Sprintf(function + "4jndj4 - error reading from socket")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	framework.Log(function + " - Decoded Resp: " + resp + "\n")

	return "ok", nil
}

func getVideoMute(socketKey string, output string, muteType string) (string, error){
	function := "getVideoMute"

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = getVideoMuteDo(socketKey, output, muteType)
		if value == `"unknown"` { // Something went wrong - perhaps try again
			framework.Log(function + " - g39dk2 retrying video mute operation")
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func getVideoMuteDo(socketKey string, output string, muteType string) (string, error){
	function := "getVideoMuteDo"
	cmdString := "B"
	sent := convertAndSend(socketKey, cmdString)

	if !sent {
		errMsg := fmt.Sprintf(function + " - fk4jdk - error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp, err := readAndConvert(socketKey)

	if err != nil{
		errMsg := fmt.Sprintf(function + " - 3ido3d - error reading from socket")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp2 := strings.Split(resp, " ")
	framework.Log(fmt.Sprintf("%s", resp2))

	outputInt, err := strconv.Atoi(output)

	msg:= ""
	if resp2[outputInt-1] == "0"{
		msg = "false"
	}else if resp2[outputInt-1] == "1"{
		if muteType == "video"{
			msg = "true"
		}else{
			msg = "false"
		}
	}else if resp2[outputInt-1] == "2"{
		if muteType == "videosync"{
			msg = "true"
		}else{
			msg = "false"
		}
	}else{
		errMsg := fmt.Sprintf(function + " - n3jdn4 - Can't interpret videomute response")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	// for index := range resp2 {
	// 	if resp2[index] == "0" {
	// 		msg = msg + "Input " + fmt.Sprint(index + 1) + ": Video mute disabled. "
	// 	}else if resp2[index] == "1" {
	// 		msg = msg + "Input " + fmt.Sprint(index + 1) + ": Video mute enabled. "
	// 	}else if resp2[index] == "2" {
	// 		msg = msg + "Input " + fmt.Sprint(index + 1) + ": Video and sync mute enabled. "
	// 	}
	// }

	// msg = strings.Trim(msg, " ")

	// If we got here, the response was good, so successful return with the state indication
	framework.Log("Response for output " + fmt.Sprint(outputInt) + ":")
	framework.Log("Mute type: " + muteType)
	framework.Log(msg)
	return `"` + fmt.Sprint(msg) + `"`, nil
}

func setRawCommand(socketKey string, cmd string) (string, error) {
	function := "setRawCommand"

	value := `"unknown"`
	err := error(nil)
	maxRetries := 2
	for maxRetries > 0 {
		value, err = setRawCommandDo(socketKey, cmd)
		if value == `"unknown"` { // Something went wrong - perhaps try again
			framework.Log(function + " - h3kdj4 retrying raw command operation")
			maxRetries--
			time.Sleep(1 * time.Second)
		} else { // Succeeded
			maxRetries = 0
		}
	}

	return value, err
}

func setRawCommandDo(socketKey string, cmd string) (string, error){
	function := "setRawCommandDo"

	cmdString := strings.Trim(cmd, "\"")

	sent := convertAndSend(socketKey, cmdString)

	if !sent {
		errMsg := fmt.Sprintf(function+" - mgk4kd error sending command")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	resp, err := readAndConvert(socketKey)
	if err != nil{
		errMsg := fmt.Sprintf(function + "3kdkd4 - error reading from socket")
		framework.AddToErrors(socketKey, errMsg)
		return errMsg, errors.New(errMsg)
	}

	framework.Log(function + " - decoded resp: " + resp + "\n")

	if len(resp) > 0{
		return resp, nil
	} else {
		framework.Log(function + ": No response")
		return "ok", nil
	}
}