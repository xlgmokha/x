// Package scim parses SCIM filter expressions (RFC 7644, Section 3.4.2.2).
//
// Parsing runs in two stages: a PEG grammar matches the text into a raw token
// tree, then newExpression decodes that tree into a typed Expression AST.
//
//	text
//	  |
//	  v   Grammar.Parse          grammar.go (rules) + combinators.go (helpers)
//	raw peg.Token tree
//	  |
//	  v   newExpression          decode.go
//	Expression                   expression.go (Comparison/Logical/Not/ValuePath)
//	  |
//	  +-> type switch            consumer code
//	  +-> Visit[T]               visitor.go
//	  +-> String()               expression.go (canonical, re-parseable text)
//
// Consume the returned Expression with a Go type switch, or fold it to a single
// value with a Visitor[T]. Errors live in errors.go; AttrPath in attr_path.go.
package scim
