// Демонстрационная программа для тестирования CFG анализа
package main

import (
	"fmt"

	"symbolic-execution-course/internal/cfg"
)

func main() {
	fmt.Println("=== CFG Analyzer Demo ===")

	// Пример исходного кода для анализа
	_ = `
package main

func testFunction(x int) int {
	if x > 0 {
		return x * 2
	} else {
		return x * -1
	}
}
`

	// Создаём builder для CFG
	_ = cfg.NewBuilder()

	// TODO: Раскомментируйте после реализации методов в builder.go
	/*
		// Строим CFG из исходного кода
		graph, err := builder.BuildCFGFromSource(source, "testFunction")
		if err != nil {
			log.Fatalf("Ошибка построения CFG: %v", err)
		}

		fmt.Printf("CFG построен для функции с %d блоками\n", len(graph.Blocks))

		// Создаём visualizer
		visualizer := cfg.NewVisualizer()

		// Экспортируем в текстовый формат
		textOutput := visualizer.ExportToText(graph)
		fmt.Println("=== Текстовое представление CFG ===")
		fmt.Println(textOutput)

		// Экспортируем в DOT формат
		dotOutput := visualizer.ExportToDot(graph)
		fmt.Println("\n=== DOT представление CFG ===")
		fmt.Println(dotOutput)

		// Показываем статистику
		fmt.Println("\n=== Статистика ===")
		visualizer.PrintStatistics(graph)

		// Анализируем обратные рёбра
		backEdges := graph.FindBackEdges()
		fmt.Printf("Найдено обратных рёбер (циклов): %d\n", len(backEdges))

		// Анализируем достижимость
		reachable := graph.GetReachableBlocks()
		unreachable := graph.GetUnreachableBlocks()
		fmt.Printf("Достижимых блоков: %d\n", len(reachable))
		fmt.Printf("Недостижимых блоков: %d\n", len(unreachable))
	*/

	fmt.Println("Реализуйте методы в cfg пакете для запуска анализа!")
}
