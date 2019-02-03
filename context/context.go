package context

import (
    //"fmt"
)

// Each module keeps its data here
type Context struct {
    Config interface {}
    DbInfo interface {}
    RestInfo interface {}
}

func CreateContext () Context {
    return Context{}
}
