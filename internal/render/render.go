package render

import (
	"strings"
)

func Markdown(src string) string {
	var b strings.Builder
	for _, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(line, "# ") {
			b.WriteString(strings.ToUpper(strings.TrimPrefix(line, "# ")))
			b.WriteByte('\n')
			continue
		}
		b.WriteString(renderCodeSpans(line))
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

func renderCodeSpans(line string) string {
	var b strings.Builder
	for {
		start := strings.IndexByte(line, '`')
		if start < 0 {
			b.WriteString(line)
			break
		}
		end := strings.IndexByte(line[start+1:], '`')
		if end < 0 {
			b.WriteString(line)
			break
		}
		end += start + 1
		b.WriteString(line[:start])
		b.WriteString("[")
		b.WriteString(line[start+1 : end])
		b.WriteString("]")
		line = line[end+1:]
	}
	return b.String()
}

func HighlightGo(src string) string {
	replacer := strings.NewReplacer("func ", "«func» ", "return ", "«return» ")
	return replacer.Replace(src)
}

func Complete(prefix string, candidates []string) []string {
	var out []string
	for _, candidate := range candidates {
		if strings.HasPrefix(candidate, prefix) {
			out = append(out, candidate)
		}
	}
	return out
}
