package storage

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Toggly/core/app/data"
)

// TestStorageCreation tests storage creation
func skipTestStorageCreation(t *testing.T) {
	url := "mongodb://developer:password123@ds131551.mlab.com:31551/toggly"

	storage, _ := NewMongoStorage(url)

	pl, err := storage.Projects("ow1").List()
	if err != nil {
		log.Print(err)
	}

	fmt.Print(pl)
}

func TestStorage2(t *testing.T) {

	url := "mongodb://developer:password123@ds131551.mlab.com:31551/toggly"

	storage, _ := NewMongoStorage(url)

	var err error

	err = storage.Projects("ow1").Save(data.Project{
		Code:        "project1",
		Description: "Project 1 description",
		RegDate:     time.Now(),
		Status:      data.ProjectStatusActive,
	})
	if err != nil {
		if _, ok := err.(*UniqueIndexError); ok {
			log.Printf("[ERROR] Unique index: %v", err)
		}
	}

	storage.Projects("ow1").Save(data.Project{
		Code:        "project2",
		Description: "Project 2 description",
		RegDate:     time.Now(),
		Status:      data.ProjectStatusActive,
	})
	if err != nil {
		if _, ok := err.(*UniqueIndexError); ok {
			log.Printf("[ERROR] Unique index: %v", err)
		}
	}
	storage.Projects("ow2").Save(data.Project{
		Code:        "project3",
		Description: "Project 3 description",
		RegDate:     time.Now(),
		Status:      data.ProjectStatusActive,
	})
	if err != nil {
		if _, ok := err.(*UniqueIndexError); ok {
			log.Printf("[ERROR] Unique index: %v", err)
		}
	}

}
