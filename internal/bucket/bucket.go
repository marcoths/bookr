package bucket

var Seats = bucket{[]byte("bookr_seats")}

type bucket struct {
	name []byte
}

func (b *bucket) Name() []byte {
	return b.name
}
