package commander

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (bt *BuildTool) handleSystemctl(params interface{}) error {
	return bt.executeCommand("systemctl", params)
}

func (bt *BuildTool) handleSystemctlFile(params interface{}) error {
	paramMap, ok := params.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid params for systemctl-file")
	}

	name, _ := paramMap["name"].(string)
	description, _ := paramMap["description"].(string)
	execStart, _ := paramMap["ExecStart"].(string)
	restart, _ := paramMap["Restart"].(string)
	tmpPath, _ := paramMap["tmp"].(string)
	location, _ := paramMap["location"].(string)

	// Create a new SystemctlService instance
	service := NewSystemctlService(name, description, execStart, restart)

	// Generate the service file
	serviceFilePath, err := service.GenerateServiceFile(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to generate systemctl service file: %v", err)
	}

	// Move the service file to the specified location
	err = MoveServiceFile(serviceFilePath, location)
	if err != nil {
		return fmt.Errorf("failed to move systemctl service file: %v", err)
	}

	fmt.Printf("Systemctl service file created and moved to: %s\n", filepath.Join(location, filepath.Base(serviceFilePath)))

	return nil
}

// SystemctlService represents a systemd service file.
type SystemctlService struct {
	Name        string
	Description string
	ExecStart   string
	Restart     string
}

// NewSystemctlService creates a new SystemctlService instance.
func NewSystemctlService(name, description, execStart, restart string) *SystemctlService {
	return &SystemctlService{
		Name:        name,
		Description: description,
		ExecStart:   execStart,
		Restart:     restart,
	}
}

// GenerateServiceFile generates a systemd service file based on the SystemctlService configuration.
func (s *SystemctlService) GenerateServiceFile(tmpPath string) (string, error) {
	content := fmt.Sprintf(`[Unit]
Description=%s

[Service]
ExecStart=%s
Restart=%s

[Install]
WantedBy=multi-user.target`, s.Description, s.ExecStart, s.Restart)

	serviceFilePath := filepath.Join(tmpPath, fmt.Sprintf("%s.service", s.Name))
	err := ioutil.WriteFile(serviceFilePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write service file: %w", err)
	}

	return serviceFilePath, nil
}

// MoveServiceFile moves the generated systemd service file to a specified location.
func MoveServiceFile(serviceFilePath, destinationPath string) error {
	destination := filepath.Join(destinationPath, filepath.Base(serviceFilePath))
	err := os.Rename(serviceFilePath, destination)
	if err != nil {
		return fmt.Errorf("failed to move service file: %w", err)
	}
	return nil
}