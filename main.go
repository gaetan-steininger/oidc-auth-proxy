package main

import (
	"context"
	"log"
	"oidc-auth-proxy/server"
	"oidc-auth-proxy/version"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Print version",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					version.PrintVersion(
						cmd.Bool("json"),
					)
					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "json",
						OnlyOnce: true,
						Value:    false,
						Usage:    "print as json",
					},
				},
			},
			{
				Name:  "server",
				Usage: "Run server",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					server.Run(
						cmd.String("session-cookie-name"),
						cmd.String("session-secret"),
						int(cmd.Int("session-max-age")),
						cmd.String("backend-url"),
						cmd.StringSlice("oidc-scopes"),
						cmd.String("oidc-redirect-url"),
						cmd.String("oidc-client-id"),
						cmd.String("oidc-client-secret"),
						cmd.String("oidc-provider-url"),
						cmd.String("oidc-user-claim"),
						cmd.String("oidc-groups-claim"),
						cmd.StringSlice("allowed-users"),
						cmd.StringSlice("allowed-groups"),
					)
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "session-cookie-name",
						OnlyOnce: true,
						Value:    "session",
						Usage:    "Session cookie name",
						Sources:  cli.EnvVars("SESSION_COOKIE_NAME"),
					},
					&cli.StringFlag{
						Name:     "session-secret",
						OnlyOnce: true,
						Required: true,
						Usage:    "Session secret",
						Sources:  cli.EnvVars("SESSION_SECRET"),
					},
					&cli.IntFlag{
						Name:     "session-max-age",
						OnlyOnce: true,
						Value:    3600,
						Usage:    "Session max age (seconds)",
						Sources:  cli.EnvVars("SESSION_MAX_AGE"),
					},
					&cli.StringFlag{
						Name:     "backend-url",
						OnlyOnce: true,
						Value:    "http://localhost:9000",
						Usage:    "Backend server URL",
						Sources:  cli.EnvVars("BACKEND_URL"),
					},
					&cli.StringSliceFlag{
						Name:     "oidc-scopes",
						OnlyOnce: true,
						Value:    []string{"openid", "profile", "email"},
						Usage:    "OIDC scopes list",
						Sources:  cli.EnvVars("OIDC_SCOPES"),
					},
					&cli.StringFlag{
						Name:     "oidc-redirect-url",
						OnlyOnce: true,
						Value:    "http://localhost:8080/callback",
						Usage:    "Post login redirect url",
						Sources:  cli.EnvVars("OIDC_REDIRECT_URL"),
					},
					&cli.StringFlag{
						Name:     "oidc-provider-url",
						OnlyOnce: true,
						Required: true,
						Usage:    "OIDC provider URL",
						Sources:  cli.EnvVars("OIDC_PROVIDER_URL"),
					},
					&cli.StringFlag{
						Name:     "oidc-client-id",
						OnlyOnce: true,
						Required: true,
						Usage:    "OIDC client id",
						Sources:  cli.EnvVars("OIDC_CLIENT_ID"),
					},
					&cli.StringFlag{
						Name:     "oidc-client-secret",
						OnlyOnce: true,
						Required: true,
						Usage:    "OIDC client secret",
						Sources:  cli.EnvVars("OIDC_CLIENT_SECRET"),
					},
					&cli.StringFlag{
						Name:     "oidc-user-claim",
						OnlyOnce: true,
						Value:    "preferred_username",
						Usage:    "OIDC claim to get username",
						Sources:  cli.EnvVars("OIDC_USER_CLAIM"),
					},
					&cli.StringFlag{
						Name:     "oidc-groups-claim",
						OnlyOnce: true,
						Value:    "",
						Usage:    "OIDC claim to get groups (ignored if \"\")",
						Sources:  cli.EnvVars("OIDC_GROUPS_CLAIM"),
					},
					&cli.StringSliceFlag{
						Name:     "allowed-users",
						OnlyOnce: true,
						Value:    []string{},
						Usage:    "List of authorized users (ignored if not set)",
						Sources:  cli.EnvVars("ALLOWED_USERS"),
					},
					&cli.StringSliceFlag{
						Name:     "allowed-groups",
						OnlyOnce: true,
						Value:    []string{},
						Usage:    "List of authorized groups (ignored if not set)",
						Sources:  cli.EnvVars("ALLOWED_GROUPS"),
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
