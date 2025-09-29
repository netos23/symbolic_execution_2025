// Package cfg предоставляет функции для визуализации CFG
package cfg

import (
	"fmt"
)

// Visualizer отвечает за экспорт CFG в различные форматы
type Visualizer struct{}

// NewVisualizer создаёт новый экземпляр Visualizer
func NewVisualizer() *Visualizer {
	return &Visualizer{}
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// ExportToDot экспортирует CFG в формат DOT для Graphviz
func (v *Visualizer) ExportToDot(cfg *CFG) string {
	// TODO: Реализовать
	// Создать DOT представление графа для визуализации в Graphviz

	// Пример структуры:
	// digraph CFG {
	//   node [shape=box];
	//   BB0 [label="Block 0\ninstruction1\ninstruction2"];
	//   BB1 [label="Block 1\ninstruction3"];
	//   BB0 -> BB1;
	// }

	// Подсказки:
	// - Используйте fmt.Sprintf для форматирования
	// - Каждый блок должен показывать свои инструкции
	// - Рёбра должны отражать переходы между блоками
	// - Выделите entry и exit блоки разными цветами

	panic("не реализовано")
}

// ExportToText создаёт текстовое представление CFG
func (v *Visualizer) ExportToText(cfg *CFG) string {
	// TODO: Реализовать
	// Создать удобочитаемое текстовое представление

	// Пример формата:
	// CFG для функции: functionName
	// Entry block: BB0
	// Exit blocks: BB3, BB4
	//
	// Block BB0:
	//   Instructions: 3
	//   Successors: BB1, BB2
	//   Predecessors: none
	//
	// Block BB1:
	//   ...

	panic("не реализовано")
}

// HighlightBackEdges выделяет обратные рёбра в DOT формате
func (v *Visualizer) HighlightBackEdges(cfg *CFG, dotString string) string {
	// TODO: Реализовать (опционально)
	// Добавить специальное выделение для обратных рёбер (циклов)
	// Например, красным цветом и пунктирной линией
	panic("не реализовано")
}

// PrintStatistics выводит статистику CFG
func (v *Visualizer) PrintStatistics(cfg *CFG) {
	// TODO: Реализовать
	// Вывести полезную статистику:
	// - Общее количество блоков
	// - Количество рёбер
	// - Количество обратных рёбер (циклов)
	// - Максимальная глубина вложенности
	// - Количество недостижимых блоков

	fmt.Printf("=== Статистика CFG ===\n")
	// TODO: добавить реализацию
	panic("не реализовано")
}
