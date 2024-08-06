package utils

import "math/rand/v2"

func GetRandomString(stringLen int) string {
	var characterSet = []byte("1234567890     ABCDEFGHIKLMNOPQRSTVXYZ")
	data := make([]byte, 1+rand.Int()%(stringLen-1))
	rand.Shuffle(len(characterSet), func(i, j int) {
		characterSet[i], characterSet[j] = characterSet[j], characterSet[i]
	})
	copy(data, characterSet)
	data = append(data, []byte("\n")...)
	return string(data)
}
