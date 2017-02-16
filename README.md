# goblin
micro web framework for go


## To do list
- [x] router
- [ ] controller
- [ ] view

## sample
### main.go
```
package main

import (
     "goblin"
     _ "sumi/routers"
)

func main() {
    goblin.Run(":9090")
}

```
### routers/routers.go
```
package routers

import (
    "goblin"
    "sumi/controllers"
)

func init() {
    goblin.Get("/books",    &controllers.BooksController{}, "Index")
    goblin.Get("/books/1",  &controllers.BooksController{}, "Show")
}

```
### controllers/BooksController.go
```
package controllers

import (
    "goblin"
    "log"
 )

 type BooksController struct {
     goblin.Controller
 }

 func (this *BooksController) Index() {
     log.Println("=> in book list")
 }

 func (this *BooksController) Show() {
     log.Println("=> in book show")
 }

 func (this *BooksController) Create() {
     log.Println("=> create new book")
 }

```
