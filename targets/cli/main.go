package main

import (
	"encoding/json"
	"fmt"

	"github.com/Contra-Culture/cli"
	"github.com/Contra-Culture/report"
	"github.com/Contra-Culture/spg"
)

func main() {
	host := spg.New(
		"test",
		"http://example.com",
		func(hc *spg.HostCfgr) {

		})
	app, _ := cli.New(func(app *cli.AppCfgr) {
		app.Title("spg")
		app.Version("0.0.1 (alpha)")
		app.Description("Static pages generator.")
		app.HandleErrorsWith(
			func(err error) {

			})
		app.Default(
			func(cmd *cli.CommandCfgr) {
				cmd.Title("help")
				cmd.Description("Shows help info.")
				cmd.HandleWith(
					func(params map[string]string) (err error) {
						return
					})
			})
		app.Command(
			"apply",
			func(cmd *cli.CommandCfgr) {
				cmd.Title("apply")
				cmd.Description("Applies new data object: adds it to data graph and renders all the dependent pages.")
				cmd.Param(
					func(prm *cli.ParamCfgr) {
						prm.Name("schema")
						prm.Description("Data object schema name.")
						prm.CheckWith(func(report report.Node, v string) (ok bool) {
							return true
						})
						prm.Param(
							func(prm *cli.ParamCfgr) {
								prm.Name("object")
								prm.Description("data object.")
								prm.CheckWith(
									func(report report.Node, v string) (ok bool) {
										return true
									})
							})
					})
				cmd.HandleWith(
					func(params map[string]string) (err error) {
						var (
							id     string
							parsed map[string]string
							s      = params["schema"]
							o      = params["schema_object"]
						)
						json.Unmarshal([]byte(o), &parsed)
						id, err = host.Update(s, parsed)
						fmt.Printf("Object %s %s applied", s, id)
						return
					})
			})
	})
	r := app.Handle()
	if r.HasErrors() || r.HasDeprecations() || r.HasWarns() {
		fmt.Print(report.ToString(r))
	}
}
