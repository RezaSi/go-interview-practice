package challenge6

import (
	"regexp"
	"strings"
)

func CountWordFrequency(text string) map[string]int {
	// 1. 特殊处理：去掉撇号（'），将 "Let's" 变成 "Lets"
	//    必须在正则拆分前完成这一步，否则 "Let's" 会被拆成 "Let" 和 "s"
	cleanedText := strings.ReplaceAll(text, "'", "")

	// 2. 正则查找：只匹配连续的字母或数字
	//    这一步隐式地解决了分隔符问题：
	//    - "new-lines" 中的 "-" 不匹配，所以会被拆分为 "new" 和 "lines"
	//    - "dog." 中的 "." 不匹配，所以被丢弃
	re := regexp.MustCompile(`[a-zA-Z0-9]+`)
	words := re.FindAllString(cleanedText, -1)

	wordMap := make(map[string]int, len(words))
	for _, word := range words {
		// 3. 转小写并统计
		wordMap[strings.ToLower(word)]++
	}

	return wordMap
}