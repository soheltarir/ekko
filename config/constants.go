package config

// ConsumerStatus is the enumeration for defining the running status of the Ping consumer
type ConsumerStatus string

const (
	NotStarted ConsumerStatus = "Not Started"
	Running                   = "Running"
	Stopped                   = "Stopped"
)
