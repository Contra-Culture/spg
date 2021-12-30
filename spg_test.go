package spg_test

import (
	"fmt"
	"os"
	"strings"

	. "github.com/Contra-Culture/spg"
	"github.com/Contra-Culture/spg/data"
	"github.com/Contra-Culture/spg/node"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("spg", func() {
	It("works", func() {
		path, err := os.Executable()
		if err != nil {
			panic(err)
		}
		chunks := strings.Split(path, "/")
		path = strings.Join(chunks[:len(chunks)-1], "/")
		Expect(
			New(
				"test",
				"example.com",
				func(cfg *HostCfgr) {
					cfg.DataGraph(
						fmt.Sprintf("%s/test/data", path),
						func(cfg *data.GraphCfgr) {
							cfg.Node(
								"author",
								func(cfg *data.NodeCfgr) {
									cfg.String("full-name")
									cfg.String("slug")
									cfg.String("bio")
									cfg.Link(
										"publications",
										[]string{"publication"},
										func(cfgr *data.LinkCfgr) {

										})
								})
							cfg.Node(
								"publication",
								func(cfg *data.NodeCfgr) {
									cfg.PK([]string{"slug"})
									cfg.String("title")
									cfg.String("slug")
									cfg.String("rubric-slug")
									cfg.String("content")
									cfg.String("author-name")
									cfg.String("published-at")
									cfg.Link(
										"authors",
										[]string{"user"},
										func(cfgr *data.LinkCfgr) {

										})
									cfg.Link(
										"rubric",
										[]string{"rubric"},
										func(cfgr *data.LinkCfgr) {

										})
								})
							cfg.Node(
								"rubric",
								func(cfg *data.NodeCfgr) {
									cfg.String("title")
									cfg.Link(
										"publications",
										[]string{"publication"},
										func(cfgr *data.LinkCfgr) {

										})
								})
						})
					cfg.Root(
						func(cfg *node.NodeCfgr) {
							cfg.Node(
								"publication",
								func(cfg *node.NodeCfgr) {

								})
							cfg.Layout([]string{})
							cfg.Screen([]string{})
							cfg.Node(
								"rubric",
								func(cfg *node.NodeCfgr) {
									cfg.Node(
										"rubric",
										func(cfg *node.NodeCfgr) {

										})
									cfg.Path(
										func(o *data.Object) []string {
											title, err := o.Prop("title")
											if err != nil {
												panic(err)
											}
											title = strings.NewReplacer(" ", "-", "&", "-and-", "?", "").Replace(title)
											return []string{"rubric", title}
										})
								})
							cfg.Node(
								"publication",
								func(cfg *node.NodeCfgr) {
									cfg.Node(
										"publication",
										func(cfg *node.NodeCfgr) {

										})
									cfg.Path(
										func(o *data.Object) []string {
											title, err := o.Prop("title")
											if err != nil {
												panic(err)
											}
											title = strings.NewReplacer(" ", "-", "&", "-and-", "?", "").Replace(title)
											return []string{"publication", title}
										})
								})
						})
				})).NotTo(BeNil())
	})
})
