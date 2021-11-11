package data_test

import (
	. "github.com/Contra-Culture/spg/data"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("data", func() {
	Describe("*Graph", func() {
		Describe("creation | New()", func() {
			Context("when valid specification", func() {
				It("returns *Graph", func() {
					g, err := New(func(cfg *GraphCfgr) {})
					Expect(err).NotTo(HaveOccurred())
					Expect(g).NotTo(BeNil())
				})
			})
			Describe(".Schema() specification", func() {
				Context("when valid specification", func() {
					It("adds Schema specification to the graph", func() {
						g, err := New(func(cfg *GraphCfgr) {
							cfg.Schema("account", func(cfg *SchemaCfgr) {
								cfg.ID([]string{"login"})
								cfg.Attribute("login")
								cfg.Attribute("firstName")
								cfg.Attribute("lastName")
								cfg.Attribute("bio")
							})
							cfg.Schema("rubric", func(cfg *SchemaCfgr) {
								cfg.ID([]string{"slug"})
								cfg.Attribute("slug")
								cfg.Attribute("title")
								cfg.Attribute("description")
							})
							cfg.Schema("publication",func(cfg *SchemaCfgr){
								cfg.ID([]string{"slug"})
								cfg.Attribute("slug")
								cfg.Attribute("title")
								cfg.Attribute("content")
								cfg.Attribute("createdAt")
								cfg.Attribute("updatedAt")
								cfg.Attribute("accountLogin")
							})
						})
						Expect(err).To(BeNil())
						Expect(g).NotTo(BeNil())
					})
				})
			})

		})
	})
})
