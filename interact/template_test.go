package interact

import (
	"fmt"
	"os"
	"testing"
	"text/template"
)

type Person struct {
	Name                string
	nonExportedAgeField string
}

func TestTemplateExc(tt *testing.T) {
	t := template.New("hello")
	t, _ = t.Parse("hello {{.Name}}!")
	p := Person{Name: "Mary", nonExportedAgeField: "31"}
	if err := t.Execute(os.Stdout, p); err != nil {
		fmt.Println("There was an error:", err.Error())
	}
}

func TestTemplateValidation(t *testing.T) {
	//为了确保模板定义语法是正确的，使用 Must 函数处理 Parse 的返回结果
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("run time panic: %v", err)
		}
	}()
	tOk := template.New("ok")
	//a valid template, so no panic with Must:
	template.Must(tOk.Parse("/* and a comment */ some static text: {{ .Name }}"))
	fmt.Println("The first one parsed OK.")
	fmt.Println("The next one ought to fail.")
	tErr := template.New("error_template")
	template.Must(tErr.Parse(" some static text {{ .Name }"))
}

func TestTemplateIfElse(t *testing.T) {
	//对管道数据的输出结果用 if-else-end 设置条件约束：如果管道是空的
	tEmpty := template.New("template test")
	tEmpty = template.Must(tEmpty.Parse("Empty pipeline if demo: {{if ``}} Will not print. {{end}}\n")) //empty pipeline following if
	tEmpty.Execute(os.Stdout, nil)

	tWithValue := template.New("template test")
	tWithValue = template.Must(tWithValue.Parse("Non empty pipeline if demo: {{if `anything`}} Will print. {{end}}\n")) //non empty pipeline following if condition
	tWithValue.Execute(os.Stdout, nil)

	tIfElse := template.New("template test")
	tIfElse = template.Must(tIfElse.Parse("if-else demo: {{if `anything`}} Print IF part. {{else}} Print ELSE part.{{end}}\n")) //non empty pipeline following if condition
	tIfElse.Execute(os.Stdout, nil)
}

func TestTemplateWithEnd(tt *testing.T) {
	//点号（.）可以在 Go 模板中使用：其值 {{.}} 被设置为当前管道的值。
	//with 语句将点号设为管道的值。如果管道是空的，那么不管 with-end 块之间有什么，都会被忽略。
	// 在被嵌套时，点号根据最近的作用域取得值。
	t := template.New("test")
	t, _ = t.Parse("{{with `hello`}}{{.}}{{end}}!\n")
	t.Execute(os.Stdout, nil)

	t, _ = t.Parse("{{with `hello`}}{{.}} {{with `Mary`}}{{.}}{{end}}{{end}}!\n")
	t.Execute(os.Stdout, nil)
}

func TestTemplateDollar(tt *testing.T) {
	//可以在模板内为管道设置本地变量，变量名以 $ 符号作为前缀。变量名只能包含字母、数字和下划线。
	t := template.New("test")
	t = template.Must(t.Parse("{{with $3 := `hello`}}{{$3}}{{end}}!\n"))
	t.Execute(os.Stdout, nil)

	t = template.Must(t.Parse("{{with $x3 := `hola`}}{{$x3}}{{end}}!\n"))
	t.Execute(os.Stdout, nil)

	t = template.Must(t.Parse("{{with $x_1 := `hey`}}{{$x_1}} {{.}} {{$x_1}}{{end}}!\n"))
	t.Execute(os.Stdout, nil)
}

func TestTemplateRangeEnd(tt *testing.T) {
	//range-end 结构格式为：{{range pipeline}} T1 {{else}} T0 {{end}}。
	//range 被用于在集合上迭代：管道的值必须是数组、切片或 map。
	// 如果管道的值长度为零，点号的值不受影响，且执行 T0；
	// 否则，点号被设置为数组、切片或 map 内元素的值，并执行 T1。
	t := template.New("test")
	t = template.Must(t.Parse(`{{range .}}aa{{.}}bb{{end}}`))
	s := []int{1, 2, 3, 4}
	t.Execute(os.Stdout, s)
}

var appText = `{{range .}}
	{{with .Author}}
		<p><b>{{html .}}</b> wrote:</p>
	{{else}}
		<p>An anonymous person wrote:</p>
	{{end}}
	<pre>{{html .Content}}</pre>
	<pre>{{html .Date}}</pre>
{{end}}`

func TestTemplateFunc(tt *testing.T) {
	t := template.New("test")
	t = template.Must(t.Parse("{{with $x := `hello`}}{{printf `%s %s` $x `Mary`}}{{end}}!\n"))
	t.Execute(os.Stdout, nil)
}
