# Calculadora Fyne

Calculadora de escritorio desarrollada en Go usando el framework Fyne.

## Características

- Suma, resta, multiplicación, división y residuo.
- Soporte completo de jerarquía de operaciones y paréntesis anidados.
- Permite números negativos en cualquier posición válida.
- Validación robusta de errores de sintaxis y caracteres inválidos (mensajes claros en inglés).
- Historial de operaciones compacto e interactivo.
- Entrada unificada: puedes usar tanto los botones como el teclado físico para ingresar operaciones completas, sin necesidad de enfocar el campo de entrada.
- Captura global de teclas: la aplicación responde a las teclas numéricas y operadores aunque el campo de entrada no tenga el foco.
- Display y panel de historial con fondo blanco, tamaño fijo y fuente controlada para mejor legibilidad.
- El historial muestra cada operación en una sola línea, con fuente pequeña; si la operación es muy larga, solo se muestra el resultado.
- Al hacer clic en un resultado del historial, puedes continuar una nueva operación desde ese valor (sin modificar la operación original).
- El botón `C` limpia la entrada y el resultado.
- Si se intenta dividir por cero, se muestra el número que se intentó dividir (ejemplo: `25/0 = 25`).
- Teclado organizado en formato tradicional (4x4).

## Instalación

1. Clona el repositorio:

   ```bash
   git clone https://github.com/MiguelP-Dev/calculadora-fyne.git
   cd calculadora-fyne
   ```

2. Instala las dependencias:

   ```bash
   go mod tidy
   ```

3. Ejecuta la aplicación:

   ```bash
   go run main.go
   ```

## Uso

- Ingresa operaciones completas usando los botones o el teclado físico (por ejemplo: `12+7*3/2`, `-(2+3)*4`, `2*(3+4*2)`).
- Puedes usar paréntesis y números negativos en cualquier posición válida.
- Presiona `=` o la tecla Enter para ver el resultado.
- Si la expresión es inválida, se mostrará un mensaje de error en inglés.
- El historial muestra todas las operaciones realizadas de forma compacta.
- Haz clic en cualquier resultado del historial para usarlo como punto de partida para una nueva operación (no modifica la operación original en el historial).
- El botón `C` limpia la entrada y el resultado.
- Si divides por cero, se mostrará el número original como resultado.

## Requisitos

- Go 1.18 o superior
- Dependencias de Fyne (ver documentación oficial para requisitos de sistema)

## Documentación del código

### Estructura principal

- **main.go**: Contiene toda la lógica y la interfaz gráfica de la calculadora.
- **main_test.go**: Incluye pruebas unitarias para la lógica de la calculadora.

### Estructuras y lógica

- `Operation`: Estructura que almacena una operación realizada (expresión y resultado).
- `Calculator`: Estructura que mantiene el estado de la calculadora, el historial y referencias a los widgets de la interfaz.
- Motor de evaluación:
  - `evalExpr(expr string)`: Evalúa expresiones matemáticas completas, soportando paréntesis, jerarquía de operaciones y números negativos. Devuelve error si la sintaxis es inválida o hay caracteres no permitidos.
  - Validación robusta: nunca se produce un panic por acceso fuera de rango.
  - Los espacios en blanco se ignoran, pero cualquier otro carácter inválido genera error.

### Interfaz gráfica (Fyne)

- Se utiliza un panel vertical para organizar el display, el teclado y el historial.
- El display y el historial tienen fondo blanco, tamaño fijo y fuente controlada con `canvas.Text`.
- El teclado numérico y los botones de operación están organizados en un grid tradicional (4x4).
- El historial es interactivo: al hacer clic en una operación, su resultado se usa como entrada para una nueva operación.

### Entrada y teclado

- Toda la entrada se maneja en un solo campo (`Entry`), sin placeholder ni altura excesiva.
- La app captura teclas globalmente usando `SetOnTypedRune` y `SetOnTypedKey`, permitiendo ingresar operaciones completas desde el teclado físico.
- No es necesario enfocar el campo de entrada para usar el teclado físico.

### Pruebas

- El archivo `main_test.go` contiene pruebas unitarias para:
  - Operaciones con jerarquía y paréntesis.
  - Números negativos y anidados.
  - Errores de sintaxis y caracteres inválidos.
  - División por cero (comportamiento especial).
- Los tests han sido simplificados y adaptados a la lógica moderna de la calculadora.
- Para ejecutar los tests:

  ```bash
  go test
  ```

### Extensión y personalización

- Puedes agregar más operaciones o mejorar la interfaz editando el archivo `main.go`.
- Para cambiar el estilo visual, revisa la documentación de [Fyne](https://developer.fyne.io/).

## Compilación para Windows y Linux

Puedes compilar la calculadora para diferentes sistemas operativos utilizando las herramientas de cross-compiling de Go. A continuación se muestran ejemplos para compilar en Windows y Linux desde cualquier sistema:

### Compilar para Linux

```bash
GOOS=linux GOARCH=amd64 go build -o calculadora-linux main.go
```

Esto generará un ejecutable llamado `calculadora-linux` para sistemas Linux de 64 bits.

### Compilar para Windows

```bash
GOOS=windows GOARCH=amd64 go build -o calculadora-win.exe main.go
```

Esto generará un ejecutable llamado `calculadora-win.exe` para sistemas Windows de 64 bits.

> **Nota:** Si necesitas compilar para otras arquitecturas (por ejemplo ARM) o sistemas, puedes cambiar los valores de `GOOS` y `GOARCH` según la [documentación oficial de Go](https://golang.org/doc/install/source#environment).

## Licencia

MIT - ver [LICENSE](LICENSE)

---
Desarrollado por Miguel Portillo - [github.com/MiguelP-Dev](https://github.com/MiguelP-Dev)
