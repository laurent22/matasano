package charfreq

import (
	"strings"
	"fmt"
)

const (
	LETTER = 1
	BIGRAM = 2
	TRIGRAM = 3
	WORD_START = 4
	WORD_END = 5
)

type ItemFrequency struct {
	Item string
	Count int
	Score float64
}

func (this ItemFrequency) String() string {
	return fmt.Sprintf("%s:%d:%f", this.Item, this.Count, this.Score)
}

type CharFrequencies struct {
	letterData []*ItemFrequency
	bigramData []*ItemFrequency
	trigramData []*ItemFrequency
	wordEndLetters []*ItemFrequency
	wordStartLetters []*ItemFrequency
}

func normalizeScores(d []*ItemFrequency, maxScore float64) {
	var mod float64 = 0
	for i := 0; i < len(d); i++ {
		if mod == 0 {
			mod = maxScore / d[i].Score 
		}
		d[i].Score *= mod
	}
}

func NewCharFrequencies() *CharFrequencies {
	f := new(CharFrequencies)
	
	// Build some statistics for letters, bigrams and trigrams 
	
	f.letterData = []*ItemFrequency{
		{" ", 72327800, 18.74},
		{"e", 37047647, 9.60},
		{"t", 27083970, 7.02},
		{"a", 23944887, 6.21},
		{"o", 22536157, 5.84},
		{"i", 20133224, 5.22},
		{"n", 20088720, 5.21},
		{"h", 18774883, 4.87},
		{"s", 18415648, 4.77},
		{"r", 17103717, 4.43},
		{"d", 13580739, 3.52},
		{"l", 12350767, 3.20},
		{"u", 8682289, 2.25},
		{"m", 7496355, 1.94},
		{"c", 7248810, 1.88},
		{"w", 7022120, 1.82},
		{"g", 6396495, 1.66},
		{"f", 6262477, 1.62},
		{"y", 6005496, 1.56},
		{"p", 5065887, 1.31},
		{",", 4784859, 1.24},
		{".", 4680323, 1.21},
		{"b", 4594147, 1.19},
		{"k", 2853307, 0.74},
		{"v", 2745322, 0.71},
		{"\"", 2566376, 0.67},
		{"'", 1699273, 0.44},
		{"-", 1000071, 0.26},
		{"?", 469889, 0.12},
		{"x", 454572, 0.12},
		{"j", 448397, 0.12},
		{";", 311385, 0.08},
		{"!", 300580, 0.08},
		{"q", 275136, 0.07},
		{"z", 268771, 0.07},
		{":", 96752, 0.03},
		{"1", 63148, 0.02},
		{"0", 40105, 0.01},
		{")", 38729, 0.01},
		{"*", 38475, 0.01},
		{"(", 38220, 0.01},
		{"2", 36981, 0.01},
		{"`", 36256, 0.01},
		{"3", 25790, 0.01},
		{"9", 24985, 0.01},
		{"5", 21865, 0.01},
		{"4", 21181, 0.01},
		{"8", 18853, 0.00},
		{"7", 17124, 0.00},
		{"6", 17007, 0.00},
		{"/", 16757, 0.00},
		{"_", 11605, 0.00},
		{"[", 11568, 0.00},
		{"»", 11551, 0.00},
		{"]", 11535, 0.00},
		{"«", 11187, 0.00},
		{"=", 9899, 0.00},
		{"´", 8807, 0.00},
		{" ", 5326, 0.00},
		{">", 4507, 0.00},
		{"~", 4067, 0.00},
		{"<", 3995, 0.00},
		{"#", 3170, 0.00},
		{"·", 2793, 0.00},
		{"&", 2690, 0.00},
		{"{", 2258, 0.00},
		{"}", 2142, 0.00},
		{"^", 1712, 0.00},
		{"|", 1512, 0.00},
		{"\\", 1366, 0.00},
		{"@", 1354, 0.00},
		{"%", 1165, 0.00},
		{"$", 1050, 0.00},
		{"ñ", 1005, 0.00},
	}
	
	f.bigramData = []*ItemFrequency{
		{"th", 92535489, 3.882543},
		{"he", 87741289, 3.681391},
		{"in", 54433847, 2.283899},
		{"er", 51910883, 2.178042},
		{"an", 51015163, 2.140460},
		{"re", 41694599, 1.749394},
		{"nd", 37466077, 1.571977},
		{"on", 33802063, 1.418244},
		{"en", 32967758, 1.383239},
		{"at", 31830493, 1.335523},
		{"ou", 30637892, 1.285484},
		{"ed", 30406590, 1.275779},
		{"ha", 30381856, 1.274742},
		{"to", 27877259, 1.169655},
		{"or", 27434858, 1.151094},
		{"it", 27048699, 1.134891},
		{"is", 26452510, 1.109877},
		{"hi", 26033632, 1.092302},
		{"es", 26033602, 1.092301},
		{"ng", 25106109, 1.053385},
	}
	
	f.trigramData = []*ItemFrequency{
		{"the", 59623899, 3.508232},
		{"and", 27088636, 1.593878},
		{"ing", 19494469, 1.147042},
		{"her", 13977786, 0.822444},
		{"hat", 11059185, 0.650715},
		{"his", 10141992, 0.596748},
		{"tha", 10088372, 0.593593},
		{"ere", 9527535,  0.560594},
		{"for", 9438784,  0.555372},
		{"ent", 9020688,  0.530771},
		{"ion", 8607405,  0.506454},
		{"ter", 7836576,  0.461099},
		{"was", 7826182,  0.460487},
		{"you", 7430619,  0.437213},
		{"ith", 7329285,  0.431250},
		{"ver", 7320472,  0.430732},
		{"all", 7184955,  0.422758},
		{"wit", 6752112,  0.397290},
		{"thi", 6709729,  0.394796},
		{"tio", 6425262,  0.378058},	
	}
	
	f.wordStartLetters = []*ItemFrequency{
		{"t", 0, 0.1594},
		{"a", 0, 0.1550},
		{"i", 0, 0.0823},
		{"s", 0, 0.0775},
		{"o", 0, 0.0712},
		{"c", 0, 0.0597},
		{"m", 0, 0.0426},
		{"f", 0, 0.0408},
		{"p", 0, 0.0400},
		{"w", 0, 0.0382},
	}
	
	f.wordEndLetters = []*ItemFrequency{
		{"e", 0, 0.1917},
		{"s", 0, 0.1435},
		{"d", 0, 0.0923},
		{"t", 0, 0.0864},
		{"n", 0, 0.0786},
		{"y", 0, 0.0730},
		{"r", 0, 0.0693},
		{"o", 0, 0.0467},
		{"l", 0, 0.0456},
		{"f", 0, 0.0408},
	}
	
	// Normalize the scores - it's assumed that a bigram should have
	// twice as much weight as a single letter, and a trigram three times more. 
	normalizeScores(f.letterData, 10)
	normalizeScores(f.bigramData, 20)
	normalizeScores(f.trigramData, 30)
	normalizeScores(f.wordStartLetters, 20)
	normalizeScores(f.wordEndLetters, 20)
	
	return f
}

func (this CharFrequencies) GetItemScore(itemType int, item string) float64 {
	if len(item) == 0 { return 0 } 
	
	var data []*ItemFrequency
	var modifier float64 = 1
	if itemType == LETTER {
		data = this.letterData
	} else if (itemType == BIGRAM) {
		data = this.bigramData
		modifier = 6
	} else if (itemType == TRIGRAM) {
		data = this.trigramData
		modifier = 6
	} else if (itemType == WORD_START) {
		data = this.wordStartLetters
		modifier = 12
	} else if (itemType == WORD_END) {
		data = this.wordEndLetters
		modifier = 12
	}
	
	for i := 0; i < len(data); i++ {
		if data[i].Item == item {
			return data[i].Score * modifier
		}
	}
	return 0	
}

func (this CharFrequencies) ScorePlainTextByItem(s []byte, itemType int) float64 {
	var output float64 = 0
	if itemType == LETTER {
		for i := 0; i < len(s); i++ {
			output += this.GetItemScore(LETTER, strings.ToLower(string(s[i])))
		}
	}
	return output
}

func (this CharFrequencies) ScorePlainText(s []byte) float64 {
	var output float64 = 0
	for i := 0; i < len(s); i++ {
		output += this.GetItemScore(LETTER, strings.ToLower(string(s[i])))
	}
	for i := 0; i < len(s) - 1; i++ {
		bigram := strings.ToLower(string(s[i:i+2]))
		s := this.GetItemScore(BIGRAM, bigram)
		output += s
		if s > 0 { i += 1 }
	}
	for i := 0; i < len(s) - 2; i++ {
		trigram := strings.ToLower(string(s[i:i+3]))
		s := this.GetItemScore(TRIGRAM, trigram)
		output += s
		if s > 0 { i += 2 }
	}
	return output
}