package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

/*
Functionality:
  - Load routes from config
  - Themes
  - Generate
*/

var ErrEmptyTemplate = errors.New("can't have nothing as a template")

func main() {

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "start the routing process",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					log.Println("Starting routes")
					err := StartApplication()
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "add",
				Usage:   "add resources",
				Aliases: []string{"a"},
				Commands: []*cli.Command{
					{
						Name:    "route",
						Aliases: []string{"r"},
						Usage:   "add a route",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if cmd.Args().First() == "" {
								return ErrEmptyTemplate
							}

							log.Printf("Adding route: %s\n", cmd.Args().First())
							err := AddRoute(cmd.Args().First())
							if err != nil {
								return err
							}

							return nil
						},
					},
					{
						Name:    "template",
						Aliases: []string{"t"},
						Usage:   "add a template",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							if cmd.Args().First() == "" {
								return errors.New("template name required")
							}

							log.Println("new task template: ", cmd.Args().First())

							err := AddTemplate(cmd.Args().First())
							if err != nil {
								return err
							}

							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func StartApplication() error {
	err := StartRoutes()
	if err != nil {
		return err
	}

	return nil
}
