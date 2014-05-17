package main

import (
  "fmt"
  "github.com/aforward/webl/api"
  "gopkg.in/fatih/set.v0"
)

type Edge struct {
  FromName string
  ToName string
}


func showVersion() {
  webl.INFO.Println(fmt.Sprintf("weblui %s", webl.Version()))
}

func flatten(edges []Edge, node *webl.Resource, alreadyProcessed *set.Set) []Edge {
  // root.Links = append(root.Links,node.Links...)
  for _,link := range node.Links {
    edges = append(edges,Edge{ FromName: node.Name, ToName: link.Name })
    if !alreadyProcessed.Has(link.Url) {
      alreadyProcessed.Add(link.Url)
      flatten(edges,&link,alreadyProcessed)
    }
  } 
  return edges
}