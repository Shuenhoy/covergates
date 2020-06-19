package models

import (
	"testing"

	"github.com/code-devel-cover/CodeCover/core"
)

func TestRepoFind(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	session := db.Session()
	repo := &Repo{
		Name:      "test_repo_find",
		NameSpace: "gitea",
		SCM:       string(core.Gitea),
		URL:       "http://gitea/test_repo_find",
	}
	session.Create(repo)
	if repo.ID <= 0 {
		t.Fail()
	}
	store := &RepoStore{DB: db}
	rst, err := store.Find(&core.Repo{
		Name:      "test_repo_find",
		NameSpace: "gitea",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if rst.URL != repo.URL {
		t.Fail()
	}
	rst, err = store.Find(&core.Repo{ID: repo.ID})
	if err != nil {
		t.Error(err)
		return
	}
	if rst.Name != "test_repo_find" {
		t.Fail()
	}
}

func TestRepoFinds(t *testing.T) {
	ctrl, db := getDatabaseService(t)
	defer ctrl.Finish()
	session := db.Session()
	names := []string{"a", "b", "c"}
	urls := make([]string, len(names))
	for i, name := range names {
		urls[i] = "http://gitea/finds" + name
		session.Create(&Repo{
			Name:      "finds" + name,
			NameSpace: "gitea",
			URL:       urls[i],
		})
	}
	store := &RepoStore{DB: db}
	repositories, err := store.Finds(urls...)
	if err != nil {
		t.Error(err)
		return
	}
	for i, repo := range repositories {
		if repo.Name != "finds"+names[i] {
			t.Fail()
		}
		if repo.ID <= 0 {
			t.Fail()
		}
	}
}