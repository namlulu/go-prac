package dictionary

import (
	"errors"
)

type Dictionary map[string]string

type Money int

var errNotFound = errors.New("could not find the word you were looking for")
var errWordExists = errors.New("cannot add word because it already exists")
var errCannotUpdate = errors.New("cannot update word because it does not exist")
var errCannotDelete = errors.New("cannot delete word because it does not exist")

func (d Dictionary) Add(word, definition string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = definition
	case nil:
		return errWordExists
	}

	return nil
}

func (d Dictionary) Search(word string) (string, error) {
	definition, ok := d[word]
	if !ok {
		return "", errNotFound
	}
	return definition, nil
}

func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		return errCannotUpdate
	case nil:
		d[word] = definition
	}
	return nil
}

func (d Dictionary) Delete(word string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		return errCannotDelete
	case nil:
		delete(d, word)
	}
	return nil
}
