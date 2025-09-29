// Package cfg предоставляет функции для построения CFG из SSA представления
package cfg

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/ssa"
)

// Builder отвечает за построение CFG из исходного кода Go
type Builder struct {
	fset *token.FileSet
}

// NewBuilder создаёт новый экземпляр Builder
func NewBuilder() *Builder {
	return &Builder{
		fset: token.NewFileSet(),
	}
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// ParseAndBuildSSA парсит исходный код Go и создаёт SSA представление
// Возвращает SSA программу и функцию по имени
func (b *Builder) ParseAndBuildSSA(source string, funcName string) (*ssa.Function, error) {
	// TODO: Реализовать
	// Шаги:
	// 1. Парсинг исходного кода с помощью go/parser
	// 2. Создание SSA программы
	// 3. Поиск нужной функции по имени

	// Подсказки:
	// - Используйте parser.ParseFile для парсинга
	// - Создайте packages.Config и загрузите пакет
	// - Используйте ssautil.CreateProgram для создания SSA
	// - Найдите функцию в SSA программе

	panic("не реализовано")
}

// BuildCFG строит CFG из SSA функции
func (b *Builder) BuildCFG(ssaFunc *ssa.Function) (*CFG, error) {
	if ssaFunc == nil {
		return nil, fmt.Errorf("SSA функция не может быть nil")
	}

	cfg := NewCFG(ssaFunc)

	// TODO: Реализовать построение CFG
	// Шаги:
	// 1. Проанализировать базовые блоки SSA функции (ssaFunc.Blocks)
	// 2. Создать BasicBlock для каждого SSA блока
	// 3. Установить связи между блоками на основе инструкций перехода
	// 4. Заполнить инструкции в каждом блоке

	// Подсказки:
	// - ssaFunc.Blocks содержит все базовые блоки
	// - Каждый ssa.BasicBlock содержит инструкции и информацию о переходах
	// - Обратите внимание на инструкции: If, Jump, Return, Panic
	// - Используйте методы AddBlock и AddSuccessor

	// Пока возвращаем пустой CFG, студенты должны реализовать логику
	return cfg, nil
}

// BuildCFGFromSource - удобный метод для создания CFG из исходного кода
func (b *Builder) BuildCFGFromSource(source string, funcName string) (*CFG, error) {
	// TODO: Реализовать
	// Объединить ParseAndBuildSSA и BuildCFG
	panic("не реализовано")
}
