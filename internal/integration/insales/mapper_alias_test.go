package insales

import "testing"

// TestMapProperty_HandleAlias — MapProperty отдаёт handle (permalink = алиас), не голый permalink.
func TestMapProperty_HandleAlias(t *testing.T) {
	resp := InSalesPropertyResponse{ID: 1, Title: "Цвет", Permalink: "cvet"}
	m := MapProperty(&resp)
	if m["handle"] != "cvet" {
		t.Errorf("handle = %v, want cvet", m["handle"])
	}
	if m["permalink"] != "cvet" {
		t.Errorf("permalink = %v, want cvet (alias)", m["permalink"])
	}
}

// TestMapCharacteristic_HandleAlias — тот же контракт для характеристик.
func TestMapCharacteristic_HandleAlias(t *testing.T) {
	resp := InSalesCharacteristicResponse{ID: 2, PropertyID: 1, Title: "Красный", Permalink: "krasnyi"}
	m := MapCharacteristic(&resp)
	if m["handle"] != "krasnyi" {
		t.Errorf("handle = %v, want krasnyi", m["handle"])
	}
	if m["permalink"] != "krasnyi" {
		t.Errorf("permalink = %v, want krasnyi (alias)", m["permalink"])
	}
}

// TestMapOptionName_HandleAlias — опции.
func TestMapOptionName_HandleAlias(t *testing.T) {
	resp := InSalesOptionNameResponse{ID: 3, Title: "Размер", Permalink: "razmer"}
	m := MapOptionName(&resp)
	if m["handle"] != "razmer" || m["permalink"] != "razmer" {
		t.Errorf("handle/permalink = %v/%v, want razmer/razmer", m["handle"], m["permalink"])
	}
}

// TestMapOptionValue_HandleAlias — значения опций.
func TestMapOptionValue_HandleAlias(t *testing.T) {
	resp := InSalesOptionValueResponse{ID: 4, OptionNameID: 3, Title: "XL", Permalink: "xl"}
	m := MapOptionValue(&resp)
	if m["handle"] != "xl" || m["permalink"] != "xl" {
		t.Errorf("handle/permalink = %v/%v, want xl/xl", m["handle"], m["permalink"])
	}
}
