// +build windows

package updater

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// serviceRunning checks to see if a service is installed/registered
func IsServiceRunning(serviceName string) (bool, error) {
	// open service manager, requires admin
	m, err := mgr.Connect()
	if nil != err {
		return false, err
	}
	defer m.Disconnect()

	// open the service
	s, err := m.OpenService(serviceName)
	if nil != err {
		return false, err
	}
	defer s.Close()

	// Interrogate service
	status, err := s.Control(svc.Interrogate)
	if nil != err {
		// will return an error if the service is not running so just return false
		return false, nil
	}

	if status.State != svc.Running {
		return false, nil
	}

	return true, nil
}

// StartService starts a service
func StartService(serviceName string) error {
	// open service manager, requires admin
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	// open the service
	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()

	// start the service
	err = s.Start()
	if err != nil {
		return err
	}

	return nil
}

// StopService stops a service
func StopService(serviceName string) error {
	// logger.Debug(fmt.Sprintf("Stopping the '%s' service", serviceName))

	// open service manager, requires admin
	m, err := mgr.Connect()
	if nil != err {
		return err
	}
	defer m.Disconnect()

	// open the service
	s, err := m.OpenService(serviceName)
	if nil != err {
		return err
	}
	defer s.Close()

	// stop the service
	_, err = s.Control(svc.Stop)
	if nil != err {
		return err
	}
	// allow time to stop
	time.Sleep(5 * time.Second)

	status, err := s.Query()
	if nil != err {
		// will return an error if the service is not running so just return err
		return err
	}

	if status.State != svc.Stopped {
		err = fmt.Errorf("'%s' did not stop; status: %+v", serviceName, status)
		return err
	}

	return nil
}
