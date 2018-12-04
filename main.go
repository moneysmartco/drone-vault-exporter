package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// build number set at compile-time
var build = "0"

// Version set at compile-time
var Version string

func main() {
	if Version == "" {
		Version = fmt.Sprintf("0.0.1+%s", build)
	}

	app := cli.NewApp()
	app.Name = "Drone Vault Exporter"
	app.Usage = "Export Vault to dotenv format"
	app.Copyright = "Copyright (c) 2018 Eric Ho"
	app.Authors = []cli.Author{
		{
			Name:  "Eric Ho",
			Email: "dho.eric@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "vault-addr",
			Usage:  "Vault server address",
			EnvVar: "VAULT_ADDR",
		},
		cli.StringFlag{
			Name:   "vault-token",
			Usage:  "Vault token for auth",
			EnvVar: "VAULT_TOKEN",
		},
		cli.StringFlag{
			Name:   "vault-auth-method",
			Usage:  "Authorization method of Vault",
			EnvVar: "PLUGIN_VAULT_AUTH_METHOD",
		},
		cli.StringFlag{
			Name:   "vault-role-id",
			Usage:  "Vault role ID when using vault approle login",
			EnvVar: "VAULT_ROLE_ID",
		},
		cli.StringFlag{
			Name:   "vault-secret-id",
			Usage:  "Vault secret ID when using vault approle login",
			EnvVar: "VAULT_SECRET_ID",
		},
		cli.StringFlag{
			Name:   "vault-key-path",
			Usage:  "Vault key path to extract the secrets",
			EnvVar: "PLUGIN_VAULT_KEY_PATH",
		},
		cli.StringFlag{
			Name:   "deploy-env-path",
			Usage:  "Path to save the dotenv file",
			EnvVar: "PLUGIN_DEPLOY_ENV_PATH",
			Value:  ".deploy.env",
		},
		cli.StringFlag{
			Name:   "output-format",
			Usage:  "Format to be saved, dotenv / helm-yaml",
			EnvVar: "PLUGIN_OUTPUT_FORMAT",
			Value:  "dotenv",
		},
		cli.StringFlag{
			Name:   "helm-env-key",
			Usage:  "Helm key for env vars",
			EnvVar: "PLUGIN_HELM_ENV_KEY",
			Value:  "envs",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	// Override a template
	cli.AppHelpTemplate = `
NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
REPOSITORY:
    Github: https://github.com/moneysmartco/drone-vault-exporter
`

	if err := app.Run(os.Args); err != nil {
		fmt.Println("drone-vault-exporter Error: ", err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Config: Config{
			VaultAddr:       c.String("vault-addr"),
			VaultToken:      c.String("vault-token"),
			VaultAuthMethod: c.String("vault-auth-method"),
			VaultRoleID:     c.String("vault-role-id"),
			VaultSecretID:   c.String("vault-secret-id"),
			VaultKeyPath:    c.String("vault-key-path"),
			DeployEnvPath:   c.String("deploy-env-path"),
			HelmEnvKey:      c.String("helm-env-key"),
			OutputFormat:    c.String("output-format"),
		},
	}

	return plugin.Exec()
}
