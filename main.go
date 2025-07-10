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

// evalExpr evalúa una expresión matemática simple (sin paréntesis, solo + - * / %)
func evalExpr(expr string) (float64, error) {
	tokens := tokenize(expr)
	if len(tokens) == 0 {
		return 0, fmt.Errorf("expresión vacía")
	}
	// Primero multiplicación, división, módulo
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "*" || tokens[i] == "/" || tokens[i] == "%" {
			left, _ := strconv.ParseFloat(tokens[i-1], 64)
			right, _ := strconv.ParseFloat(tokens[i+1], 64)
			var res float64
			switch tokens[i] {
			case "*":
				res = left * right
			case "/":
				if right == 0 {
					// Si se divide por cero, devolver el número que se intentó dividir
					res = left
				} else {
					res = left / right
				}
			case "%":
				if right == 0 {
					// Si se divide por cero, devolver el número que se intentó dividir
					res = left
				} else {
					res = float64(int(left) % int(right))
				}
			}
			tokens = append(tokens[:i-1], append([]string{fmt.Sprintf("%v", res)}, tokens[i+2:]...)...)
			i--
		}
	}
	// Luego suma y resta
	res, _ := strconv.ParseFloat(tokens[0], 64)
	for i := 1; i < len(tokens); i += 2 {
		op := tokens[i]
		num, _ := strconv.ParseFloat(tokens[i+1], 64)
		switch op {
		case "+":
			res += num
		case "-":
			res -= num
		}
	}
	return res, nil
}

// tokenize separa la expresión en números y operadores
func tokenize(expr string) []string {
	expr = strings.ReplaceAll(expr, " ", "")
	tokens := []string{}
	num := ""
	for _, r := range expr {
		if unicode.IsDigit(r) || r == '.' {
			num += string(r)
		} else {
			if num != "" {
				tokens = append(tokens, num)
				num = ""
			}
			tokens = append(tokens, string(r))
		}
	}
	if num != "" {
		tokens = append(tokens, num)
	}
	return tokens
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
