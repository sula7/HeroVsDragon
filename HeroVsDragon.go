/* Author - sula7 (@sulafpv)
Please read a manual on GitHub on how to run this application
*/

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var isGameStart bool = false
var isGameEnd bool = false

var selectedMenuLang bool = true //true = RU, false = EN
var inputMainMenuItem string     //input
var inputLangMenuItem string     //input
var hpHero int = 100
var hpDragon int = 100
var weaponHero int = 0 //input

const harmOfSword = 10
const harmOfArrow = 15
const harmOfFrstn = 30

var heroName string //input

type menuData struct { //Тексты главного меню
	point1, point2, point3, langRU, langEN, loading, bye, incorrectInp string
}

type gameDataText struct { //Игровые тексты
	hp, dragonMiss, harmHeroToDragon, harmDragonToHero, hero, dragon, entHeroName,
	gameOver, selWeapon, weapon1, weapon2, weapon3, winner, step, standoff string
}

type gameDataLogic struct { //Игровая логика

}

func setGetMainMenuTextRU() menuData { //Тексты главного меню
	mainMenu := menuData{
		point1:       "1. Начать новую игру",
		point2:       "2. Выбрать язык",
		point3:       "3. Выход",
		langRU:       "1. Русский",
		langEN:       "2. English",
		loading:      "Загружается...",
		bye:          "До скорой встречи!",
		incorrectInp: "Неверный выбор, повторите снова",
	}
	return mainMenu
}

func setGetMainMenuTextEN() menuData { //Тексты главного меню
	mainMenu := menuData{
		point1:       "1. Start a new game",
		point2:       "2. Choose lang",
		point3:       "3. Exit",
		loading:      "Loading...",
		bye:          "Good bye :)",
		incorrectInp: "Incorrect selection, try again",
	}
	return mainMenu
}

func setGetGameDataTextRU() gameDataText { //Игровые тексты
	data := gameDataText{
		hp:               "hp",
		dragonMiss:       "Дракон промахнулся и не нанёс вам урона :)",
		harmHeroToDragon: "Вы нанесли урон Дракону",
		harmDragonToHero: "Дракон нанёс Вам урон",
		hero:             "Герой",
		dragon:           "Дракон Драконыч",
		entHeroName:      "Введите имя Героя: (нажмите Enter для случайного)",
		gameOver:         "Игра завершилась",
		selWeapon:        "Выберите оружие:",
		weapon1:          "1. Меч",
		weapon2:          "2. Стрела",
		weapon3:          "3. Огненный камень",
		winner:           "Победил",
		step:             "Ход #",
		standoff:         "Победила дружба :)",
	}
	return data
}

func setGetGameDataTextEN() gameDataText { //Игровые тексты
	data := gameDataText{
		hp:               "hp",
		dragonMiss:       "Dragon missed, no harm to Hero :)",
		harmHeroToDragon: "You damaged Dragon to",
		harmDragonToHero: "Dragon damaged You to",
		hero:             "Hero",
		dragon:           "Dragon Dragoner",
		entHeroName:      "Enter Hero name: (press Enter to generate random)",
		gameOver:         "Game over",
		selWeapon:        "Select weapon:",
		weapon1:          "1. Sword",
		weapon2:          "2. Arrow",
		weapon3:          "3. Firestone",
		winner:           "Winner is",
		step:             "Step #",
		standoff:         "Friendship won :)",
	}
	return data
}

func setGetGameDataLogic() gameDataLogic { //Игровая логика
	data := gameDataLogic{}
	return data
}

func showLangMenu() {
	fmt.Println(setGetMainMenuTextRU().langRU)
	fmt.Println(setGetMainMenuTextRU().langEN)
}

func showMainMenu() {
	fmt.Println(checkLangShowText(setGetMainMenuTextRU().point1, setGetMainMenuTextEN().point1))
	fmt.Println(checkLangShowText(setGetMainMenuTextRU().point2, setGetMainMenuTextEN().point2))
	fmt.Println(checkLangShowText(setGetMainMenuTextRU().point3, setGetMainMenuTextEN().point3))
}

func selectLangMenuItem() {
	someThing := bufio.NewScanner(os.Stdin)
	someThing.Scan()
	inputLangMenuItem = strings.TrimSpace(someThing.Text()) //Убирает пробелы в начале и в конце
	switch inputLangMenuItem {
	case "1":
		selectedMenuLang = true
	case "2":
		selectedMenuLang = false
	default:
		fmt.Println(checkLangShowText(setGetMainMenuTextRU().incorrectInp, setGetMainMenuTextEN().incorrectInp))
		showLangMenu()
		selectLangMenuItem()
	}
}

func checkLangShowText(textRU, textEN string) string {
	if selectedMenuLang == true {
		return textRU
	}
	return textEN
}

func selectMainMenuItem() {
	someThing := bufio.NewScanner(os.Stdin)
	someThing.Scan()
	inputMainMenuItem = strings.TrimSpace(someThing.Text()) //Убирает пробелы в начале и в конце
	switch inputMainMenuItem {
	case "1": //Начать новую игру
		fmt.Println(checkLangShowText(setGetMainMenuTextRU().loading, setGetMainMenuTextEN().loading))
		isGameStart = true
	case "2": //Выбрать язык
		showLangMenu()
		selectLangMenuItem()
	case "3": //Выход
		fmt.Println(checkLangShowText(setGetMainMenuTextRU().bye, setGetMainMenuTextEN().bye))
		os.Exit(0)
	default:
		fmt.Println(checkLangShowText(setGetMainMenuTextRU().incorrectInp, setGetMainMenuTextEN().incorrectInp))
		selectMainMenuItem()
	}
}

func gameStart() {
	fmt.Println(checkLangShowText(setGetGameDataTextRU().entHeroName, setGetGameDataTextEN().entHeroName))
	inputHeroName()
	for i := 0; ; i++ {
		if !isGameEnd {
			showGameResult()
			showStep(i)
			showWeaponHero()
			selectWeapon()
			attackToDragon()
			checkCurrentHp()
		} else if isGameEnd {
			gameEnd()
			showWinner()
			break
		}
	}
}

func inputHeroName() {
	someThing := bufio.NewScanner(os.Stdin)
	someThing.Scan()
	if someThing.Text() == `` { //для обработки пустой строки
		data := map[string]string{}
		resp, err := http.Get("https://uinames.com/api/?amount=1&gender=male&region=kyrgyz+republic")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &data)
		heroName = data["name"]
	} else {
		heroName = strings.TrimSpace(someThing.Text()) //Убирает пробелы в начале и в конце
	}
	fmt.Println("")
}

func showGameResult() {
	fmt.Println(checkLangShowText(setGetGameDataTextRU().hero, setGetGameDataTextEN().hero), heroName,
		"\t\t\t", checkLangShowText(setGetGameDataTextRU().dragon, setGetGameDataTextEN().dragon))
	showCurrentHP()
}

func showCurrentHP() {
	fmt.Println(hpHero, setGetGameDataTextRU().hp, "\t\t\t\t", hpDragon, setGetGameDataTextEN().hp)
}

func checkCurrentHp() {
	if hpHero|hpDragon < 1 {
		isGameEnd = true
	}
}

func showStep(step int) {
	fmt.Println(checkLangShowText(setGetGameDataTextRU().step, setGetGameDataTextEN().step), step+1)
}

func gameEnd() {
	fmt.Println("\n\n")
	fmt.Println(checkLangShowText(setGetGameDataTextRU().gameOver, setGetGameDataTextEN().gameOver))
	showGameResult()
}

func showWinner() {
	if hpHero > hpDragon {
		fmt.Println("")
		fmt.Println(checkLangShowText(setGetGameDataTextRU().winner, setGetGameDataTextEN().winner),
			checkLangShowText(setGetGameDataTextRU().hero, setGetGameDataTextEN().hero), heroName)
	} else if hpDragon > hpHero {
		fmt.Println(checkLangShowText(setGetGameDataTextRU().winner, setGetGameDataTextEN().winner),
			checkLangShowText(setGetGameDataTextRU().dragon, setGetGameDataTextEN().dragon))
	} else if hpDragon == hpHero {
		fmt.Println(checkLangShowText(setGetGameDataTextRU().standoff, setGetGameDataTextEN().standoff))
	}
}

func showWeaponHero() {
	fmt.Println("\n")
	fmt.Println(checkLangShowText(setGetGameDataTextRU().selWeapon, setGetGameDataTextEN().selWeapon))
	weaponHero := [3]string{
		checkLangShowText(setGetGameDataTextRU().weapon1, setGetGameDataTextEN().weapon1),
		checkLangShowText(setGetGameDataTextRU().weapon2, setGetGameDataTextEN().weapon2),
		checkLangShowText(setGetGameDataTextRU().weapon3, setGetGameDataTextEN().weapon3),
	}
	for i := 0; i < len(weaponHero); i++ {
		fmt.Println(weaponHero[i])
	}
}

func selectWeapon() {
	fmt.Scan(&weaponHero)
	fmt.Println("\n")
}

func randomize(min int, max int) int {
	return min + rand.Intn(max-min)
}

func attackToDragon() {
	rand.Seed(time.Now().UTC().UnixNano())
	switch weaponHero {
	case 1:
		randomized := randomize(0, 20)
		casesAttackToDragon(randomized, harmOfSword)
	case 2:
		randomized := randomize(10, 30)
		casesAttackToDragon(randomized, harmOfArrow)
	case 3:
		randomized := randomize(20, 40)
		casesAttackToDragon(randomized, harmOfFrstn)
	default:
		fmt.Println(checkLangShowText(setGetMainMenuTextRU().incorrectInp, setGetMainMenuTextEN().incorrectInp))
	}
}

func casesAttackToDragon(randomized int, harm int) {
	hpHero = hpHero - randomized
	hpDragon = hpDragon - harm
	fmt.Println(checkLangShowText(setGetGameDataTextRU().harmHeroToDragon, setGetGameDataTextEN().harmHeroToDragon), harm)
	if randomized == 0 {
		fmt.Println(checkLangShowText(setGetGameDataTextRU().dragonMiss, setGetGameDataTextEN().dragonMiss))
	} else {
		fmt.Println(checkLangShowText(setGetGameDataTextRU().harmDragonToHero, setGetGameDataTextEN().harmDragonToHero), randomized)
	}
	fmt.Println("\n")
}

func main() {
	for {
		showMainMenu()
		selectMainMenuItem()
		if isGameStart {
			break
		}
	}
	gameStart()
}
