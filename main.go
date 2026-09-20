package main

import(
"fmt"
"blog-aggregator/internal/config"
)




func main() {
	var c config.Config
	cfg, err := c.ReadFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	newCfg, err := cfg.SetUser("Grimahed")
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg2, err := newCfg.ReadFile()
	if err != nil {
                fmt.Println(err)
                return
        }
	fmt.Println(cfg2)
	fmt.Println("I compiled :)")

}
