package main

import (
	"fmt"
)

func main() {
	// Sample texts and patterns
	testCases := []struct {
		text    string
		pattern string
	}{
		{"ABABDABACDABABCABAB", "ABABCABAB"},
		{"AABAACAADAABAABA", "AABA"},
		{"GEEKSFORGEEKS", "GEEK"},
		{"AAAAAA", "AA"},
	}

	// Test each pattern matching algorithm
	for i, tc := range testCases {
		fmt.Printf("Test Case %d:\n", i+1)
		fmt.Printf("Text: %s\n", tc.text)
		fmt.Printf("Pattern: %s\n", tc.pattern)

		// Test naive pattern matching
		naiveResults := NaivePatternMatch(tc.text, tc.pattern)
		fmt.Printf("Naive Pattern Match: %v\n", naiveResults)

		// Test KMP algorithm
		kmpResults := KMPSearch(tc.text, tc.pattern)
		fmt.Printf("KMP Search: %v\n", kmpResults)

		// Test Rabin-Karp algorithm
		rkResults := RabinKarpSearch(tc.text, tc.pattern)
		fmt.Printf("Rabin-Karp Search: %v\n", rkResults)

		fmt.Println("------------------------------")
	}
}

// NaivePatternMatch выполняет поиск паттерна в тексте методом перебора.
// Возвращает набор всех начальных индексов, в которых найден паттерн.
func NaivePatternMatch(text, pattern string) []int {
	matches := []int{}

	// Обработка крайних случаев
	if len(pattern) == 0 || len(text) < len(pattern) {
		return matches
	}

	// Проверяем каждое возможное положение в тексте
	for i := 0; i <= len(text)-len(pattern); i++ {
		j := 0

		// Проверьте каждое возможное положение в тексте
		for j < len(pattern) && text[i+j] == pattern[j] {
			j++
		}

		// Если j достиг конца паттерна, мы нашли совпадение
		if j == len(pattern) {
			matches = append(matches, i)
		}
	}

	return matches
}

// KMPSearch реализует алгоритм Кнута-Морриса-Пратта для поиска паттерна в тексте.
// Возвращает набор всех начальных индексов, в которых найден паттерн
func KMPSearch(text, pattern string) []int {
	matches := []int{}

	// Обработка крайних случаев
	if len(pattern) == 0 || len(text) < len(pattern) {
		return matches
	}

	n := len(text)
	m := len(pattern)

	// Предварительная обработка паттерна
	lps := computeLPSArray(pattern)

	i := 0 // Индекс для текста
	j := 0 // Индекс для паттерна

	for i < n {
		// Если текущие символы совпадают, переместить оба указателя вперед
		if pattern[j] == text[i] {
			i++
			j++
		}

		//  Если найдено полное совпадение
		if j == m {
			matches = append(matches, i-j)
			// Использовать lps для сдвига шаблона для следующего совпадения
			j = lps[j-1]
		} else if i < n && pattern[j] != text[i] {
			// Несоответствие после j совпадений
			if j != 0 {
				// Использовать lps для сдвига шаблона
				j = lps[j-1]
			} else {
				// Соответствий не найдено, перейти к следующему символу в тексте
				i++
			}
		}
	}

	return matches
}
func computeLPSArray(pattern string) []int {
	m := len(pattern)
	lps := make([]int, m)

	// Длина предыдущего самого длинного префикса-суффикса
	length := 0
	i := 1

	// Цикл вычисляет lps[i] для i = 1 до m-1
	for i < m {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++
		} else {
			// Это хитрая часть
			if length != 0 {
				length = lps[length-1]
				// Примечание: здесь мы не увеличиваем i
			} else {
				lps[i] = 0
				i++
			}
		}
	}

	return lps
}

// RabinKarpSearch реализует алгоритм Рабина-Карпа для поиска паттерна в тексте.
// Возвращает набор всех начальных индексов, где найден паттерн.
func RabinKarpSearch(text, pattern string) []int {
	matches := []int{}

	// Обработка крайних случаев
	if len(pattern) == 0 || len(text) < len(pattern) {
		return matches
	}

	n := len(text)
	m := len(pattern)

	// Большое простое число для предотвращения коллизий хешей
	prime := 101

	// Базовое значение для хеш-функции
	base := 256

	// Хеш-значение для шаблона и начального окна
	patternHash := 0
	windowHash := 0

	// Наибольшая степень основания, которая нам нужна
	h := 1
	for i := 0; i < m-1; i++ {
		h = (h * base) % prime
	}

	// Вычислить начальные хеш-значения
	for i := 0; i < m; i++ {
		patternHash = (base*patternHash + int(pattern[i])) % prime
		windowHash = (base*windowHash + int(text[i])) % prime
	}

	// Поочередно перемещайте паттерны над текстом
	for i := 0; i <= n-m; i++ {
		// Проверить, совпадают ли хеш-значения
		if patternHash == windowHash {
			// Проверьте соответствие символ за символом
			match := true
			for j := 0; j < m; j++ {
				if text[i+j] != pattern[j] {
					match = false
					break
				}
			}
			if match {
				matches = append(matches, i)
			}
		}

		// Вычислить хеш-значение для следующего окна
		if i < n-m {
			windowHash = (base*(windowHash-int(text[i])*h) + int(text[i+m])) % prime

			// Убедитесь, что у нас есть только положительные значения хеша
			if windowHash < 0 {
				windowHash += prime
			}
		}
	}

	return matches
}
