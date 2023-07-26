package util

import (
	"github.com/rs/zerolog/log"
	"regexp"
)

var TagRegex = CreateRegexForTags()

func CreateRegexForTags() *regexp.Regexp {
	var _symbols = "·・ー_"

	var _numbers = "0-9０-９"

	var _englishLetters = "a-zA-Zａ-ｚＡ-Ｚ"

	var _japaneseLetters = "ぁ-んァ-ン一-龠"

	var _koreanLetters = "\u1100-\u11FF\uAC00-\uD7A3"

	var _spanishLetters = "áàãâéêíóôõúüçÁÀÃÂÉÊÍÓÔÕÚÜÇ"

	var _arabicLetters = "\u0621-\u064A"

	var _thaiLetters = "\u0E00-\u0E7F"

	var detectionContentLetters = _symbols +
		_numbers +
		_englishLetters +
		_japaneseLetters +
		_koreanLetters +
		_spanishLetters +
		_arabicLetters +
		_thaiLetters
	r, err := regexp.Compile("(^|\\s)([#@]([" + detectionContentLetters + "]+))")
	if err != nil {
		log.Fatal().Err(err).Msg("Regex not compiling")
	}
	return r
}
