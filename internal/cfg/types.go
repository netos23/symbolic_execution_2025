// Package cfg предоставляет типы и функции для работы с Control Flow Graph (CFG)
package cfg

import (
	"golang.org/x/tools/go/ssa"
)

// BasicBlock представляет базовый блок в CFG
type BasicBlock struct {
	ID           int               // Уникальный идентификатор блока
	Instructions []ssa.Instruction // Инструкции в этом блоке
	Successors   []*BasicBlock     // Следующие блоки
	Predecessors []*BasicBlock     // Предыдущие блоки
	Label        string            // Метка блока (опционально)
}

// CFG представляет граф потока управления функции
type CFG struct {
	Entry    *BasicBlock   // Входной блок
	Exit     *BasicBlock   // Выходной блок (может быть nil)
	Blocks   []*BasicBlock // Все блоки в графе
	Function *ssa.Function // Исходная SSA функция
}

// NewBasicBlock создаёт новый базовый блок
func NewBasicBlock(id int) *BasicBlock {
	return &BasicBlock{
		ID:           id,
		Instructions: make([]ssa.Instruction, 0),
		Successors:   make([]*BasicBlock, 0),
		Predecessors: make([]*BasicBlock, 0),
	}
}

// AddInstruction добавляет инструкцию в базовый блок
func (bb *BasicBlock) AddInstruction(instr ssa.Instruction) {
	bb.Instructions = append(bb.Instructions, instr)
}

// AddSuccessor добавляет связь к следующему блоку
func (bb *BasicBlock) AddSuccessor(successor *BasicBlock) {
	bb.Successors = append(bb.Successors, successor)
	successor.Predecessors = append(successor.Predecessors, bb)
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// String возвращает строковое представление базового блока
// Должно включать ID блока и краткое описание инструкций
func (bb *BasicBlock) String() string {
	// TODO: Реализовать
	panic("не реализовано")
}

// IsEntry проверяет, является ли блок входным
func (bb *BasicBlock) IsEntry() bool {
	// TODO: Реализовать
	// Подсказка: входной блок не имеет предшественников
	panic("не реализовано")
}

// IsExit проверяет, является ли блок выходным
func (bb *BasicBlock) IsExit() bool {
	// TODO: Реализовать
	// Подсказка: выходной блок не имеет наследников или содержит Return
	panic("не реализовано")
}

// NewCFG создаёт новый CFG для заданной SSA функции
func NewCFG(function *ssa.Function) *CFG {
	return &CFG{
		Function: function,
		Blocks:   make([]*BasicBlock, 0),
	}
}

// AddBlock добавляет блок в CFG
func (cfg *CFG) AddBlock(block *BasicBlock) {
	cfg.Blocks = append(cfg.Blocks, block)

	// Устанавливаем входной блок, если это первый блок
	if cfg.Entry == nil {
		cfg.Entry = block
	}
}

// GetBlockByID возвращает блок по его идентификатору
func (cfg *CFG) GetBlockByID(id int) *BasicBlock {
	for _, block := range cfg.Blocks {
		if block.ID == id {
			return block
		}
	}
	return nil
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// FindBackEdges находит обратные рёбра в CFG (создающие циклы)
// Возвращает пары блоков (from, to), где edge from->to является обратным
func (cfg *CFG) FindBackEdges() [][2]*BasicBlock {
	// TODO: Реализовать
	// Подсказка: используйте DFS для обнаружения back edges
	panic("не реализовано")
}

// GetReachableBlocks возвращает все блоки, достижимые из входного блока
func (cfg *CFG) GetReachableBlocks() []*BasicBlock {
	// TODO: Реализовать
	// Подсказка: используйте BFS или DFS из entry блока
	panic("не реализовано")
}

// GetUnreachableBlocks возвращает недостижимые блоки (мёртвый код)
func (cfg *CFG) GetUnreachableBlocks() []*BasicBlock {
	// TODO: Реализовать
	panic("не реализовано")
}
