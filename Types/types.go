package Types

type Conf struct {
	Counter []DeviceCounters `yaml:"device_counters"`
}

type DeviceCounters struct {
	Name     string `yaml:"name"`
	Interval int    `yaml:"interval"`
	Path     string `yaml:"path"`
}

type ConfigRequest struct {
	DeviceIP   string `yaml:"device_ip"`
	DeviceName string `yaml:"device_name"`
	Protocol   string `yaml:"protocol"`
	Configs    []Conf `yaml:"configs"`
}

// For testing:

type Schema struct {
	Entries []SchemaEntry
}

type SchemaEntry struct {
	Name      string
	Tag       string
	Namespace string
	Value     string
}

type NamespaceParser struct {
	Parent              *NamespaceParser
	Children            []*NamespaceParser
	LastParentNamespace string
}
