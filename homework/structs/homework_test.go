package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	manaShift    = 12
	healthShift  = 22
	respectShift = 8
	strengthShift
	levelShift  = 12
	houseShift  = 11
	gunShift    = 10
	familyShift = 9
	typeShift   = 7
)

const (
	manaMask       = 0b1111111111
	healthMask     = 0b1111111111
	respectMask    = 0b1111
	strengthMask   = 0b1111
	experienceMask = 0b1111
	levelMask      = 0b1111
	houseMask      = 0b1
	gunMask        = 0b1
	familyMask     = 0b1
	typeMask       = 0b11
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		for i, nameByte := range []byte(name) {
			person.PersonName[i] = nameByte
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.XCoord = int32(x)
		person.YCoord = int32(y)
		person.ZCoord = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PersonGold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PackedValues2 = uint32(mana) << manaShift
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PackedValues2 |= uint32(health) << healthShift
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PackedValues2 |= uint32(respect) << respectShift
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PackedValues2 |= uint32(strength) << strengthShift
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PackedValues2 |= uint32(experience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PackedValues1 |= uint16(level) << levelShift
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.PackedValues1 |= uint16(1) << houseShift
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.PackedValues1 |= uint16(1) << gunShift
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.PackedValues1 |= uint16(1) << familyShift
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.PackedValues1 |= uint16(personType) << typeShift
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	PersonName    [42]byte
	PackedValues1 uint16 // 4 - level, 1 - house, 1 - gun, 1 - family, 2 - type
	PackedValues2 uint32 // 10 - health, 10 - mana, 4 - respect, 4 - strength, 4 - experience
	XCoord        int32
	YCoord        int32
	ZCoord        int32
	PersonGold    uint32
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}

	for _, option := range options {
		option(&person)
	}

	return person
}

func (p *GamePerson) Name() string {
	length := 0
	for i, b := range p.PersonName {
		if b == 0 {
			break
		}
		length = i + 1
	}
	return unsafe.String(&p.PersonName[0], length)
}

func (p *GamePerson) X() int {
	return int(p.XCoord)
}

func (p *GamePerson) Y() int {
	return int(p.YCoord)
}

func (p *GamePerson) Z() int {
	return int(p.ZCoord)
}

func (p *GamePerson) Gold() int {
	return int(p.PersonGold)
}

func (p *GamePerson) Mana() int {
	return int(p.PackedValues2>>manaShift) & manaMask
}

func (p *GamePerson) Health() int {
	return int(p.PackedValues2>>healthShift) & healthMask
}

func (p *GamePerson) Respect() int {
	return int(p.PackedValues2>>respectShift) & respectMask
}

func (p *GamePerson) Strength() int {
	return int(p.PackedValues2>>strengthShift) & strengthMask
}

func (p *GamePerson) Experience() int {
	return int(p.PackedValues2) & experienceMask
}

func (p *GamePerson) Level() int {
	return int(p.PackedValues1>>levelShift) & levelMask
}

func (p *GamePerson) HasHouse() bool {
	return int(p.PackedValues1>>houseShift)&houseMask == 1
}

func (p *GamePerson) HasGun() bool {
	return int(p.PackedValues1>>gunShift)&gunMask == 1
}

func (p *GamePerson) HasFamilty() bool {
	return int(p.PackedValues1>>familyShift)&familyMask == 1
}

func (p *GamePerson) Type() int {
	return int(p.PackedValues1>>typeShift) & typeMask
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
