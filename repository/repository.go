package repository

import (
	"edot-test/domain"
	"regexp"
	"strings"
)

type CommentParser interface {
	ParseText(text string) []domain.ParsedContent
}

type CommentParserImpl struct{}

func NewCommentParser() CommentParser {
	return &CommentParserImpl{}
}

func (c *CommentParserImpl) ParseText(text string) []domain.ParsedContent {
	var parsedContents []domain.ParsedContent
	var comments []string
	var currentContent domain.ParsedContent
	var currentValue domain.Values

	singleLineCommentRegex := regexp.MustCompile(`(?m)//(.*)`)
	multiLineCommentRegex := regexp.MustCompile(`(?m)/\*(.*?)\*/`)

	lines := strings.Split(text, "\n")

	isMulti := []bool{}
	isContent := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		isMultiCheck := false
		if len(isMulti) > 0 {
			isMultiCheck = isMulti[len(isMulti)-1]
		}

		if !isMultiCheck && singleLineCommentRegex.MatchString(line) {
			match := singleLineCommentRegex.FindStringSubmatch(line)
			if match[1] != "" {
				comments = append(comments, strings.TrimSpace(match[1]))
			}

			contentParts := strings.Fields(strings.TrimSuffix(line, strings.TrimSpace(match[0])))
			if len(contentParts) >= 2 {
				if contentParts[len(contentParts)-1] == "{" {
					isContent = true
					name := strings.TrimSuffix(strings.Join(contentParts[1:], " "), " {")
					currentContent.Type = contentParts[0]
					currentContent.Name = name
				} else {
					if strings.Contains(strings.Join(contentParts, " "), ":") {
						if currentContent.Type != "" {
							relParts := strings.Split(strings.Join(contentParts, " "), ":")
							if len(relParts) == 2 {
								relationName := strings.Split(strings.TrimSpace(relParts[0]), " ")
								relationValue := strings.TrimSpace(relParts[1])

								currentValue = domain.Values{
									Type:     relationName[0],
									Name:     relationName[1],
									Value:    relationValue,
									Comments: comments,
								}
								currentContent.Values = append(currentContent.Values, currentValue)

								comments = []string{}
								isMulti = []bool{}
							}
						}
					} else {
						values := strings.TrimSuffix(strings.Join(contentParts, " "), match[0])
						if strings.Contains(values, "=") {
							parts := strings.Split(strings.ReplaceAll(values, "=", ""), " ")

							currentValue = domain.Values{
								Type:     parts[0],
								Name:     parts[1],
								Value:    parts[3],
								Comments: comments,
							}

							currentContent.Values = append(currentContent.Values, currentValue)

						}
					}
				}
			}
		} else if multiLineCommentRegex.MatchString(line) {
			match := multiLineCommentRegex.FindStringSubmatch(line)
			if match[len(match)-1] != "" && !isContent {
				comments = append(comments, strings.TrimSpace(match[1]))
			} else {
				values := strings.TrimSuffix(line, match[0])
				if strings.Contains(values, "=") {
					if currentContent.Type != "" {
						currentValue = domain.Values{
							Value:    values,
							Comments: []string{strings.TrimSpace(match[1])},
						}
						currentContent.Values = append(currentContent.Values, currentValue)
					}
				}
			}
		} else if strings.Contains(line, "{") {
			isContent = true
			contentParts := strings.Fields(strings.TrimSuffix(line, "{"))
			if contentParts[len(contentParts)-1] == "{}" {
				contentParts = strings.Fields(strings.TrimSuffix(line, "{}"))

				if len(contentParts) >= 2 {
					currentContent.Type = contentParts[0]
					currentContent.Name = strings.Join(contentParts[1:], " ")
				}

				currentContent.Comments = comments
				parsedContents = append(parsedContents, currentContent)
				currentContent = domain.ParsedContent{}
				comments = []string{}
				isMulti = []bool{}
			} else if len(contentParts) >= 2 {
				currentContent.Type = contentParts[0]
				currentContent.Name = strings.Join(contentParts[1:], " ")
			}
		} else if isContent && strings.Contains(line, "}") {
			currentContent.Comments = comments
			parsedContents = append(parsedContents, currentContent)
			currentContent = domain.ParsedContent{}
			comments = []string{}
			isMulti = []bool{}
			isContent = false
		} else if strings.Contains(line, "/*") || strings.Contains(line, "/**") {
			commentsPart := strings.Fields(strings.TrimPrefix(line, "/*"))
			if strings.Contains(line, "/**") {
				commentsPart = strings.Fields(strings.TrimPrefix(line, "/**"))
			}

			comment := strings.Join(commentsPart[0:], " ")
			comment = strings.ReplaceAll(comment, "/*", "")
			comment = strings.ReplaceAll(comment, "*", "")
			comment = strings.TrimSpace(comment)
			isStar := strings.Contains(comment, "*")

			if len(commentsPart) > 0 && (!isStar) {
				if isMultiCheck {
					comment = strings.Join([]string{comments[len(comments)-1], comment}, "\n")
					comments[len(comments)-1] = comment
				} else {
					comments = append(comments, comment)
				}
			}

			isMulti = append(isMulti, true)
		} else if isMultiCheck && strings.Contains(line, "*/") {
			isMulti = append(isMulti, false)
		} else if isMultiCheck {
			commentsPart := strings.Fields(strings.TrimPrefix(line, "*"))
			comment := strings.Join(commentsPart[0:], " ")
			comment = strings.ReplaceAll(comment, "/*", "")
			comment = strings.ReplaceAll(comment, "*", "")
			comment = strings.TrimSpace(comment)
			isStar := strings.Contains(comment, "*")

			if len(commentsPart) > 0 && (!isStar) {
				if isMultiCheck && len(comments) > 0 {
					if len(isMulti) > 2 {
						if !isMulti[len(isMulti)-2] {
							comments = append(comments, comment)
						}
					} else {
						comment = strings.Join([]string{comments[len(comments)-1], comment}, "\n")
						comments[len(comments)-1] = comment
					}
				} else {
					comments = append(comments, comment)
				}
			}
		}
	}

	if currentContent.Type != "" {
		currentContent.Comments = comments
		parsedContents = append(parsedContents, currentContent)
	}

	return parsedContents
}
