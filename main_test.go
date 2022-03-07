package main

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestFindFormTokenWhenItIsMissing(t *testing.T) {
	got, err := findFormToken(`no token here`)

	if err == nil {
		t.Fatalf(`Error should be returned instead of: %v`, got)
	}
}

func TestFindFormTokenWhenItExists(t *testing.T) {
	content, err := ioutil.ReadFile(`test-resources/step1.html`)
	if err != nil {
		t.Fatal(err)
	}

	got, err := findFormToken(string(content))
	want := `db765f2dc5f4c2070a27d8284852bf1d`

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestFindTrainingsNameWhenGymIsMissing(t *testing.T) {
	content, err := ioutil.ReadFile(`test-resources/step1.html`)
	if err != nil {
		t.Fatal(err)
	}

	got, err := findTrainings(`Invalid Name`, string(content))

	if err == nil {
		t.Fatalf(`Error should be returned instead of: %v`, got)
	}
}

func TestFindTrainingsWhenGymExists(t *testing.T) {
	content, err := ioutil.ReadFile(`test-resources/step1.html`)
	if err != nil {
		t.Fatal(err)
	}

	got, err := findTrainings(`FIT STAR Berlin-Moabit`, string(content))
	want := map[string]string{`YOGA`: `19`, `Les Mills BODYBALANCE`: `20`, `ZUMBA`: `21`, `Les Mills BODYPUMP`: `22`, `Les Mills BODYATTACK`: `23`, `PILATES`: `25`}

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf(`got %v want %v`, got, want)
	}
}
