package main

import (
	"github.com/MaximeHeckel/temaki_generator/generator"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	"fmt"
)

func main() {
	fmt.Println(generator.Generate())
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{Directory: "generator/templates"}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "main", generator.Generate())
	})

	m.Run()
}
