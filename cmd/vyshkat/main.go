package main

import (
	"flag"
	"fmt"
	"os"

	"vyshkat/internal/bootstrap"
	"vyshkat/internal/console"
	"vyshkat/internal/demo"
)

func main() {
	skipDemo := flag.Bool("no-demo", false, "запустить без демо")
	onlyDemo := flag.Bool("demo-only", false, "показать только демо")
	flag.Parse()

	if !*skipDemo || *onlyDemo {
		err := demo.Run(os.Stdout)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка при запуске демо:", err)
			os.Exit(1)
		}
	}

	if *onlyDemo {
		return
	}

	container := bootstrap.NewContainer()
	park, err := bootstrap.Fleet(container)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Не получилось создать парк:", err)
		os.Exit(1)
	}

	app := console.New(os.Stdin, os.Stdout, park)

	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка программы:", err)
		os.Exit(1)
	}
}
