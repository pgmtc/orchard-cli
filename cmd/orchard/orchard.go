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
)

var modules = map[string]common.Module{
	"source": source.Module{},
	"config": config.Module{},
}

func main() {
	args := os.Args[1:]
	if err := common.LoadConfig(); err != nil && (len(args) < 2 || args[0] != "config" || args[1] == "init") {
		color.HiRed("Error when loading config: %s", err.Error())
		color.HiRed("Try initializing config directory by running '%s config init'", os.Args[0])
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
	logger := common.ConsoleLogger{}
	actions := module.GetActions()
	actionName := "default"
	var actionArgs []string

	if len(args) > 0 {
		actionName = args[0]
		actionArgs = args[1:]
	}

	if _, ok := actions[actionName]; !ok {
		logger.Errorf("%s action does not exist. ", actionName)
		return 1
	}

	action := actions[actionName]
	if err := action.Run(logger, actionArgs...); err != nil {
		logger.Errorf("%s", err.Error())
		return 2
	}
	return 0
}

func printHelp(messages ...string) {
	availableModules := reflect.ValueOf(modules).MapKeys()

	fmt.Printf("Current profile: ")
	color.HiWhite("%s", common.CONFIG.Profile)
	for _, message := range messages {
		fmt.Printf(message)
	}

	fmt.Printf("Please provide module, available modules: ")
	color.HiWhite("%s", availableModules)
	fmt.Printf(" syntax : %s [module] [action]\n", os.Args[0])
	fmt.Printf(" example: %s local status\n", os.Args[0])
}

func main_old() {
	args := os.Args[1:]

	modules := make(map[string]func(args []string) error)
	modules["local"] = local.Parse
	//modules["source"] = source.Parse
	modules["builder"] = builder.Parse

	availableModules := reflect.ValueOf(modules).MapKeys()

	if len(args) == 0 {
		fmt.Printf("Current profile: ")
		color.HiWhite("%s", common.CONFIG.Profile)
		fmt.Printf("Please provide module, available modules: ")
		color.HiWhite("%s", availableModules)
		fmt.Printf(" syntax : %s [module] [action]\n", os.Args[0])
		fmt.Printf(" example: %s local status\n", os.Args[0])
		os.Exit(1)
	}

	module := args[0]
	handler := modules[module]

	if handler == nil {
		color.HiRed(fmt.Sprintf(" Module '%s' does not exist. Available modules = %s", module, availableModules))
		os.Exit(1)
	}

	err := handler(args[1:])
	if err != nil {
		color.HiRed(err.Error())
	} else {
		color.HiGreen("Finished OK (current profile: %s)", common.CONFIG.Profile)
	}
}
