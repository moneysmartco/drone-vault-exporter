package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

// Error message
const (
	noSuchVaultKeyPath = "no such vault key path"
	nilResponse        = "got nil response"
	emptyClientToken   = "got empty client token"
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

	// Write to file
	f, err := os.OpenFile(p.Config.DeployEnvPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	for k, v := range secret.Data {
		fmt.Fprintf(f, fmt.Sprintf("%s=%s\n", k, v))
	}

	return nil
}
