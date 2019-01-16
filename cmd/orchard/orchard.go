package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/pgmtc/orchard-cli/internal/pkg/builder"
	"github.com/pgmtc/orchard-cli/internal/pkg/common"
	"github.com/pgmtc/orchard-cli/internal/pkg/config"
	"github.com/pgmtc/orchard-cli/internal/pkg/local"
	"github.com/pgmtc/orchard-cli/internal/pkg/source"
	"os"
	"reflect"
	"strings"
)

var (
	modules = map[string]common.Module{
		"source":  source.Module{},
		"config":  config.Module{},
		"local":   local.Module{},
		"builder": builder.Module{},
	}
	cnf    = common.FileSystemConfig("~/.orchard")
	logger = common.ConsoleLogger{}
)

func main() {
	args := os.Args[1:]
	if err := cnf.LoadConfig(); err != nil && !(len(args) == 2 && args[0] == "config" && args[1] == "init") {
		logger.Errorf("%s\n", err.Error())
		logger.Errorf("Try initializing config directory by running '%s config init'\n", os.Args[0])
		os.Exit(1)
	}
	if len(args) == 0 {
		printHelp()
		os.Exit(1)
	}
	moduleName := args[0]
	if _, ok := modules[moduleName]; !ok {
		printHelp(fmt.Sprintf("Module %s does not exist", moduleName))
		os.Exit(1)
	}

	moduleArgs := args[1:]
	os.Exit(runModule(modules[moduleName], moduleArgs...))
}

func runModule(module common.Module, args ...string) int {
	actions := module.GetActions()
	actionName := "default"
	var actionArgs []string

	if len(args) > 0 {
		actionName = args[0]
		actionArgs = args[1:]
	}

	if _, ok := actions[actionName]; !ok {
		availableActions := reflect.ValueOf(actions).MapKeys()
		logger.Errorf("Missing action '%s'. Available actions: %s\n", actionName, availableActions)
		return 1
	}

	logger.Infof("Current profile: %s\n", cnf.Config().Profile)
	action := actions[actionName]
	if err := action.Run(common.Context{Log: logger, Config: cnf}, actionArgs...); err != nil {
		logger.Errorf("Action Error: %s\n", strings.Replace(err.Error(), "\n", "", -1))
		return 2
	}
	return 0
}

func printHelp(messages ...string) {
	availableModules := reflect.ValueOf(modules).MapKeys()

	fmt.Printf("Current profile: ")
	color.HiWhite("%s", cnf.Config().Profile)
	for _, message := range messages {
		fmt.Printf(message)
	}

	fmt.Printf("Please provide module, available modules: ")
	color.HiWhite("%s", availableModules)
	fmt.Printf(" syntax : %s [module] [action]\n", os.Args[0])
	fmt.Printf(" example: %s local status\n", os.Args[0])
}
