package entities

import "gocleanarchitecture/errors"

type BlogPost struct {
	ID      string
	Title   string
	Content string
}

func (bp *BlogPost) Validate() error {
	if bp.Title == "" {
		return errors.New("title cannot be empty")
	}
	if bp.Content == "" {
		return errors.New("content cannot be empty")
	}
	return nil
}
