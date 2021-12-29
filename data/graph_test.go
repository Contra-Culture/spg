package data_test

import (
	"os"
	"strings"

	"github.com/Contra-Culture/report"
	. "github.com/Contra-Culture/spg/data"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("data", func() {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	chunks := strings.Split(path, "/")
	path = strings.Join(append(chunks[:len(chunks)-2], "test/data"), "/")
	Describe("*Graph", func() {
		Describe("creation | New()", func() {
			Context("when valid specification", func() {
				It("returns *Graph", func() {
					r := report.New("graph")
					g := New(r, path, func(cfg *GraphCfgr) {})
					Expect(g).NotTo(BeNil())
					Expect(report.ToString(r)).To(Equal("| graph\n"))
				})
			})
			Describe(".Schema() specification", func() {
				Context("when valid specification", func() {
					It("adds Schema specification to the graph", func() {
						r := report.New("graph")
						g := New(r, path, func(cfg *GraphCfgr) {
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
							cfg.Schema("publication", func(cfg *SchemaCfgr) {
								cfg.ID([]string{"slug"})
								cfg.Attribute("slug")
								cfg.Attribute("title")
								cfg.Attribute("content")
								cfg.Attribute("createdAt")
								cfg.Attribute("updatedAt")
								cfg.Attribute("accountLogin")
							})
						})
						Expect(g).NotTo(BeNil())
						Expect(report.ToString(r)).To(Equal("| graph\n\t| schema: account\n\t| schema: rubric\n\t| schema: publication\n"))
					})
				})
			})
		})
		Describe("objects updating", func() {
			Describe(".Update()", func() {
				Context("when successfull", func() {
					It("returns object id", func() {
						r := report.New("graph")
						g := New(r, path, func(cfg *GraphCfgr) {
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
							cfg.Schema("publication", func(cfg *SchemaCfgr) {
								cfg.ID([]string{"slug"})
								cfg.Attribute("slug")
								cfg.Attribute("title")
								cfg.Attribute("content")
								cfg.Attribute("createdAt")
								cfg.Attribute("updatedAt")
								cfg.Attribute("accountLogin")
							})
						})
						id, err := g.Update(
							"rubric",
							map[string]string{
								"slug":        "interviews",
								"title":       "Interviews",
								"description": "Interviews with passionate people.",
							})
						Expect(err).NotTo(HaveOccurred())
						Expect(id).To(Equal("interviews"))
						o, err := g.Get("rubric", "interviews")
						Expect(err).NotTo(HaveOccurred())
						Expect(o).NotTo(BeNil())
					})
				})
			})
		})
		Describe("JSON representation", func() {
			It("returns JSON string", func() {
				r := report.New("graph")
				g := New(r, path, func(cfg *GraphCfgr) {
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
					cfg.Schema("publication", func(cfg *SchemaCfgr) {
						cfg.ID([]string{"slug"})
						cfg.Attribute("slug")
						cfg.Attribute("title")
						cfg.Attribute("content")
						cfg.Attribute("createdAt")
						cfg.Attribute("updatedAt")
						cfg.Attribute("accountLogin")
					})
				})
				Expect(g.JSONString()).To(Equal("{\"schemas\":{\"account\":{\"name\":\"account\",\"id\":[\"login\"],\"attributes\":[],\"arrows\":[]},\"rubric\":{\"name\":\"rubric\",\"id\":[\"slug\"],\"attributes\":[],\"arrows\":[]},\"publication\":{\"name\":\"publication\",\"id\":[\"slug\"],\"attributes\":[],\"arrows\":[]}}}"))
			})
		})
	})
})
