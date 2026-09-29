package integration

import (
	"net/http"
	"testing"
)

func TestMutateCatalog(t *testing.T) {
	user := createUser(t)

	a1 := loadAnthologyFixture(t, "anthology_a.json")
	a2 := loadAnthologyFixture(t, "anthology_b.json")

	a1ID := createAnthology(t, a1).ID
	a2ID := createAnthology(t, a2).ID

	c := getCatalog(t, user.ID)

	if len(c.Metas) != 0 {
		t.Errorf("expected empty catalog, got %v", c.Metas)
	}

	setCatalog(t, user, []string{a1ID, a2ID})

	c = getCatalog(t, user.ID)

	if len(c.Metas) != 2 {
		t.Errorf("expected 2 catalog entries, got %d", len(c.Metas))
	}

	if c.Metas[0].ID != a1ID {
		t.Errorf("expected first catalog entry %q, got %q", a1ID, c.Metas[0].ID)
	}

	if c.Metas[1].ID != a2ID {
		t.Errorf("expected second catalog entry %q, got %q", a2ID, c.Metas[1].ID)
	}

	setCatalog(t, user, []string{a2ID})

	c = getCatalog(t, user.ID)

	if len(c.Metas) != 1 {
		t.Fatalf("expected 1 catalog entry, got %d", len(c.Metas))
	}

	if c.Metas[0].ID != a2ID {
		t.Errorf("expected catalog entry %q, got %q", a2ID, c.Metas[0].ID)
	}
}

func TestEditingUnownedCatalogs(t *testing.T) {

	user1 := createUser(t)
	user2 := createUser(t)

	a := loadAnthologyFixture(t, "anthology_a.json")
	aID := createAnthology(t, a).ID

	setCatalog(t, user2, []string{aID})

	c := getCatalog(t, user2.ID)

	if len(c.Metas) != 1 {
		t.Fatalf("expected 1 catalog entry, got %v", c)
	}

	user2.Cookie = user1.Cookie

	resp := requestSetCatalog(t, user2, []string{})
	defer resp.Body.Close()

	requireStatus(t, resp, http.StatusUnauthorized)

	c = getCatalog(t, user2.ID)

	if len(c.Metas) == 0 {
		t.Fatalf("expected catalog to remain unchanged, got %v", c)
	}
}
