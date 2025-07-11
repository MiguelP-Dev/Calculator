package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Operation representa una operación realizada por la calculadora, guardando la expresión y el resultado.
type Operation struct {
	Expression string // Ejemplo: "2 + 2"
	Result     string // Ejemplo: "4"
}

// Calculator contiene el estado y la lógica de la calculadora, así como referencias a los widgets de la UI.
type Calculator struct {
	history     []Operation   // Historial de operaciones
	historyList *widget.List  // Referencia al widget de historial en la UI
	resultLabel *canvas.Text  // Display de resultado
	inputEntry  *widget.Entry // Campo de entrada de operaciones
}

// NewCalculator crea una nueva instancia de Calculator con valores iniciales.
func NewCalculator() *Calculator {
	return &Calculator{
		history: []Operation{},
	}
}

// evalExpr evalúa una expresión matemática con soporte para paréntesis, jerarquía de operaciones y números negativos
func evalExpr(expr string) (float64, error) {
	tokens, err := tokenizeAdvanced(expr)
	if err != nil {
		return 0, err
	}
	if len(tokens) == 0 {
		return 0, fmt.Errorf("empty expression")
	}
	res, idx, err := parseExpr(tokens, 0)
	if err != nil {
		return 0, err
	}
	if idx != len(tokens) {
		return 0, fmt.Errorf("invalid syntax: unexpected token '%s'", tokens[idx])
	}
	return res, nil
}

// tokenizeAdvanced separa la expresión en números, operadores y paréntesis, validando caracteres
func tokenizeAdvanced(expr string) ([]string, error) {
	expr = strings.ReplaceAll(expr, " ", "")
	tokens := []string{}
	num := ""
	for i, r := range expr {
		if unicode.IsDigit(r) || r == '.' {
			num += string(r)
			continue
		}
		if r == '-' && (i == 0 || expr[i-1] == '(' || isOperator(rune(expr[i-1]))) {
			// Soporta negativos al inicio, después de '(', o después de operador
			num += string(r)
			continue
		}
		if num != "" {
			tokens = append(tokens, num)
			num = ""
		}
		if isOperator(r) || r == '(' || r == ')' {
			tokens = append(tokens, string(r))
			continue
		}
		return nil, fmt.Errorf("invalid character in expression: '%c'", r)
	}
	if num != "" {
		tokens = append(tokens, num)
	}
	return tokens, nil
}

func isOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' || r == '%'
}

// parseExpr implementa la jerarquía de operaciones y paréntesis
func parseExpr(tokens []string, idx int) (float64, int, error) {
	return parseAddSub(tokens, idx)
}

func parseAddSub(tokens []string, idx int) (float64, int, error) {
	res, idx, err := parseMulDivMod(tokens, idx)
	if err != nil {
		return 0, idx, err
	}
	for idx < len(tokens) {
		tok := tokens[idx]
		if tok != "+" && tok != "-" {
			break
		}
		idx++
		right, nextIdx, err := parseMulDivMod(tokens, idx)
		if err != nil {
			return 0, idx, err
		}
		if tok == "+" {
			res += right
		} else {
			res -= right
		}
		idx = nextIdx
	}
	return res, idx, nil
}

func parseMulDivMod(tokens []string, idx int) (float64, int, error) {
	res, idx, err := parseFactor(tokens, idx)
	if err != nil {
		return 0, idx, err
	}
	for idx < len(tokens) {
		tok := tokens[idx]
		if tok != "*" && tok != "/" && tok != "%" {
			break
		}
		idx++
		right, nextIdx, err := parseFactor(tokens, idx)
		if err != nil {
			return 0, idx, err
		}
		switch tok {
		case "*":
			res *= right
		case "/":
			if right == 0 {
				res = res // mismo comportamiento: devolver el número original
			} else {
				res /= right
			}
		case "%":
			if right == 0 {
				res = res
			} else {
				res = float64(int(res) % int(right))
			}
		}
		idx = nextIdx
	}
	return res, idx, nil
}

func parseFactor(tokens []string, idx int) (float64, int, error) {
	if idx >= len(tokens) {
		return 0, idx, fmt.Errorf("unexpected end of expression")
	}
	tok := tokens[idx]
	if tok == "(" {
		res, nextIdx, err := parseExpr(tokens, idx+1)
		if err != nil {
			return 0, idx, err
		}
		if nextIdx >= len(tokens) || tokens[nextIdx] != ")" {
			return 0, idx, fmt.Errorf("unmatched parenthesis")
		}
		return res, nextIdx + 1, nil
	}
	// Soporta números negativos
	if tok == "-" && idx+1 < len(tokens) && (isNumber(tokens[idx+1]) || tokens[idx+1] == "(") {
		res, nextIdx, err := parseFactor(tokens, idx+1)
		if err != nil {
			return 0, idx, err
		}
		return -res, nextIdx, nil
	}
	if isNumber(tok) {
		val, err := strconv.ParseFloat(tok, 64)
		if err != nil {
			return 0, idx, fmt.Errorf("invalid number: %s", tok)
		}
		return val, idx + 1, nil
	}
	return 0, idx, fmt.Errorf("invalid token: %s", tok)
}

func isNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s[0] == '-' && len(s) > 1 {
		s = s[1:]
	}
	for _, r := range s {
		if !unicode.IsDigit(r) && r != '.' {
			return false
		}
	}
	return true
}

// borderedWhite crea una caja blanca con ancho fijo y padding
func borderedWhite(content fyne.CanvasObject, width float32, height float32) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(width, height))
	return container.NewStack(
		rect,
		container.NewPadded(content),
	)
}

// updateHistory refresca el widget de historial.
func (c *Calculator) updateHistory() {
	if c.historyList != nil {
		c.historyList.Refresh()
	}
}

func main() {
	calculatorApp := app.New()
	mainWindow := calculatorApp.NewWindow("Calculadora Fyne")

	calc := NewCalculator()

	// Campo de entrada único para operaciones
	entry := widget.NewEntry()
	calc.inputEntry = entry

	// Display de resultado (canvas.Text para tamaño grande)
	resultText := canvas.NewText("", color.Black)
	resultText.Alignment = fyne.TextAlignTrailing
	resultText.TextStyle = fyne.TextStyle{Bold: true}
	resultText.TextSize = 28
	calc.resultLabel = resultText
	resultCard := borderedWhite(resultText, 400, 60)

	// Función para evaluar y mostrar resultado
	evalAndShow := func() {
		expr := entry.Text
		if expr == "" {
			return
		}
		res, err := evalExpr(expr)
		if err != nil {
			resultText.Text = err.Error()
			resultText.Refresh()
			return
		}
		resultStr := strconv.FormatFloat(res, 'f', -1, 64)
		resultText.Text = resultStr
		resultText.Refresh()
		calc.history = append([]Operation{{Expression: expr, Result: resultStr}}, calc.history...)
		calc.updateHistory()
	}

	// Asignar Enter para evaluar
	entry.OnSubmitted = func(_ string) {
		evalAndShow()
	}

	// Botones y su acción: insertan símbolo en el Entry
	addToEntry := func(s string) func() {
		return func() {
			entry.SetText(entry.Text + s)
			entry.CursorColumn = len(entry.Text)
		}
	}
	setEntry := func(s string) func() {
		return func() {
			entry.SetText(s)
			entry.CursorColumn = len(s)
		}
	}

	// Botones numéricos y de operaciones
	btn7 := widget.NewButton("7", addToEntry("7"))
	btn8 := widget.NewButton("8", addToEntry("8"))
	btn9 := widget.NewButton("9", addToEntry("9"))
	btnDiv := widget.NewButton("/", addToEntry("/"))
	btn4 := widget.NewButton("4", addToEntry("4"))
	btn5 := widget.NewButton("5", addToEntry("5"))
	btn6 := widget.NewButton("6", addToEntry("6"))
	btnMul := widget.NewButton("*", addToEntry("*"))
	btn1 := widget.NewButton("1", addToEntry("1"))
	btn2 := widget.NewButton("2", addToEntry("2"))
	btn3 := widget.NewButton("3", addToEntry("3"))
	btnSub := widget.NewButton("-", addToEntry("-"))
	btn0 := widget.NewButton("0", addToEntry("0"))
	btnDot := widget.NewButton(".", addToEntry("."))
	btnEq := widget.NewButton("=", func() { evalAndShow() })
	btnAdd := widget.NewButton("+", addToEntry("+"))
	btnMod := widget.NewButton("%", addToEntry("%"))
	btnC := widget.NewButton("C", func() { setEntry("")(); resultText.Text = ""; resultText.Refresh() })

	// Teclado organizado de forma tradicional
	keypad := container.NewVBox(
		container.NewGridWithColumns(4, btn7, btn8, btn9, btnDiv),
		container.NewGridWithColumns(4, btn4, btn5, btn6, btnMul),
		container.NewGridWithColumns(4, btn1, btn2, btn3, btnSub),
		container.NewGridWithColumns(4, btn0, btnDot, btnEq, btnAdd),
		container.NewGridWithColumns(2, btnMod, btnC),
	)

	// Historial con canvas.Text para tamaño pequeño
	historyList := widget.NewList(
		func() int { return len(calc.history) },
		func() fyne.CanvasObject {
			lbl := canvas.NewText("", color.Black)
			lbl.TextSize = 12
			lbl.Alignment = fyne.TextAlignLeading
			return lbl
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			lbl := o.(*canvas.Text)
			if len(calc.history[i].Expression) > 30 {
				lbl.Text = "= " + calc.history[i].Result
			} else {
				lbl.Text = calc.history[i].Expression + " = " + calc.history[i].Result
			}
			lbl.Refresh()
		},
	)
	calc.historyList = historyList
	calc.historyList.OnSelected = func(id int) {
		setEntry(calc.history[id].Result)()
	}
	historyScroll := container.NewVScroll(borderedWhite(historyList, 400, 120))
	historyScroll.SetMinSize(fyne.NewSize(400, 120))

	mainContent := container.NewVBox(
		borderedWhite(entry, 400, 28),
		resultCard,
		keypad,
		historyScroll,
	)

	mainWindow.SetContent(mainContent)

	// Captura global de teclas físicas
	mainWindow.Canvas().SetOnTypedRune(func(r rune) {
		if unicode.IsDigit(r) || strings.ContainsRune(".+-*/%", r) {
			entry.SetText(entry.Text + string(r))
			entry.CursorColumn = len(entry.Text)
		}
		if r == '=' {
			evalAndShow()
		}
		if r == 'c' || r == 'C' {
			setEntry("")()
			resultText.Text = ""
			resultText.Refresh()
		}
	})
	mainWindow.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
			evalAndShow()
		}
	})

	mainWindow.ShowAndRun()
}
