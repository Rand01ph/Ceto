package main

import (
	"Ceto/cgroups"
	"Ceto/cgroups/subsystems"
	"Ceto/container"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `Ceto is a simple container runtime, Just made it for fun.`

func main() {
	app := cli.NewApp()
	app.Name = "Ceto"
	app.Usage = usage
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}
	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})

		log.SetOutput(os.Stdout)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var runCommand = cli.Command{
	Name:  "run",
	Usage: `Create a container with namespace and cgroups limit ceto run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cupshare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := ctx.Args().Get(0)
		tty := ctx.Bool("ti")
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("m"),
			CpuSet:      ctx.String("cpuset"),
			CpuShare:    ctx.String("cpushare"),
		}
		Run(tty, cmd, resConf)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container",
	Action: func(ctx *cli.Context) error {
		log.Infof("init come on")
		cmd := ctx.Args().Get(0)
		log.Infof("command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}

func Run(tty bool, command string, res *subsystems.ResourceConfig) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)
	parent.Wait()
	os.Exit(-1)
}
