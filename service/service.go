package service

import (
	"edot-test/domain"
	"edot-test/repository"
)

type TextToJsonService struct {
	parser repository.CommentParser
}

func NewTextToJsonService(parser repository.CommentParser) *TextToJsonService {
	return &TextToJsonService{
		parser: parser,
	}
}

func (s *TextToJsonService) ConvertTextToJson(text string) []domain.ParsedContent {
	return s.parser.ParseText(text)
}
