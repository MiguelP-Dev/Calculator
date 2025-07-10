package main

import (
	"testing"

	"fyne.io/fyne/v2/widget"
)

func TestCalculatorBasicOps(t *testing.T) {
	calc := NewCalculator()
	calc.resultLabel = widget.NewLabel("")

	// Suma
	calc.currentInput = "5"
	calc.setOperator("+")
	calc.appendInput("3")
	calc.calculate()
	if calc.currentInput != "8" {
		t.Errorf("Suma: esperado 8, obtenido %s", calc.currentInput)
	}

	// Resta
	calc.setOperator("-")
	calc.appendInput("2")
	calc.calculate()
	if calc.currentInput != "6" {
		t.Errorf("Resta: esperado 6, obtenido %s", calc.currentInput)
	}

	// Multiplicación
	calc.setOperator("*")
	calc.appendInput("4")
	calc.calculate()
	if calc.currentInput != "24" {
		t.Errorf("Multiplicación: esperado 24, obtenido %s", calc.currentInput)
	}

	// División
	calc.setOperator("/")
	calc.appendInput("6")
	calc.calculate()
	if calc.currentInput != "4" {
		t.Errorf("División: esperado 4, obtenido %s", calc.currentInput)
	}

	// Módulo
	calc.setOperator("%")
	calc.appendInput("3")
	calc.calculate()
	if calc.currentInput != "1" {
		t.Errorf("Módulo: esperado 1, obtenido %s", calc.currentInput)
	}
}

func TestCalculatorDivisionByZero(t *testing.T) {
	calc := NewCalculator()
	calc.resultLabel = widget.NewLabel("")
	calc.currentInput = "5"
	calc.setOperator("/")
	calc.appendInput("0")
	calc.calculate()
	if calc.resultLabel.Text != "Div/0" {
		t.Errorf("División por cero: esperado 'Div/0', obtenido %s", calc.resultLabel.Text)
	}
}

func TestCalculatorHistory(t *testing.T) {
	calc := NewCalculator()
	calc.resultLabel = widget.NewLabel("")
	calc.currentInput = "2"
	calc.setOperator("+")
	calc.appendInput("2")
	calc.calculate() // 2 + 2 = 4

	calc.setOperator("*")
	calc.appendInput("3")
	calc.calculate() // 4 * 3 = 12

	if len(calc.history) != 2 {
		t.Errorf("Historial: esperado 2 operaciones, obtenido %d", len(calc.history))
	}
	if calc.history[0].Result != "12" || calc.history[1].Result != "4" {
		t.Errorf("Historial: resultados incorrectos: %+v", calc.history)
	}
}

func TestContinueFromHistory(t *testing.T) {
	calc := NewCalculator()
	calc.resultLabel = widget.NewLabel("")
	calc.currentInput = "7"
	calc.setOperator("+")
	calc.appendInput("8")
	calc.calculate() // 7 + 8 = 15

	calc.currentInput = "3"
	calc.setOperator("*")
	calc.appendInput("2")
	calc.calculate() // 3 * 2 = 6

	calc.continueFromHistory(1) // Debería poner 15 en currentInput
	if calc.currentInput != "15" {
		t.Errorf("Continuar desde historial: esperado 15, obtenido %s", calc.currentInput)
	}
}
