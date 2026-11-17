package utils

import "strings"

// StripRawQueryPrefix removes query segments whose key starts with prefix,
// preserving the original order and exact encoding of all other segments
// (e.g. "?wp-admin" stays "?wp-admin", not "?wp-admin=").
func StripRawQueryPrefix(rawQuery, prefix string) string {
	if rawQuery == "" || !strings.Contains(rawQuery, prefix) {
		return rawQuery
	}
	segments := strings.Split(rawQuery, "&")
	kept := segments[:0]
	for _, seg := range segments {
		key, _, _ := strings.Cut(seg, "=")
		if strings.HasPrefix(key, prefix) {
			continue
		}
		kept = append(kept, seg)
	}
	return strings.Join(kept, "&")
}
