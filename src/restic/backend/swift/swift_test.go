package swift_test

import (
	"fmt"
	"math/rand"
	"os"
	"restic"
	"time"

	"restic/debug"
	"restic/errors"

	"restic/backend/swift"
	"restic/backend/test"
)

//go:generate go run ../test/generate_backend_tests.go

func init() {
	// only use minimal data
	test.MinimalData = true

	cfg := swift.Config{}

	for _, val := range []struct {
		s   *string
		env string
	}{
		// v2/v3 specific
		{&cfg.UserName, "OS_USERNAME"},
		{&cfg.APIKey, "OS_PASSWORD"},
		{&cfg.Region, "OS_REGION_NAME"},
		{&cfg.AuthURL, "OS_AUTH_URL"},

		// v3 specific
		{&cfg.Domain, "OS_USER_DOMAIN_NAME"},
		{&cfg.Tenant, "OS_PROJECT_NAME"},
		{&cfg.TenantDomain, "OS_PROJECT_DOMAIN_NAME"},

		// v2 specific
		{&cfg.TenantID, "OS_TENANT_ID"},
		{&cfg.Tenant, "OS_TENANT_NAME"},

		// v1 specific
		{&cfg.AuthURL, "ST_AUTH"},
		{&cfg.UserName, "ST_USER"},
		{&cfg.APIKey, "ST_KEY"},

		// Manual authentication
		{&cfg.StorageURL, "OS_STORAGE_URL"},
		{&cfg.AuthToken, "OS_AUTH_TOKEN"},

		{&cfg.DefaultContainerPolicy, "SWIFT_DEFAULT_CONTAINER_POLICY"},
	} {
		if *val.s == "" {
			*val.s = os.Getenv(val.env)
		}
	}

	cfg.Container = os.Getenv("RESTIC_TEST_SWIFT_CONTAINER")
	if cfg.Container == "" {
		SkipMessage = "RESTIC_TEST_SWIFT_CONTAINER unset, skipping test"
		return
	}

	// use a unique prefix
	rand.Seed(time.Now().UnixNano())
	cfg.Prefix = fmt.Sprintf("travis-%s-%d", os.Getenv("TRAVIS_BUILD_ID"), rand.Int63())

	debug.Log("opening swift repository at %#v", cfg)

	test.CreateFn = func() (restic.Backend, error) {
		be, err := swift.Open(cfg)
		if err != nil {
			return nil, err
		}

		exists, err := be.Test(restic.Handle{Type: restic.ConfigFile})
		if err != nil {
			return nil, err
		}

		if exists {
			return nil, errors.New("config already exists")
		}

		return be, nil
	}

	test.OpenFn = func() (restic.Backend, error) {
		return swift.Open(cfg)
	}

	test.CleanupFn = func() error {
		type deleter interface {
			Delete() error
		}

		be, err := swift.Open(cfg)
		if err != nil {
			return err
		}

		return be.(deleter).Delete()
	}
}
