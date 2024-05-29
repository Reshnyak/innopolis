package main

// Interface Formatter
type Formatter interface {
	Format() string
}

// Структура с обычным текстом
type PlainText struct {
	text string
}

// Функция возвращает текс без изменений
func (p PlainText) Format() string {
	return p.text
}

// Структура с жирным текстом
type BoltText struct {
	text string
}

// Метод возвращает **текс**
func (b BoltText) Format() string {
	return "**" + b.text + "**"
}

// Структура с текстом код
type CodeText struct {
	text string
}

// Метод возвращает `текст`
func (c CodeText) Format() string {
	return "`" + c.text + "`"
}

// Структура с курсивом
type ItalicText struct {
	text string
}

// Метод возвращает _текс_
func (i ItalicText) Format() string {
	return "_" + i.text + "_"
}
