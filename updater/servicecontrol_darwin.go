// +build darwin

package updater

// serviceRunning checks to see if a service is installed/registered
func IsServiceRunning(serviceName string) (bool, error) {
	return true, nil
}

// StartService starts a service
func StartService(serviceName string) error {
	return nil
}

// StopService stops a service
func StopService(serviceName string) error {
	return nil
}
