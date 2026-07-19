package census

import "strings"

func CanonicalizeSourceName(name string) string {
	var canonicalName strings.Builder

	canonicalName.Grow(len(name))

	separatorWritten := false

	for _, character := range name {
		if character == ' ' || character == '-' {
			if !separatorWritten {
				canonicalName.WriteByte('_')

				separatorWritten = true
			}

			continue
		}

		canonicalName.WriteRune(character)

		separatorWritten = false
	}

	return canonicalName.String()
}
