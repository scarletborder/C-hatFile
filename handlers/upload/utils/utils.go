package utils

import (
	"regexp"
	"strings"
)

func Str2Tags(s string) (parts []string) {
	re := regexp.MustCompile("[,，]")
	ori_parts := re.Split(s, -1)

	var tmp_tag string
	for _, ori := range ori_parts {
		tmp_tag = strings.TrimSpace(ori)
		if tmp_tag != "" {
			parts = append(parts, tmp_tag)
		}
	}
	return
}

func Tags2Str(tags []string) string {
	// 保证标签是trim过的
	return "[" + strings.Join(tags, ",") + "]"
}
