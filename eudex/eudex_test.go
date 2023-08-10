package eudex

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExact(t *testing.T) {
	assert.Equal(t, Eudex("JAva"), Eudex("jAva"))
	assert.Equal(t, Eudex("co!mputer"), Eudex("computer"))
	assert.Equal(t, Eudex("comp-uter"), Eudex("computer"))
	assert.Equal(t, Eudex("comp@u#te?r"), Eudex("computer"))
	assert.Equal(t, Eudex("lal"), Eudex("lel"))
	assert.Equal(t, Eudex("rindom"), Eudex("ryndom"))
	assert.Equal(t, Eudex("riiiindom"), Eudex("ryyyyyndom"))
	assert.Equal(t, Eudex("riyiyiiindom"), Eudex("ryyyyyndom"))
	assert.Equal(t, Eudex("triggered"), Eudex("TRIGGERED"))
	assert.Equal(t, Eudex("repert"), Eudex("ropert"))
}

func TestMismatch(t *testing.T) {
	assert.NotEqual(t, Eudex("reddit"), Eudex("eddit"))
	assert.NotEqual(t, Eudex("lol"), Eudex("lulz"))
	assert.NotEqual(t, Eudex("ijava"), Eudex("java"))
	// - this fails. Interestingly this also fails for the JS version
	// assert.NotEqual(t, Eudex("jiva").String(), Eudex("java").String())
	assert.NotEqual(t, Eudex("jesus"), Eudex("iesus"))
	assert.NotEqual(t, Eudex("aesus"), Eudex("iesus"))
	assert.NotEqual(t, Eudex("iesus"), Eudex("yesus"))
	assert.NotEqual(t, Eudex("rupirt"), Eudex("ropert"))
	assert.NotEqual(t, Eudex("ripert"), Eudex("ropyrt"))
	assert.NotEqual(t, Eudex("rrr"), Eudex("rraaaa"))
	assert.NotEqual(t, Eudex("randomal"), Eudex("randomai"))
}

func TestDistance(t *testing.T) {
	assert.Greater(t, StringDistance("lizzard", "wizzard"), StringDistance("rick", "rolled"))
	assert.GreaterOrEqual(t, StringDistance("bannana", "panana"), StringDistance("apple", "abple"))
	assert.Less(t, StringDistance("trump", "drumpf"), StringDistance("gangam", "style"))
}

func TestReflexivity(t *testing.T) {
	assert.Equal(t, StringDistance("a", "b"), StringDistance("b", "a"))
	assert.Equal(t, StringDistance("youtube", "facebook"), StringDistance("facebook", "youtube"))
	assert.Equal(t, StringDistance("Rust", "Go"), StringDistance("Go", "Rust"))
	assert.Equal(t, StringDistance("rick", "rolled"), StringDistance("rolled", "rick"))
}

func TestSimilar(t *testing.T) {
	assert.True(t, Similar("yay", "yuy"))
	assert.Less(t, StringHammingDistance("crack", "crakk"), uint32(10))
	assert.True(t, Similar("what", "wat"))
	assert.True(t, Similar("jesus", "jeuses"))
	assert.True(t, Similar("", ""))
	assert.True(t, Similar("jumpo", "jumbo"))
	assert.True(t, Similar("lol", "lulz"))
	assert.True(t, Similar("goth", "god"))
	assert.True(t, Similar("maier", "meyer"))
	assert.True(t, Similar("java", "jiva"))
	assert.True(t, Similar("möier", "meyer"))
	assert.True(t, Similar("fümlaut", "fymlaut"))
	assert.Less(t, StringHammingDistance("schmid", "schmidt"), uint32(14))

	assert.False(t, Similar("youtube", "reddit"))
	assert.False(t, Similar("yet", "vet"))
	assert.False(t, Similar("hacker", "4chan"))
	assert.False(t, Similar("awesome", "me"))
	assert.False(t, Similar("prisco", "vkisco"))
	assert.False(t, Similar("no", "go"))
	assert.False(t, Similar("horse", "norse"))
	assert.False(t, Similar("nice", "mice"))
}

func BenchmarkDict(b *testing.B) {
	var hashes []EudexHash
	for i := 0; i < b.N; i++ {
		hashes = nil
		file, _ := os.Open("/usr/share/dict/web2")
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			hashes = append(hashes, Eudex(line))
		}
	}
}
