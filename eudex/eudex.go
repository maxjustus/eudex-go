package eudex

import (
	"fmt"
	"math/bits"
)

const LETTERS uint64 = 26

var A uint64 = uint64('a')
var Z uint64 = uint64('z')

func charCode(s string, idx int) uint64 {
	return uint64(s[idx])
}

//! The raw Eudex API.

// / The sound table.
// /
// / The first bit each describes a certain property of the phone:
// /
// / | Position | Modifier | Property     | Phones                   |
// / |----------|---------:|--------------|:------------------------:|
// / | 1        | 1        | Discriminant | (for tagging duplicates) |
// / | 2        | 2        | Nasal        | mn                       |
// / | 3        | 4        | Fricative    | fvsjxzhct                |
// / | 4        | 8        | Plosive      | pbtdcgqk                 |
// / | 5        | 16       | Dental       | tdnzs                    |
// / | 6        | 32       | Liquid       | lr                       |
// / | 7        | 64       | Labial       | bfpv                     |
// / | 8        | 128      | Confident¹   | lrxzq                    |
// /
// / ¹hard to misspell.
// /
// / Vowels are, to maxize the XOR distance, represented by 0 and 1 (open and close, respectively).
var PHONES = []uint8{
	0, // a
	//    +--------- Confident
	//    |+-------- Labial
	//    ||+------- Liquid
	//    |||+------ Dental
	//    ||||+----- Plosive
	//    |||||+---- Fricative
	//    ||||||+--- Nasal
	//    |||||||+-- Discriminant
	//    ||||||||
	0b01001000, // b
	0b00001100, // c
	0b00011000, // d
	0,          // e
	0b01000100, // f
	0b00001000, // g
	0b00000100, // h
	1,          // i
	0b00000101, // j
	0b00001001, // k
	0b10100000, // l
	0b00000010, // m
	0b00010010, // n
	0,          // o
	0b01001001, // p
	0b10101000, // q
	0b10100001, // r
	0b00010100, // s
	0b00011101, // t
	1,          // u
	0b01000101, // v
	0b00000000, // w
	0b10000100, // x
	1,          // y
	0b10010100, // z
}

// Non ASCII-phones.
//
// Starts 0xDF (ß). These are all aproixmated sounds, since they can vary a lot between languages.
var PHONES_C1 = []uint8{
	PHONES[('s'-'a')] ^ 1, // ß
	0,                     // à
	0,                     // á
	0,                     // â
	0,                     // ã
	0,                     // ä [æ]
	1,                     // å [oː]
	0,                     // æ [æ]
	PHONES[('z'-'a')] ^ 1, // ç [t͡ʃ]
	1,                     // è
	1,                     // é
	1,                     // ê
	1,                     // ë
	1,                     // ì
	1,                     // í
	1,                     // î
	1,                     // ï
	0b00010101,            // ð [ð̠] (represented as a non-plosive T)
	0b00010111,            // ñ [nj] (represented as a combination of n and j)
	0,                     // ò
	0,                     // ó
	0,                     // ô
	0,                     // õ
	1,                     // ö [ø]
	255,                   // ÷
	1,                     // ø [ø]
	1,                     // ù
	1,                     // ú
	1,                     // û
	1,                     // ü
	1,                     // ý
	0b00010101,            // þ [ð̠] (represented as a non-plosive T)
	1,                     // ÿ
}

// / An _injective_ phone table.
// /
// / The table is derived the following way:
// /
// / | Position | Modifier | Property (vowel)    | Property (consonant)                              |
// / |----------|---------:|---------------------|---------------------------------------------------|
// / | 1        | 1        | Discriminant        | (property 2 from the phone table) or discriminant |
// / | 2        | 2        | Is it open-mid?     | (property 3 from the phone table)                 |
// / | 3        | 4        | Is it central?      | (property 4 from the phone table)                 |
// / | 4        | 8        | Is it close-mid?    | (property 5 from the phone table)                 |
// / | 5        | 16       | Is it front?        | (property 6 from the phone table)                 |
// / | 6        | 32       | Is it close?        | (property 7 from the phone table)                 |
// / | 7        | 64       | More close than [ɜ] | (property 8 from the phone table)                 |
// / | 8        | 128      | Vowel?                                                                  |
// /
// / If it is a consonant, the rest of the bits are simply a right truncated version of the
// / [`PHONES`](./const.PHONES.html) table, with the LSD used as discriminant.
var INJECTIVE_PHONES = []uint8{
	//    +--------- Vowel
	//    |+-------- Closer than ɜ
	//    ||+------- Close
	//    |||+------ Front
	//    ||||+----- Close-mid
	//    |||||+---- Central
	//    ||||||+--- Open-mid
	//    |||||||+-- Discriminant
	//    ||||||||   (*=vowel)
	0b10000100, // a*
	0b00100100, // b
	0b00000110, // c
	0b00001100, // d
	0b11011000, // e*
	0b00100010, // f
	0b00000100, // g
	0b00000010, // h
	0b11111000, // i*
	0b00000011, // j
	0b00000101, // k
	0b01010000, // l
	0b00000001, // m
	0b00001001, // n
	0b10010100, // o*
	0b00100101, // p
	0b01010100, // q
	0b01010001, // r
	0b00001010, // s
	0b00001110, // t
	0b11100000, // u*
	0b00100011, // v
	0b00000000, // w
	0b01000010, // x
	0b11100100, // y*
	0b01001010, // z
}

// / Non-ASCII injective phone table.
// /
// / Starting at C1.
var INJECTIVE_PHONES_C1 = []uint8{
	INJECTIVE_PHONES[('s'-'a')] ^ 1, // ß
	INJECTIVE_PHONES[('a'-'a')] ^ 1, // à
	INJECTIVE_PHONES[('a'-'a')] ^ 1, // á
	//    +--------- Vowel
	//    |+-------- Closer than ɜ
	//    ||+------- Close
	//    |||+------ Front
	//    ||||+----- Close-mid
	//    |||||+---- Central
	//    ||||||+--- Open-mid
	//    |||||||+-- Discriminant
	//    ||||||||
	0b10000000,                      // â
	0b10000110,                      // ã
	0b10100110,                      // ä [æ]
	0b11000010,                      // å [oː]
	0b10100111,                      // æ [æ]
	0b01010100,                      // ç [t͡ʃ]
	INJECTIVE_PHONES[('e'-'a')] ^ 1, // è
	INJECTIVE_PHONES[('e'-'a')] ^ 1, // é
	INJECTIVE_PHONES[('e'-'a')] ^ 1, // ê
	0b11000110,                      // ë [ə] or [œ]
	INJECTIVE_PHONES[('i'-'a')] ^ 1, // ì
	INJECTIVE_PHONES[('i'-'a')] ^ 1, // í
	INJECTIVE_PHONES[('i'-'a')] ^ 1, // î
	INJECTIVE_PHONES[('i'-'a')] ^ 1, // ï
	0b00001011,                      // ð [ð̠] (represented as a non-plosive T)
	0b00001011,                      // ñ [nj] (represented as a combination of n and j)
	INJECTIVE_PHONES[('o'-'a')] ^ 1, // ò
	INJECTIVE_PHONES[('o'-'a')] ^ 1, // ó
	INJECTIVE_PHONES[('o'-'a')] ^ 1, // ô
	INJECTIVE_PHONES[('o'-'a')] ^ 1, // õ
	0b11011100,                      // ö [œ] or [ø]
	255,                             // ÷
	0b11011101,                      // ø [œ] or [ø]
	INJECTIVE_PHONES[('u'-'a')] ^ 1, // ù
	INJECTIVE_PHONES[('u'-'a')] ^ 1, // ú
	INJECTIVE_PHONES[('u'-'a')] ^ 1, // û
	INJECTIVE_PHONES[('y'-'a')] ^ 1, // ü
	INJECTIVE_PHONES[('y'-'a')] ^ 1, // ý
	0b00001011,                      // þ [ð̠] (represented as a non-plosive T)
	INJECTIVE_PHONES[('y'-'a')] ^ 1, // ÿ
}

type EudexHash struct {
	h uint64
}

func Eudex(sequence string) EudexHash {
	if len(sequence) == 0 {
		return EudexHash{0}
	}

	entry := (charCode(sequence, 0) | 32 - A) & 0xFF
	var firstByte uint64
	if entry < LETTERS {
		firstByte = uint64(INJECTIVE_PHONES[entry])
	} else {
		if 0xDF <= entry && entry < 0xFF {
			firstByte = uint64(INJECTIVE_PHONES_C1[entry-0xDF])
		}
	}

	var res uint64
	n, b := 0, 1
	for n < 8 && b < len(sequence) {
		entry := (charCode(sequence, b) | 32 - A) & 0xFF

		var x uint64
		if entry <= Z {
			if entry < LETTERS {
				x = uint64(PHONES[entry])
			} else if 0xDF <= entry && entry < 0xFF {
				x = uint64(PHONES_C1[entry-0xDF])
			} else {
				b++
				continue
			}

			if (res & 0xFE) != (x & 0xFE) {
				res = res << 8
				res |= x
				n++
			}
		}

		b++
	}

	return EudexHash{res | (firstByte << 56)}
}

func (e EudexHash) Sub(other EudexHash) uint64 {
	return e.h ^ other.h
}

func (e EudexHash) Dist(other EudexHash) uint32 {
	d := e.Sub(other)

	return uint32(bits.OnesCount8(uint8(d))) +
		uint32(bits.OnesCount8(uint8(d>>8)))*2 +
		uint32(bits.OnesCount8(uint8(d>>16)))*3 +
		uint32(bits.OnesCount8(uint8(d>>24)))*5 +
		uint32(bits.OnesCount8(uint8(d>>32)))*8 +
		uint32(bits.OnesCount8(uint8(d>>40)))*13 +
		uint32(bits.OnesCount8(uint8(d>>48)))*21 +
		uint32(bits.OnesCount8(uint8(d>>56)))*34
}

func (e EudexHash) HammingDist(other EudexHash) uint32 {
	d := e.Sub(other)

	return uint32(bits.OnesCount64(d))
}

func (e EudexHash) String() string {
	return fmt.Sprintf("%064b", e.h)
}

func (e EudexHash) Similar(other EudexHash) bool {
	return e.Dist(other) < 15
}

func StringDistance(a, b string) uint32 {
	return Eudex(a).Dist(Eudex(b))
}

func StringHammingDistance(a, b string) uint32 {
	return Eudex(a).HammingDist(Eudex(b))
}

func Similar(a, b string) bool {
	return Eudex(a).Similar(Eudex(b))
}
