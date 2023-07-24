package util

import (
	"fmt"
	"github.com/google/wire"
)

// Monster 怪兽
type Monster struct {
	Name string
}

// NewMonster 创建怪兽
func NewMonster() Monster {
	return Monster{Name: "kitty"}
}

type Player struct {
	Name string
}

// NewPlayer 创建勇士
func NewPlayer(name string) (Player, error) {
	return Player{Name: name}, nil
}

// Mission 勇士打怪兽
type Mission struct {
	Player  Player
	Monster Monster
}

// NewMission 创建一个勇士打怪兽的任务
func NewMission(p Player, m Monster) Mission {
	return Mission{p, m}
}

func (m Mission) Start() {
	fmt.Printf("勇士 %v 打怪兽 %v\n", m.Player.Name, m.Monster.Name)
}

var MPset = wire.NewSet(NewMonster, NewPlayer)

type EndingA struct {
	Player  Player
	Monster Monster
}

func (e EndingA) Appear() {
	fmt.Printf("%s defeats %s, world peace!\n", e.Player.Name, e.Monster.Name)
}

type EndingB struct {
	Player  Player
	Monster Monster
}

func (e EndingB) Appear() {
	fmt.Printf("%s defeats %s, but become monster, world darker!\n", e.Player.Name, e.Monster.Name)
}

type Person struct {
	Name string
	Age  int
}

type PersonI interface {
	GetName() string
}

func (p Person) GetName() string {
	return p.Name
}

var PSet = wire.NewSet(wire.Struct(new(Person), "*"), wire.Bind(new(PersonI), new(Person)))

/**
var EndingASet = wire.NewSet(MPset, wire.Struct(new(EndingA), "*"))
var EndingBSet = wire.NewSet(MPset, wire.Struct(new(EndingB), "*"))
var MissionSet = wire.NewSet(wire.Struct(new(Mission), "*"))
*/
