package credential

import "strconv"

type Key string

func (record Record) Key() Key {
	return Key(record.session.String() + "#" + strconv.FormatUint(record.revision, 10) + "#" + record.handle)
}

func (key Key) String() string {
	return string(key)
}
