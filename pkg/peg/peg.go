package peg

type ASTNode any
type Token map[string]ASTNode
type Parser func(c *Context) (ASTNode, bool)
