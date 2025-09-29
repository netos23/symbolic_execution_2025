# Домашнее задание 2: Символьные выражения и трансляция в SMT

## Цель
Создать систему символьных выражений для представления вычислений и реализовать их трансляцию в SMT-формулы для Z3 solver. Фокус на выражениях (expressions), а не на инструкциях (statements).

## ⚠️ Важно
Весь код должен быть реализован в общей кодовой базе: `../internal/`
Это задание развивает модули symbolic и translator, которые работают совместно с CFG из ДЗ1.

## Теоретические основы

### Символьные выражения
Символьные выражения представляют вычисления в абстрактном виде, где переменные могут иметь символьные (неконкретные) значения.

### SMT (Satisfiability Modulo Theories)
SMT - расширение SAT для работы с различными теориями (арифметика, массивы, строки и т.д.).

## Задания

### Задание 2.1: Проектирование иерархии символьных выражений

Создайте интерфейс и базовые типы символьных выражений:

```go
// Базовый интерфейс для всех символьных выражений
type SymbolicExpression interface {
    // Возвращает тип выражения
    Type() ExpressionType
    
    // Возвращает строковое представление
    String() string
    
    // Принимает visitor для обхода
    Accept(visitor Visitor) interface{}
}

// Типы выражений
type ExpressionType int

const (
    IntType ExpressionType = iota
    BoolType
    ArrayType
    // Добавьте другие типы по необходимости
)
```

**Требуемые типы выражений:**

1. **Переменные:**
```go
type SymbolicVariable struct {
    Name string
    Type ExpressionType
}
```

2. **Константы:**
```go
type IntConstant struct {
    Value int64
}

type BoolConstant struct {
    Value bool
}
```

3. **Арифметические операции:**
```go
type BinaryOperation struct {
    Left     SymbolicExpression
    Right    SymbolicExpression
    Operator BinaryOperator
}

type BinaryOperator int
// ADD, SUB, MUL, DIV, MOD, EQ, LT, LE, GT, GE, NE
```

4. **Логические операции:**
```go
type LogicalOperation struct {
    Operands []SymbolicExpression
    Operator LogicalOperator
}

type LogicalOperator int
// AND, OR, NOT, IMPLIES
```

### Задание 2.2: Visitor Pattern для обхода выражений

Реализуйте паттерн Visitor для работы с символьными выражениями:

```go
type Visitor interface {
    VisitVariable(expr *SymbolicVariable) interface{}
    VisitIntConstant(expr *IntConstant) interface{}
    VisitBoolConstant(expr *BoolConstant) interface{}
    VisitBinaryOperation(expr *BinaryOperation) interface{}
    VisitLogicalOperation(expr *LogicalOperation) interface{}
    // Добавьте методы для других типов выражений
}

// Пример реализации - вывод выражения в строку
type StringPrinter struct {
    result strings.Builder
}

func (sp *StringPrinter) VisitVariable(expr *SymbolicVariable) interface{} {
    // Реализация вывода переменной
}

// И т.д. для других типов
```

### Задание 2.3: Трансляция в Z3/SMT

Создайте translator для конвертации символьных выражений в Z3 формулы:

```go
type Z3Translator struct {
    ctx     *z3.Context
    solver  *z3.Solver
    vars    map[string]z3.Value  // Кэш переменных
}

func NewZ3Translator(ctx *z3.Context, solver *z3.Solver) *Z3Translator {
    return &Z3Translator{
        ctx:    ctx,
        solver: solver,
        vars:   make(map[string]z3.Value),
    }
}

func (zt *Z3Translator) VisitVariable(expr *SymbolicVariable) interface{} {
    // Если переменная уже создана, вернуть её
    if v, exists := zt.vars[expr.Name]; exists {
        return v
    }
    
    // Создать новую переменную в Z3
    switch expr.Type {
    case IntType:
        v := zt.ctx.IntConst(expr.Name)
        zt.vars[expr.Name] = v
        return v
    case BoolType:
        v := zt.ctx.BoolConst(expr.Name)
        zt.vars[expr.Name] = v
        return v
    }
    return nil
}

func (zt *Z3Translator) VisitBinaryOperation(expr *BinaryOperation) interface{} {
    left := expr.Left.Accept(zt)
    right := expr.Right.Accept(zt)
    
    // Приведение типов и выполнение операции
    switch expr.Operator {
    case ADD:
        return left.(z3.Int).Add(right.(z3.Int))
    case EQ:
        return left.(z3.Int).Eq(right.(z3.Int))
    // И т.д. для других операций
    }
    return nil
}
```

### Задание 2.4: Построение символьных выражений

Создайте компонент для построения символьных выражений из простых значений:

```go
type ExpressionBuilder struct {
    // Вспомогательные методы для создания выражений
}

func (eb *ExpressionBuilder) BuildVariable(name string, t ExpressionType) SymbolicExpression {
    return &SymbolicVariable{Name: name, Type: t}
}

func (eb *ExpressionBuilder) BuildArithmeticExpr(left, right SymbolicExpression, op BinaryOperator) SymbolicExpression {
    return &BinaryOperation{
        Left:     left,
        Right:    right,
        Operator: op,
    }
}

func (eb *ExpressionBuilder) BuildLogicalExpr(operands []SymbolicExpression, op LogicalOperator) SymbolicExpression {
    return &LogicalOperation{
        Operands: operands,
        Operator: op,
    }
}
```

### Задание 2.5: Интеграция и тестирование

Создайте тесты для демонстрации работы с выражениями:

```go
func TestExpressionTranslation(t *testing.T) {
    // 1. Создаём символьные выражения вручную
    builder := &ExpressionBuilder{}
    
    x := builder.BuildVariable("x", IntType)
    y := builder.BuildVariable("y", IntType)
    five := &IntConstant{Value: 5}
    
    // Выражение: x + y > 5
    sum := builder.BuildArithmeticExpr(x, y, ADD)
    condition := builder.BuildArithmeticExpr(sum, five, GT)
    
    // 2. Транслировать в Z3
    ctx := z3.NewContext(z3.NewConfig())
    solver := z3.NewSolver(ctx)
    translator := NewZ3Translator(ctx, solver)
    
    z3Expr := condition.Accept(translator)
    solver.Assert(z3Expr.(z3.Bool))
    
    // 3. Проверить выполнимость
    if solver.Check() == z3.True {
        model := solver.Model()
        fmt.Printf("Найдено решение: %v\n", model)
    }
}
```

## Что сдавать

**Реализация в общей кодовой базе:**
- `../internal/symbolic/expressions.go` - дополните методы для символьных выражений
- `../internal/translator/z3_translator.go` - реализуйте трансляцию выражений в Z3

**В папке homework2:**
- `main.go` - демонстрационная программа (уже создана)
- `examples/test_functions.go` - тестовые функции (уже созданы)
- `REPORT.md` - отчёт о выполненной работе

```  
homework2/
├── examples/
│   └── test_functions.go    // Функции для тестирования
├── main.go                  // Демонстрация (использует internal/*)
├── REPORT.md                // Ваш отчёт
└── README.md
```
