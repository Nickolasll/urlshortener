package domain

import "math/rand"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateSlug(size int) string {
	result := make([]byte, size)
	for i := range result {
		result[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return "/" + string(result)
}
