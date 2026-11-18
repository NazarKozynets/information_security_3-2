package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var alphabetMap = make(map[rune]int)
var alphabet = []rune("АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯабвгґдеєжзиіїйклмнопрстуфхцчшщьюя")
var homophones = map[rune][]string{
	'А': {"25", "62"},
	'Б': {"99", "36"},
	'В': {"73", "10"},
	'Г': {"47", "84"},
	'Ґ': {"21", "58"},
	'Д': {"95", "32"},
	'Е': {"69", "06"},
	'Є': {"43", "80"},
	'Ж': {"17", "54", "91", "28", "65"},
	'З': {"02", "39"},
	'И': {"76", "13", "50", "87", "24"},
	'І': {"61", "98"},
	'Ї': {"35", "72"},
	'Й': {"09", "46"},
	'К': {"83", "20", "57", "94", "31"},
	'Л': {"68", "05"},
	'М': {"42", "79"},
	'Н': {"16", "53"},
	'О': {"90", "27", "64", "01", "38"},
	'П': {"75", "12", "49", "86", "23"},
	'Р': {"60", "97"},
	'С': {"34", "71"},
	'Т': {"08", "45", "82", "19"},
	'У': {"56", "93", "30", "67"},
	'Ф': {"04", "41", "78", "15"},
	'Х': {"52", "89", "26", "63", "00"},
	'Ц': {"37", "74", "11", "48", "85"},
	'Ч': {"22", "59"},
	'Ш': {"96", "33", "70", "07"},
	'Щ': {"44", "81", "18", "55"},
	'Ь': {"92", "29", "66", "03", "40"},
	'Ю': {"77", "14"},
	'Я': {"51", "88"},
}

func main() {
	const key = 8

	for i, r := range alphabet {
		alphabetMap[r] = i
	}

	input := bufio.NewReader(os.Stdin)
	fmt.Println("1 — Шифрування")
	fmt.Println("2 — Розшифрування")
	fmt.Println("Виберіть режим (1 або 2):")
	choice, ok := readLine(input)
	if !ok {
		return
	}
	mode := strings.TrimSpace(choice)

	switch mode {
	case "1":
		fmt.Println("Введіть текст:")
	case "2":
		fmt.Println("Введіть шифротекст:")
	default:
		fmt.Fprintln(os.Stderr, "Некоректний режим. Виберіть 1 або 2.")
		return
	}

	text, ok := readLine(input)
	if !ok {
		return
	}

	if mode == "1" {
		fmt.Println("Зашифрований текст:", homophonicEncrypt(caesar(text, key)))
		return
	}

	shifted, err := homophonicDecrypt(text)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println("Розшифрований текст:", caesar(shifted, -key))
}

func readLine(input *bufio.Reader) (string, bool) {
	line, err := input.ReadString('\n')

	if err == io.EOF && line == "" {
		return "", false
	}
	if err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, "Помилка читання:", err)
		return "", false
	}

	line = strings.TrimSuffix(line, "\n")
	return strings.TrimSuffix(line, "\r"), true
}

func caesar(str string, key int) string {
	var builder strings.Builder
	alphabetSize := len(alphabet) / 2

	for _, r := range str {
		alphabetIndex := findRuneIndex(r)
		if alphabetIndex < 0 {
			builder.WriteRune(r)
			continue
		}

		// Великі та малі літери обертаються в межах своєї частини абетки.
		offset := (alphabetIndex / alphabetSize) * alphabetSize
		newIndex := (alphabetIndex - offset + key) % alphabetSize
		if newIndex < 0 {
			newIndex += alphabetSize
		}

		builder.WriteRune(alphabet[offset+newIndex])
	}

	return builder.String()
}

func homophonicEncrypt(s string) string {
	occurrences := make(map[rune]int)
	var tokens []string

	for _, r := range s {
		letter := unicode.ToUpper(r)
		codes, ok := homophones[letter]

		if !ok || (r != letter && r != unicode.ToLower(letter)) {
			tokens = append(tokens, fmt.Sprintf("U+%04X", r))
			continue
		}

		code := codes[occurrences[letter]%len(codes)]
		occurrences[letter]++

		if r != letter {
			code = "l" + code
		}

		tokens = append(tokens, code)
	}

	return strings.Join(tokens, " ")
}

func homophonicDecrypt(s string) (string, error) {
	lettersByCode := make(map[string]rune)
	for letter, codes := range homophones {
		for _, code := range codes {
			lettersByCode[code] = letter
			lettersByCode["l"+code] = unicode.ToLower(letter)
		}
	}

	var result strings.Builder
	for _, token := range strings.Fields(s) {
		if letter, ok := lettersByCode[token]; ok {
			result.WriteRune(letter)
			continue
		}

		hex, ok := strings.CutPrefix(token, "U+")
		if ok && len(hex) >= 4 && len(hex) <= 6 && strings.Trim(hex, "0123456789ABCDEF") == "" {
			value, err := strconv.ParseUint(hex, 16, 32)
			if err == nil && utf8.ValidRune(rune(value)) {
				result.WriteRune(rune(value))
				continue
			}
		}

		return "", fmt.Errorf("Некоректний токен: %s", token)
	}

	return result.String(), nil
}

// Для пошуку індексу букви в алфавіті.
func findRuneIndex(target rune) int {
	if index, exists := alphabetMap[target]; exists {
		return index
	}
	return -1
}
