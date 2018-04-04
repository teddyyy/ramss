package systemd

import (
	"fmt"
	"os"

	"../model"
	"github.com/godbus/dbus"
)

const destBus = "org.freedesktop.systemd1"
const objectPath = "/org/freedesktop/systemd1"
const getMethod = "org.freedesktop.DBus.Properties.Get"
const mngerMethod = "org.freedesktop.systemd1.Manager"
const destUnit = "org.freedesktop.systemd1.Unit"
const destService = "org.freedesktop.systemd1.Service"

// Get ...
func Get(ID string) *model.Unit {
	var dstService = ID + ".service"

	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to system bus:", err)
		return nil
	}

	var path dbus.ObjectPath
	err = conn.Object(destBus, objectPath).Call(mngerMethod+".GetUnit", 0, dstService).Store(&path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get object path:", err)
		return nil
	}

	var desc string
	err = conn.Object(destBus, path).Call(getMethod, 0, destUnit, "Description").Store(&desc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get description:", err)
		return nil
	}

	var loadStatus string
	err = conn.Object(destBus, path).Call(getMethod, 0, destUnit, "LoadState").Store(&loadStatus)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get load status:", err)
		return nil
	}

	var activeStatus string
	err = conn.Object(destBus, path).Call(getMethod, 0, destUnit, "ActiveState").Store(&activeStatus)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get active status:", err)
		return nil
	}

	var unitFileState string
	err = conn.Object(destBus, path).Call(getMethod, 0, destUnit, "UnitFileState").Store(&unitFileState)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get unit file status:", err)
		return nil
	}

	var pid int
	err = conn.Object(destBus, path).Call(getMethod, 0, destService, "MainPID").Store(&pid)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get main pid:", err)
		return nil
	}

	u := &model.Unit{
		ID:            ID,
		Description:   desc,
		LoadState:     loadStatus,
		ActiveState:   activeStatus,
		UnitFileState: unitFileState,
		MainPID:       pid,
	}

	return u
}

// Post ...
func Post(ID string, action string, mode string) error {
	var dstService = ID + ".service"

	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to system bus:", err)
		return nil
	}

	var path dbus.ObjectPath

	obj := conn.Object(destBus, objectPath)

	switch action {
	case "start":
		err = obj.Call(mngerMethod+".StartUnit", 0, dstService, mode).Store(&path)
	case "restart":
		err = obj.Call(mngerMethod+".ReStartUnit", 0, dstService, mode).Store(&path)
	case "stop":
		err = obj.Call(mngerMethod+".StopUnit", 0, dstService, mode).Store(&path)
	case "reload":
		err = obj.Call(mngerMethod+".ReloadUnit", 0, dstService, mode).Store(&path)
	default:
		fmt.Fprintln(os.Stderr, "Unknown action:", err)
		return nil
	}

	return nil
}
