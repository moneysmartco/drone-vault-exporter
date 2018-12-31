package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/api"
)

// Error message
const (
	noSuchVaultKeyPath = "no such vault key path"
	nilResponse        = "got nil response"
	emptyClientToken   = "got empty client token"
	noSuchOutputFormat = "not support such output format"
)

type (
	// Config for the plugin.
	Config struct {
		VaultAddr       string
		VaultToken      string
		VaultAuthMethod string
		VaultRoleID     string
		VaultSecretID   string
		VaultKeyPath    string
		DeployEnvPath   string
		HelmEnvKey      string
		OutputFormat    string
	}

	// Plugin structure
	Plugin struct {
		Config Config
	}
)

func (p Plugin) useVaultAppRole(c *api.Client) error {
	resp, err := c.Logical().Write("auth/approle/login", map[string]interface{}{
		"role_id":   p.Config.VaultRoleID,
		"secret_id": p.Config.VaultSecretID,
	})

	if err != nil {
		return err
	}
	if resp == nil || resp.Auth == nil {
		return fmt.Errorf(nilResponse)
	}
	if resp.Auth.ClientToken == "" {
		return fmt.Errorf(emptyClientToken)
	}

	c.SetToken(resp.Auth.ClientToken)
	return nil
}

// Exec executes the plugin.
func (p Plugin) Exec() error {
	fmt.Println("================================")
	fmt.Println("= Here is drone-vault-exporter =")
	fmt.Println("================================")

	// Create vault client
	vaultConfig := api.DefaultConfig()
	vaultClient, err := api.NewClient(vaultConfig)
	if err != nil {
		return err
	}

	if p.Config.VaultAuthMethod == "APPROLE" {
		err := p.useVaultAppRole(vaultClient)
		if err != nil {
			return err
		}
	}

	// Get secret from vault
	secret, err := vaultClient.Logical().Read(p.Config.VaultKeyPath)
	if err != nil {
		return err
	}
	if secret == nil && err == nil {
		return fmt.Errorf(noSuchVaultKeyPath)
	}

	// Sort
	secretKeys := make([]string, 0)
	for k := range secret.Data {
		secretKeys = append(secretKeys, k)
	}
	sort.Strings(secretKeys)

	fmt.Printf("%d items in '%s'\n", len(secret.Data), p.Config.VaultKeyPath)
	for _, k := range secretKeys {
		fmt.Printf("  - %s\n", k)
	}

	// Write to file
	f, err := os.OpenFile(p.Config.DeployEnvPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// Output file in different formats
	fmt.Printf("Writing into '%s' (format: %s)...\n", p.Config.DeployEnvPath, p.Config.OutputFormat)
	if p.Config.OutputFormat == "dotenv" {
		// dotenv format
		for _, k := range secretKeys {
			v := secret.Data[k]

			// Convert to string if boolean
			if _, ok := v.(bool); ok {
				v = strconv.FormatBool(v.(bool))
			}

			v = strings.Replace(fmt.Sprintf("%s", v), "\n", "\\n", -1)
			v = strings.Replace(fmt.Sprintf("%s", v), "\r", "\\r", -1)
			v = strings.Replace(fmt.Sprintf("%s", v), "%", "%%", -1)
			fmt.Fprintf(f, fmt.Sprintf("%s='%s'\n", k, v))
		}
	} else if p.Config.OutputFormat == "helm-yaml" {
		// yaml format, basically for helm values.yaml
		fmt.Fprintf(f, fmt.Sprintf("%s:\n", p.Config.HelmEnvKey))
		for _, k := range secretKeys {
			v := secret.Data[k]

			v = strings.Replace(fmt.Sprintf("%s", v), "\n", "\\n", -1)
			fmt.Fprintf(f, fmt.Sprintf("  %s: \"%s\"\n", k, v))
		}
	} else {
		return fmt.Errorf(noSuchOutputFormat)
	}

	return nil
}
