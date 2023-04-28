package cortex

import (
	"context"
	"errors"
	"fmt"
	"github.com/dghubble/sling"
	"net/url"
)

const (
	// MajorVersion is the major version
	MajorVersion = 0
	// MinorVersion is the minor version
	MinorVersion = 0
	// PatchVersion is the patch version
	PatchVersion = 1

	// UserAgentPrefix is the prefix of the User-Agent header that all terraform REST calls perform
	UserAgentPrefix = "cortex-terraform-provider"
)

// Version is the semver of this provider
var Version = fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)

type Client struct {
	ctx     context.Context
	client  *sling.Sling
	baseUrl string
	token   string
	version string
}

type OptionDelegator func(c *Client) error

// NewClient initializes a new API client for Cortex
func NewClient(opts ...OptionDelegator) (*Client, error) {
	c := &Client{}
	for _, f := range opts {
		if err := f(c); err != nil {
			return nil, err
		}
	}

	c.client = sling.New().Base(c.baseUrl).
		Set("User-Agent", fmt.Sprintf("%s (%s)", UserAgentPrefix, c.version)).
		Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return c, nil
}

func WithVersion(version string) func(*Client) error {
	return func(c *Client) error {
		if version == "" {
			return errors.New("cannot specify empty version")
		}
		c.version = version
		return nil
	}
}

// WithContext Specify the context for the cortex client to use.
func WithContext(ctx context.Context) func(*Client) error {
	return func(c *Client) error {
		c.ctx = ctx
		return nil
	}
}

// WithURL Specify the base URL for the cortex client to connect to.
func WithURL(baseUrl string) func(*Client) error {
	return func(c *Client) error {
		if baseUrl == "" {
			return errors.New("cannot specify empty API Base URL")
		}
		if _, err := url.Parse(baseUrl); err != nil {
			return err
		}
		c.baseUrl = baseUrl
		return nil
	}
}

// WithToken Specify the API token for the cortex client to use.
func WithToken(token string) func(*Client) error {
	return func(c *Client) error {
		if token == "" {
			return errors.New("cannot specify empty token")
		}
		c.token = token
		return nil
	}
}
