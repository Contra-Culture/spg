package spg

import "github.com/Contra-Culture/go2html"

type (
	PageType      string
	PageGenerator struct {
		Name     string
		Layout   *go2html.Template
		Children map[string]*PageGenerator
	}
	PageGeneratorCfgr struct {
		pageGenerator *PageGenerator
	}
)

func (c *PageGeneratorCfgr) Page(cfg func(*PageGeneratorCfgr)) {

}
func (c *PageGeneratorCfgr) Layout(cfg func(*go2html.TemplateConfiguringProxy)) {

}
func (c *PageGeneratorCfgr) Screen(cfg func(*go2html.TemplateConfiguringProxy)) {

}
func (c *PageGeneratorCfgr) Schema(n string, arrows []string, cfg func(*go2html.TemplateConfiguringProxy)) {

}
