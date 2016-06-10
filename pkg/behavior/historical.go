package behavior

import (
    "fmt"
    "time"
)

/*
Historical object can generate a history of states.  Each state
  saved is referred to as a Record.  Records record what something
  was at a certain time when something happened (event).

Historical records may be queried to reconstruct timelines, provie
  feaures such as undo, etc...
*/

// Historical Types need an Type (GetType) name and a mechanism for deternining
//   whether to create a record (RecordHistory).
type Historical interface {

    // Historical currently uses the [Item's Type]_history as the
    //  name of the store for the records
    GetType() string

    // Historical items must report themselves as historical for events
    //  using this func, which returns true of false
    // Conditional historicity can be implemented in this
    //  method.  For example, we may only want to store history on
    //  Items with a certain value, such as owned by a certain user
    IsRecordableEvent(string) bool
}

// A record of an _event_ that happened on a _Historical Object_
//
//{
//  ID
//  Type string // the type of the object this history record covers
//  Event string // what happened at this point in history
//  Date time.Date // when did this happen
//  Record interface{} // the state of the document
//}
type HistoricalRecord struct {
    ID    interface{}
    Event string
    Date  time.Time
    Item  interface{}

    // User to be implemented here
}

// Record takes an event string and Historical Item and writes
//  the record in history
//  Records are written to the ItemType_history store
func (h HistoricalRecord) Record(e string, i Historical) error {
    fmt.Println("Creating record", time.Now(), i.IsRecordableEvent(e))

    h.Event = e
    h.Date = time.Now()
    h.Item = i

    return nil

}
