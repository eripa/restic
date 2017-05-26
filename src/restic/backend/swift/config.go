package swift

import (
	"os"
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

// ApplyEnvironment saves values from the environment to the config.
func ApplyEnvironment(cfg interface{}) error {
	c := cfg.(*Config)
	for _, val := range []struct {
		s   *string
		env string
	}{
		// v2/v3 specific
		{&c.UserName, "OS_USERNAME"},
		{&c.APIKey, "OS_PASSWORD"},
		{&c.Region, "OS_REGION_NAME"},
		{&c.AuthURL, "OS_AUTH_URL"},

		// v3 specific
		{&c.Domain, "OS_USER_DOMAIN_NAME"},
		{&c.Tenant, "OS_PROJECT_NAME"},
		{&c.TenantDomain, "OS_PROJECT_DOMAIN_NAME"},

		// v2 specific
		{&c.TenantID, "OS_TENANT_ID"},
		{&c.Tenant, "OS_TENANT_NAME"},

		// v1 specific
		{&c.AuthURL, "ST_AUTH"},
		{&c.UserName, "ST_USER"},
		{&c.APIKey, "ST_KEY"},

		// Manual authentication
		{&c.StorageURL, "OS_STORAGE_URL"},
		{&c.AuthToken, "OS_AUTH_TOKEN"},

		{&c.DefaultContainerPolicy, "SWIFT_DEFAULT_CONTAINER_POLICY"},
	} {
		if *val.s == "" {
			*val.s = os.Getenv(val.env)
		}
	}
	return nil
}
