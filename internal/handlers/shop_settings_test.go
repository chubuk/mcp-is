package handlers

import "testing"

func TestDiffSettings(t *testing.T) {
	item := func(id, title string) SettingsItem {
		return SettingsItem{ID: id, Title: title}
	}

	t.Run("first run: empty prev marks everything added", func(t *testing.T) {
		curr := SettingsSnapshot{
			"product_fields": {Title: "Дополнительные поля: Товары", Items: []SettingsItem{item("1", "Состав")}},
		}
		d := diffSettings(nil, curr)
		var pf *SettingsGroupDiff
		for i := range d {
			if d[i].Title == "Дополнительные поля: Товары" {
				pf = &d[i]
			}
		}
		if pf == nil {
			t.Fatal("product_fields diff missing")
		}
		if len(pf.Added) != 1 || pf.Added[0].ID != "1" {
			t.Fatalf("added = %+v, want [1]", pf.Added)
		}
		if len(pf.Removed) != 0 {
			t.Fatalf("removed = %+v, want empty", pf.Removed)
		}
	})

	t.Run("added and removed detected by id", func(t *testing.T) {
		prev := SettingsSnapshot{
			"product_fields": {Items: []SettingsItem{item("1", "Состав"), item("2", "Старое")}},
		}
		curr := SettingsSnapshot{
			"product_fields": {Items: []SettingsItem{item("1", "Состав"), item("3", "Новое")}},
		}
		d := diffSettings(prev, curr)
		pf := d[0]
		if pf.Title != "Дополнительные поля: Товары" {
			t.Fatalf("unexpected first group %q", pf.Title)
		}
		if len(pf.Added) != 1 || pf.Added[0].ID != "3" {
			t.Fatalf("added = %+v, want [3]", pf.Added)
		}
		if len(pf.Removed) != 1 || pf.Removed[0].ID != "2" {
			t.Fatalf("removed = %+v, want [2]", pf.Removed)
		}
		if pf.Total != 2 {
			t.Fatalf("total = %d, want 2", pf.Total)
		}
	})

	t.Run("no changes produce empty diff", func(t *testing.T) {
		snap := SettingsSnapshot{
			"price_kinds": {Items: []SettingsItem{item("7", "Розница")}},
		}
		d := diffSettings(snap, snap)
		for _, g := range d {
			if len(g.Added) != 0 || len(g.Removed) != 0 {
				t.Fatalf("group %q: unexpected diff %+v", g.Title, g)
			}
		}
	})

	t.Run("error group is not compared", func(t *testing.T) {
		prev := SettingsSnapshot{
			"option_names": {Items: []SettingsItem{item("1", "Цвет")}},
		}
		curr := SettingsSnapshot{
			"option_names": {Error: "boom"},
		}
		d := diffSettings(prev, curr)
		if len(d) != len(settingsGroupDefs) {
			t.Fatalf("len = %d, want %d", len(d), len(settingsGroupDefs))
		}
		for _, g := range d {
			if g.Title == "Свойства вариантов" {
				if len(g.Added) != 0 || len(g.Removed) != 0 {
					t.Fatalf("error group must not diff: %+v", g)
				}
				if g.Error != "boom" {
					t.Fatalf("error = %q, want boom", g.Error)
				}
			}
		}
	})
}
