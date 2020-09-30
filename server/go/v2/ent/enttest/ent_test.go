package enttest

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/facebook/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
	"github.com/liov/hoper/go/v2/ent"
	"github.com/liov/hoper/go/v2/ent/car"
	"github.com/liov/hoper/go/v2/ent/user"
)

func TestEnt(t *testing.T) {
	//get.GetDB().Migrator().CreateTable(&ent.User{},&ent.Car{},&ent.Group{})
	client := Open(t, dialect.MySQL, "web:123456@tcp(192.168.1.212:3306)/test2?charset=utf8mb4&parseTime=True&loc=Local")
	ctx := context.Background()
	u, err := QueryUser(ctx, client, "钟尉")
	if err != nil {
		t.Error(err)
	}
	u, err = CreateCars(ctx, client)
	if err != nil {
		t.Error(err)
	}
	log.Println(u)
}

func CreateUser(ctx context.Context, client *ent.Client, age int8, name string) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(age).
		SetName(name).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %v", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client, name string) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.NameEQ(name)).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user: %v", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	// Create a new car with model "Tesla".
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %v", err)
	}

	// Create a new car with model "Ford".
	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %v", err)
	}
	log.Println("car was created: ", ford)

	// Create a new user, and add it the 2 cars.
	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %v", err)
	}
	log.Println("user was created: ", a8m)
	return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %v", err)
	}
	log.Println("returned cars:", cars)

	// What about filtering specific cars.
	ford, err := a8m.QueryCars().
		Where(car.ModelEQ("Ford")).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %v", err)
	}
	log.Println(ford)
	return nil
}
