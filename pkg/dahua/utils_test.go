package dahua

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_AuthParam(t *testing.T) {
	hash := (AuthParam{
		Realm:      "Login to a0c50bcd05b2f03d067e530d9bf069af",
		Random:     "1172275829",
		Encryption: "Default",
	}).HashPassword("admin", "123")

	assert.Equal(t, "2E9AD6D2DB08E0882F376A622BC76B9A", hash)
}

func Test_Timestamp(t *testing.T) {
	data := []Timestamp{
		"2023-02-06 00:00:00",
		"2023-02-06 03:09:09",
		"2023-02-06 23:59:59",
	}

	for _, date := range data {
		from, err := date.Parse(time.Local)
		assert.Nil(t, err, nil)
		assert.Equal(t, date, NewTimestamp(from, time.Local))
	}
}

func Test_extractFilePathTags(t *testing.T) {
	data := []struct {
		Path string
		Tags []string
	}{
		{
			Path: "/mnt/sd/2023-04-09/001/jpg/07/12/04[M][0@0][0][].jpg",
			Tags: []string{"M", "0@0", "0", ""},
		},
		{
			Path: "04[M][0@0][0][].jpg",
			Tags: []string{"M", "0@0", "0", ""},
		},
		{
			Path: "04M]0@0][0][].jpg",
			Tags: []string{"0", ""}},
		{
			Path: "/mnt/dvr/mmc0p2_0/2023-04-09/0/jpg/09/44/34[M][0@0][7136][0].jpg",
			Tags: []string{"M", "0@0", "7136", "0"},
		},
		{
			Path: "/mnt/dvr/mmc0p2_0/2023-04-09/0/jpg/09/44/34[M][0@0][7136][0.jpg",
			Tags: []string{"M", "0@0", "7136"},
		},
		{
			Path: "/mnt/dvr/mmc0p2_0/2023-04-09/0/jpg/09/44/34M][0@0][7136].jpg",
			Tags: []string{"0@0", "7136"},
		},
	}

	for _, d := range data {
		tags := ExtractFilePathTags(d.Path)
		assert.Equal(t, d.Tags, tags)
	}
}
