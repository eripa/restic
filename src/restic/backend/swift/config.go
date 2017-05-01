package swift

import (
	"restic/errors"
	"strings"
)

// Config contains basic configuration needed to specify swift location for a swift server
type Config struct {
	UserName     string
	Domain       string
	APIKey       string
	AuthURL      string
	Region       string
	Tenant       string
	TenantID     string
	TenantDomain string
	TrustID      string

	StorageURL string
	AuthToken  string

	Container              string
	Prefix                 string
	DefaultContainerPolicy string
}

// ParseConfig parses the string s and extract swift's container name and prefix.
func ParseConfig(s string) (interface{}, error) {
	data := strings.SplitN(s, ":", 3)
	if len(data) != 3 {
		return nil, errors.New("invalid URL, expected: swift:container-name:/[prefix]")
	}

	scheme, container, prefix := data[0], data[1], data[2]
	if scheme != "swift" {
		return nil, errors.Errorf("unexpected prefix: %s", data[0])
	}

	if len(prefix) == 0 {
		return nil, errors.Errorf("prefix is empty")
	}

	if prefix[0] != '/' {
		return nil, errors.Errorf("prefix does not start with slash (/)")
	}
	prefix = prefix[1:]

	cfg := Config{
		Container: container,
		Prefix:    prefix,
	}

	return cfg, nil
}
