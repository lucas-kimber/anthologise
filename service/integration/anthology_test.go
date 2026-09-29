package integration

import "testing"

func TestCreateAndGetAnthology(t *testing.T) {
	user := createUser(t)
	want := loadAnthologyFixture(t, "anthology_a.json")

	created := createAnthology(t, want)
	got := getAnthology(t, user.ID, created.Type, created.ID)

	if got.ID != created.ID {
		t.Errorf("expected ID %q, got %q", created.ID, got.ID)
	}

	if got.Name != created.Name {
		t.Errorf("expected name %q, got %q", created.Name, got.Name)
	}
}
