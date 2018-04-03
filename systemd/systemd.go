package systemd

import (
	"fmt"
	"os"

	"../model"
	"github.com/godbus/dbus"
)

// Get ...
func Get(ID string) *model.Unit {

	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to system bus:", err)
		return nil
	}

	var path dbus.ObjectPath
	err = conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1").Call("org.freedesktop.systemd1.Manager.GetUnit", 0, ID+".service").Store(&path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get object path:", err)
		return nil
	}

	var desc string
	err = conn.Object("org.freedesktop.systemd1", path).Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Unit", "Description").Store(&desc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get description:", err)
		return nil
	}

	var loadStatus string
	err = conn.Object("org.freedesktop.systemd1", path).Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Unit", "LoadState").Store(&loadStatus)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get load status:", err)
		return nil
	}

	var activeStatus string
	err = conn.Object("org.freedesktop.systemd1", path).Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Unit", "ActiveState").Store(&activeStatus)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get active status:", err)
		return nil
	}

	var unitFileState string
	err = conn.Object("org.freedesktop.systemd1", path).Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Unit", "UnitFileState").Store(&unitFileState)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get unit file status:", err)
		return nil
	}

	var pid int
	err = conn.Object("org.freedesktop.systemd1", path).Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.systemd1.Service", "MainPID").Store(&pid)
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
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to system bus:", err)
		return nil
	}

	var path dbus.ObjectPath

	obj := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")

	switch action {
	case "start":
		err = obj.Call("org.freedesktop.systemd1.Manager.StartUnit", 0, ID+".service", mode).Store(&path)
	case "restart":
		err = obj.Call("org.freedesktop.systemd1.Manager.ReStartUnit", 0, ID+".service", mode).Store(&path)
	case "stop":
		err = obj.Call("org.freedesktop.systemd1.Manager.StopUnit", 0, ID+".service", mode).Store(&path)
	case "reload":
		err = obj.Call("org.freedesktop.systemd1.Manager.ReloadUnit", 0, ID+".service", mode).Store(&path)
	default:
		fmt.Fprintln(os.Stderr, "Unknown action:", err)
		return nil
	}

	return nil
}
