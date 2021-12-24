package constobjs

import (
    "strconv"

    "github.com/Ramen2X/outrun/enums"
    "github.com/Ramen2X/outrun/obj"
)

var DefaultPlayerStateItems = func() []obj.Item {
    items := []obj.Item{}
    baseNum := 120000
    for i := range make([]byte, 11) { // 0...10
        n := baseNum + i
        s := strconv.Itoa(n)
        item := obj.NewItem(s, 0)
        if int64(n) == enums.ItemIDLaser { // free item at the start of the game
            item.Amount = 1
        }
        items = append(items, item)
    }
    return items
}()
