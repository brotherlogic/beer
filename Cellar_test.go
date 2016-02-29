package main

import "log"
import "math"
import "testing"

type LineCounterPrinter struct {
	lcount int
}

func (lineCounter *LineCounterPrinter) Println(output string) {
	lineCounter.lcount = lineCounter.lcount + 1
}

func (lineCounter LineCounterPrinter) GetCount() int {
	return lineCounter.lcount
}

func TestPrint(t *testing.T) {
	printer := StdOutPrint{}

	//Line below should not fail
	printer.Println("Made up")
}

func TestBuildCellar(t *testing.T) {
	cellar := BuildCellar("testdata/simplecellar.cellar")
	if cellar.Size() != 3 {
		t.Errorf("Cellar has %d entries, should have 3\n", cellar.Size())
	}
}

func TestBuildCellarWithBadSize(t *testing.T) {
	cellar := BuildCellar("testdata/badsize.cellar")
	if cellar.Size() != 2 {
		t.Errorf("Cellar has %d entries, should have 2\n", cellar.Size())
	}
}

func TestBuildCellarFileOpenFail(t *testing.T) {
	cellar := BuildCellar("testdata/makdeupfilename.cellar")
	if cellar != nil {
		t.Errorf("Cellar has been built with a non-existant file\n")
	}
}

func TestComputeMixSizes(t *testing.T) {
	cellar := NewCellar("testing_cellar")
	beer1, err1 := NewBeer("1234~01/01/16~bomber")
	beer2, err2 := NewBeer("1235~01/01/16~small")

	if err1 != nil || err2 != nil {
		t.Errorf("Parse issue %v,%v\n", err1, err2)
	}

	cellar.AddBeer(beer1)

	cost := cellar.ComputeInsertCost(beer2)
	if cost >= 0 {
		t.Errorf("Beer sizes have been mixed: %v\n", cost)
	}
}

func TestComputeEmptyCellarCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer, _ := NewBeer("1234~01/01/16")

	ic := cellar.ComputeInsertCost(beer)

	if ic != math.MaxInt16 {
		t.Errorf("Insert cost should be really high here; in fact it's %d\n", ic)
	}
}

func TestComputeBackCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")

	cellar.AddBeer(beer1)
	ic := cellar.ComputeInsertCost(beer2)

	if ic != 1 {
		t.Errorf("Insert costs should have been one, in fact it's %d\n", ic)
	}
}

func TestComputeMiddleCost(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, _ := NewBeer("1234~01/03/16~bomber")

	cellar.AddBeer(beer1)
	cellar.AddBeer(beer3)
	ic := cellar.ComputeInsertCost(beer2)

	if ic != 1 {
		t.Errorf("Insert costs should have been one, in fact it's %d\n", ic)
	}
}

func TestCellarInsert(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16~bomber")
	beer2, _ := NewBeer("1234~01/02/16~bomber")
	beer3, _ := NewBeer("1234~01/03/16~bomber")

	cellar.AddBeer(beer1)
	cellar.AddBeer(beer2)
	cellar.AddBeer(beer3)

	beert1 := cellar.GetNext()
	beert2 := cellar.GetNext()
	beert3 := cellar.GetNext()

	if beert1 != beer1 {
		t.Errorf("Beer1 has been inserted in the wrong position\n")
	}
	if beert2 != beer2 {
		t.Errorf("Beer2 has been inserted in the wrong position\n")
	}
	if beert3 != beer3 {
		t.Errorf("Beer3 has been inserted in the wrong position\n")
	}
}

func TestPrintOutCellar(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, _ := NewBeer("1234~01/01/16")
	beer2, _ := NewBeer("1234~01/02/16")
	beer3, _ := NewBeer("1234~01/03/16")

	cellar.AddBeer(beer1)
	cellar.AddBeer(beer2)
	cellar.AddBeer(beer3)

	linecounter := &LineCounterPrinter{}
	cellar.PrintCellar(linecounter)

	if linecounter.GetCount() != 4 {
		t.Errorf("Print cellar has printed the wrong number of lines: %v\n", linecounter.GetCount())
	}
}

func TestSaveCellar(t *testing.T) {
	cellar := NewCellar("test_cellar")
	beer1, err1 := NewBeer("1234~01/01/16~bomber")
	beer2, err2 := NewBeer("1234~01/02/16~bomber")
	beer3, err3 := NewBeer("1234~01/03/16~bomber")

	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("Problems loading beers: %v, %v, %v\n", err1, err2, err3)
	}

	cellar.AddBeer(beer1)
	cellar.AddBeer(beer2)
	cellar.AddBeer(beer3)

	log.Printf("Pre save size: %v\n", cellar.Size())

	cellar.Save()

	cellar2 := BuildCellar("test_cellar")
	if cellar2.Size() != 3 {
		t.Errorf("Reloading cellar is not the right size: %v\n", cellar2.Size())
	}
}

func TestSaveBadCellar(t *testing.T) {
	cellar := NewCellar("madeupdirectory/blah/blah/cellar")
	cellar.Save()
}
