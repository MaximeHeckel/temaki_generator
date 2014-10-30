package main

import (
  //"fmt"
  "generator"
  "github.com/go-martini/martini"
  "github.com/codegangsta/martini-contrib/render"
)

func main() {

  m := martini.Classic()
  m.Use(render.Renderer(render.Options{Directory: "src/generator/templates"}))

  m.Get("/", func(r render.Render){
      r.HTML(200,"main",generator.Generate())
  })

  m.Run()
  
}
