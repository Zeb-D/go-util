package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// ast文本内容
var src = `
 package main
 import "fmt"
 func main() {
     fmt.Println("Hello, World!")
 }
`

func TestAst(t *testing.T) {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST. //具体输出内容如下，额外人为注释，使用 //=
	ast.Print(fset, f)
}

//0  *ast.File {
//     1  .  Package: 2:1  //= 包声明
//     2  .  Name: *ast.Ident {
//     3  .  .  NamePos: 2:9 //= 包的位置
//     4  .  .  Name: "main" //= 包的名称
//     5  .  }
//     6  .  Decls: []ast.Decl (len = 2) {
//     7  .  .  0: *ast.GenDecl {
//     8  .  .  .  TokPos: 4:1
//     9  .  .  .  Tok: import //= 引入声明
//    10  .  .  .  Lparen: -
//    11  .  .  .  Specs: []ast.Spec (len = 1) { //=引入声明数组
//    12  .  .  .  .  0: *ast.ImportSpec {
//    13  .  .  .  .  .  Path: *ast.BasicLit {
//    14  .  .  .  .  .  .  ValuePos: 4:8
//    15  .  .  .  .  .  .  Kind: STRING
//    16  .  .  .  .  .  .  Value: "\"fmt\"" //=引入声明的情况
//    17  .  .  .  .  .  }
//    18  .  .  .  .  .  EndPos: -
//    19  .  .  .  .  }
//    20  .  .  .  }
//    21  .  .  .  Rparen: -
//    22  .  .  }
//    23  .  .  1: *ast.FuncDecl { //=函数声明
//    24  .  .  .  Name: *ast.Ident {
//    25  .  .  .  .  NamePos: 6:6
//    26  .  .  .  .  Name: "main" //=函数名称
//    27  .  .  .  .  Obj: *ast.Object {
//    28  .  .  .  .  .  Kind: func
//    29  .  .  .  .  .  Name: "main"
//    30  .  .  .  .  .  Decl: *(obj @ 23)
//    31  .  .  .  .  }
//    32  .  .  .  }
//    33  .  .  .  Type: *ast.FuncType { //= 函数类型
//    34  .  .  .  .  Func: 6:1
//    35  .  .  .  .  Params: *ast.FieldList { //=函数声明的请求参数
//    36  .  .  .  .  .  Opening: 6:10
//    37  .  .  .  .  .  Closing: 6:11
//    38  .  .  .  .  } //=没出现Results，表示无返回结果
//    39  .  .  .  }
//    40  .  .  .  Body: *ast.BlockStmt { //=函数体信息
//    41  .  .  .  .  Lbrace: 6:13
//    42  .  .  .  .  List: []ast.Stmt (len = 1) { //=len = 1 代码行
//    43  .  .  .  .  .  0: *ast.ExprStmt { //=第一个代码行内容
//    44  .  .  .  .  .  .  X: *ast.CallExpr { //= CallExpr 函数调用
//    45  .  .  .  .  .  .  .  Fun: *ast.SelectorExpr { //= SelectorExpr 调用函数的包名及函数名
//    46  .  .  .  .  .  .  .  .  X: *ast.Ident {
//    47  .  .  .  .  .  .  .  .  .  NamePos: 7:5
//    48  .  .  .  .  .  .  .  .  .  Name: "fmt"
//    49  .  .  .  .  .  .  .  .  }
//    50  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
//    51  .  .  .  .  .  .  .  .  .  NamePos: 7:9
//    52  .  .  .  .  .  .  .  .  .  Name: "Println"
//    53  .  .  .  .  .  .  .  .  }
//    54  .  .  .  .  .  .  .  }
//    55  .  .  .  .  .  .  .  Lparen: 7:16
//    56  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) { //= 参数信息
//    57  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
//    58  .  .  .  .  .  .  .  .  .  ValuePos: 7:17
//    59  .  .  .  .  .  .  .  .  .  Kind: STRING
//    60  .  .  .  .  .  .  .  .  .  Value: "\"Hello World!\""
//    61  .  .  .  .  .  .  .  .  }
//    62  .  .  .  .  .  .  .  }
//    63  .  .  .  .  .  .  .  Ellipsis: -
//    64  .  .  .  .  .  .  .  Rparen: 7:31
//    65  .  .  .  .  .  .  }
//    66  .  .  .  .  .  }
//    67  .  .  .  .  }
//    68  .  .  .  .  Rbrace: 8:1
//    69  .  .  .  }
//    70  .  .  }
//    71  .  }
//    72  .  Scope: *ast.Scope {
//    74  .  .  Objects: map[string]*ast.Object (len = 1) {
//    74  .  .  .  "main": *(obj @ 27)
//    75  .  .  }
//    76  .  }
//    77  .  Imports: []*ast.ImportSpec (len = 1) {
//    78  .  .  0: *(obj @ 12)
//    79  .  }
//    80  .  Unresolved: []*ast.Ident (len = 1) {
//    81  .  .  0: *(obj @ 46)
//    82  .  }
//    83  }
