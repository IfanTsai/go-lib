package set_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/IfanTsai/go-lib/set"
)

func TestSet_Add(t *testing.T) {
	t.Parallel()

	s := set.NewSet()
	require.NotNil(t, s)
}

func TestSet_Size(t *testing.T) {
	t.Parallel()

	s := set.NewSet()
	require.Equal(t, s.Size(), 0)

	s.Add(1)
	s.Add(2)
	require.Equal(t, s.Size(), 2)

	s.Remove(1)
	require.Equal(t, s.Size(), 1)
}

func TestSet_Contains(t *testing.T) {
	t.Parallel()

	s := set.NewSet(1, "2")
	s.Add("3")

	require.True(t, s.Contains(1))
	require.True(t, s.Contains("2"))
	require.True(t, s.Contains("3"))
	require.False(t, s.Contains(4))
}

func TestSet_Remove(t *testing.T) {
	t.Parallel()

	s := set.NewSet(1, 2)
	s.Remove(1)

	require.False(t, s.Contains(1))
	require.True(t, s.Contains(2))
}

func TestSet_Clear(t *testing.T) {
	t.Parallel()

	s := set.NewSet(1, 2, 3, 4, 5)
	require.False(t, s.IsEmpty())
	s.Clear()
	require.True(t, s.IsEmpty())
}

func TestSet_Isempty(t *testing.T) {
	t.Parallel()

	s := set.NewSet()
	require.True(t, s.IsEmpty())

	s.Add("1")
	require.False(t, s.IsEmpty())
}

func TestSet_Compare(t *testing.T) {
	t.Parallel()

	s1 := set.NewSet(
		"2271",
		"2272",
		"2273",
		"2274",
		"2275",
		"2276",
		"2277",
		"2278",
		"2279",
		"22710",
	)

	s2 := set.NewSet(
		"2271",
		"22711",
		"2273",
		"2274",
		"2275",
		"2276",
		"2277",
		"2278",
		"2279",
		"22710",
	)

	s3 := set.NewSet(
		"2272",
		"2271",
		"2273",
		"2275",
		"2276",
		"2278",
		"2277",
		"2274",
		"2279",
		"22710",
	)

	require.False(t, s1.Compare(s2))
	require.True(t, s1.Compare(s3))
	require.False(t, s2.Compare(s3))
}

func TestSet_AddSet(t *testing.T) {
	t.Parallel()

	s1 := set.NewSet(1, 2)
	require.Equal(t, s1.Size(), 2)
	require.True(t, s1.Contains(1))
	require.False(t, s1.Contains("1"))

	s2 := set.NewSet("1", "2", "3")
	s1.AddSet(s2)
	require.Equal(t, s1.Size(), 5)
	require.True(t, s1.Contains("1"))
}

func TestSet_ConvertSlice(t *testing.T) {
	t.Parallel()

	s := set.NewSet()
	items := s.ConvertSlice()
	require.Len(t, items, 0)

	s = set.NewSet("aaa", "bbb", "ccc")
	items = s.ConvertSlice()
	require.Len(t, items, 3)

	for index := range items {
		require.True(t, s.Contains(items[index]))
	}
}
