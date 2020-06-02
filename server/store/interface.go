package store

import (
	"fmt"
	"github.com/mmathys/acfts/common"
	"github.com/mmathys/acfts/common/funset"
	"sync"
)

type UTXOMap struct {}

const (
	DefaultMapType = TypeInsertOnly
	DefaultMapTypeString = "insertOnly"

	TypeSyncMap    = 0
	TypeInsertOnly = 1
)

var mapType = TypeInsertOnly
var syncMap sync.Map
var funSetMap *funset.FunSet

func (m *UTXOMap) SetType(_mapType int) {
	if _mapType != TypeSyncMap && _mapType != TypeInsertOnly {
		panic("Map type must be TypeSyncMap or TypeInsertOnly")
	}
	mapType = _mapType
}

func (m *UTXOMap) Init() {
	if mapType == TypeInsertOnly {
		fmt.Printf("initializing map with len = %e...\n", funset.ArrayLength)
		funSetMap = funset.NewFunSet()
	}
}

// returns true if the identifier has been doublespent. returns false if not.
func (m *UTXOMap) Store(id common.Identifier) bool {
	if mapType == -1 {
		panic("Map must be initialized first with Init() before using Store()")
	}

	if mapType == TypeSyncMap {
		_, spent := syncMap.LoadOrStore(id, true)
		return spent
	} else if mapType == TypeInsertOnly {
		inserted := funSetMap.Insert(id)
		return !inserted
	} else {
		panic("unrecognized map type")
	}
}
