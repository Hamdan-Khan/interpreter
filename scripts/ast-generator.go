package scripts

import (
	"fmt"
	"os"
	"strings"
)

// generates ast classes
func GenerateAst() {
    types := []string{
        "Binary   : left Expr, operator Token, right Expr",
        "Grouping : expression Expr",
        "Literal  : value any",
        "Unary    : operator Token, right Expr",
    }

    var builder strings.Builder

    builder.WriteString("package syntax\n\n")

    // define visitor interface
    defineVisitor(&builder, "Expr", types)

    builder.WriteString("type Expr interface {\n")
    builder.WriteString("    Accept(visitor Visitor) any\n")
    builder.WriteString("}\n\n")

    // generate each concrete type
    for _, typeStr := range types {
        parts := strings.Split(typeStr, ":")
        className := strings.TrimSpace(parts[0])
        fieldList := strings.TrimSpace(parts[1])
        
        defineType(&builder, className, fieldList)
    }
    
    os.WriteFile("syntax/syntax.go", []byte(builder.String()), 0644)
}

func defineVisitor(builder *strings.Builder, baseName string, types []string) {
    builder.WriteString("type Visitor interface {\n")

    for _, typeStr := range types {
        typeName := strings.TrimSpace(strings.Split(typeStr, ":")[0])
        fmt.Fprintf(builder, "    Visit%s%s(%s *%s) any\n",
                    typeName, baseName, strings.ToLower(baseName), typeName)
    }
    
    builder.WriteString("}\n\n")
}

func defineType(builder *strings.Builder, className, fieldList string) {
    fmt.Fprintf(builder, "type %s struct {\n", className)
    
    if fieldList != "" {
        fields := strings.SplitSeq(fieldList, ", ")
        for field := range fields {
            parts := strings.Fields(field)
            if len(parts) >= 2 {
                fieldName := capitalize(parts[0])
                fieldType := parts[1]
                fmt.Fprintf(builder, "    %s %s\n", fieldName, fieldType)
            }
        }
    }

    builder.WriteString("}\n\n")

    // implement the Accept method (Visitor pattern)
    fmt.Fprintf(builder, "func (e *%s) Accept(visitor Visitor) any {\n", className)
    fmt.Fprintf(builder, "    return visitor.Visit%sExpr(e)\n", className)
    builder.WriteString("}\n\n")
}

func capitalize(s string) string {
    if s == "" {
        return s
    }
    return strings.ToUpper(s[:1]) + s[1:]
}
