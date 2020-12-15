package main

import (
	"fmt"
	"strconv"

	sample "test/custom/flatbuffer/MyGame/Sample"

	flatbuffers "github.com/google/flatbuffers/go"
)

func main() {
	builder := flatbuffers.NewBuilder(0)

	// Create some weapons for our Monster ("Sword" and "Axe").
	weaponOne := builder.CreateString("Sword")
	weaponTwo := builder.CreateString("Axe")

	sample.WeaponStart(builder)
	sample.WeaponAddName(builder, weaponOne)
	sample.WeaponAddDamage(builder, 3)
	sword := sample.WeaponEnd(builder)

	sample.WeaponStart(builder)
	sample.WeaponAddName(builder, weaponTwo)
	sample.WeaponAddDamage(builder, 5)
	axe := sample.WeaponEnd(builder)

	// Serialize the FlatBuffer data.
	name := builder.CreateString("Orc")

	sample.MonsterStartInventoryVector(builder, 10)
	// Note: Since we prepend the bytes, this loop iterates in reverse.
	for i := 9; i >= 0; i-- {
		builder.PrependByte(byte(i))
	}
	inv := builder.EndVector(10)

	sample.MonsterStartWeaponsVector(builder, 2)
	// Note: Since we prepend the weapons, prepend in reverse order.
	builder.PrependUOffsetT(axe)
	builder.PrependUOffsetT(sword)
	weapons := builder.EndVector(2)

	pos := sample.CreateVec3(builder, 1.0, 2.0, 3.0)

	sample.MonsterStart(builder)
	sample.MonsterAddPos(builder, pos)
	sample.MonsterAddHp(builder, 300)
	sample.MonsterAddName(builder, name)
	sample.MonsterAddInventory(builder, inv)
	sample.MonsterAddColor(builder, sample.ColorRed)
	sample.MonsterAddWeapons(builder, weapons)
	sample.MonsterAddEquippedType(builder, sample.EquipmentWeapon)
	sample.MonsterAddEquipped(builder, axe)
	orc := sample.MonsterEnd(builder)

	builder.Finish(orc)
	buf := builder.FinishedBytes()
	monster := sample.GetRootAsMonster(buf, 0)
	assert(monster.Mana() == 150, "`monster.Mana()`", strconv.Itoa(int(monster.Mana())), "150")
	assert(monster.Hp() == 300, "`monster.Hp()`", strconv.Itoa(int(monster.Hp())), "300")
	assert(string(monster.Name()) == "Orc", "`string(monster.Name())`", string(monster.Name()),
		"\"Orc\"")
	assert(monster.Color() == sample.ColorRed, "`monster.Color()`",
		strconv.Itoa(int(monster.Color())), strconv.Itoa(int(sample.ColorRed)))
	assert(monster.Pos(nil).X() == 1.0, "`monster.Pos(nil).X()`",
		strconv.FormatFloat(float64(monster.Pos(nil).X()), 'f', 1, 32), "1.0")
	assert(monster.Pos(nil).Y() == 2.0, "`monster.Pos(nil).Y()`",
		strconv.FormatFloat(float64(monster.Pos(nil).Y()), 'f', 1, 32), "2.0")
	assert(monster.Pos(nil).Z() == 3.0, "`monster.Pos(nil).Z()`",
		strconv.FormatFloat(float64(monster.Pos(nil).Z()), 'f', 1, 32), "3.0")
	for i := 0; i < monster.InventoryLength(); i++ {
		assert(monster.Inventory(i) == byte(i), "`monster.Inventory(i)`",
			strconv.Itoa(int(monster.Inventory(i))), strconv.Itoa(int(byte(i))))
	}
	expectedWeaponNames := []string{"Sword", "Axe"}
	expectedWeaponDamages := []int{3, 5}
	weapon := new(sample.Weapon) // We need a `sample.Weapon` to pass into `monster.Weapons()`
	// to capture the output of that function.
	for i := 0; i < monster.WeaponsLength(); i++ {
		if monster.Weapons(weapon, i) {
			assert(string(weapon.Name()) == expectedWeaponNames[i], "`weapon.Name()`",
				string(weapon.Name()), expectedWeaponNames[i])
			assert(int(weapon.Damage()) == expectedWeaponDamages[i],
				"`weapon.Damage()`", strconv.Itoa(int(weapon.Damage())),
				strconv.Itoa(expectedWeaponDamages[i]))
		}
	}
	assert(monster.EquippedType() == sample.EquipmentWeapon, "`monster.EquippedType()`",
		strconv.Itoa(int(monster.EquippedType())), strconv.Itoa(int(sample.EquipmentWeapon)))

	unionTable := new(flatbuffers.Table)
	if monster.Equipped(unionTable) {
		// An example of how you can appropriately convert the table depending on the
		// FlatBuffer `union` type. You could add `else if` and `else` clauses to handle
		// other FlatBuffer `union` types for this field. (Similarly, this could be
		// done in a switch statement.)
		if monster.EquippedType() == sample.EquipmentWeapon {
			unionWeapon := new(sample.Weapon)
			unionWeapon.Init(unionTable.Bytes, unionTable.Pos)

			assert(string(unionWeapon.Name()) == "Axe", "`unionWeapon.Name()`",
				string(unionWeapon.Name()), "Axe")
			assert(int(unionWeapon.Damage()) == 5, "`unionWeapon.Damage()`",
				strconv.Itoa(int(unionWeapon.Damage())), strconv.Itoa(5))
		}
	}

	fmt.Printf("The FlatBuffer was successfully created and verified!\n")
}

func assert(assertPassed bool, codeExecuted string, actualValue string, expectedValue string) {
	if assertPassed == false {
		panic("Assert failed! " + codeExecuted + " (" + actualValue +
			") was not equal to " + expectedValue + ".")
	}
}
