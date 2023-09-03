package inner_storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInnerStorage(t *testing.T) {
	is := NewInnerStorage()

	data := "string"

	id, resData := is.Create(data)
	require.Equal(t, data, resData)

	resData, err := is.Get(id)
	require.NoError(t, err)
	require.Equal(t, data, resData)

	updData := "new data"

	resData, err = is.Update(id, updData)
	require.NoError(t, err)
	require.Equal(t, updData, resData)

	resData, err = is.Get(id)
	require.NoError(t, err)
	require.Equal(t, updData, resData)

	resData, err = is.Delete(id)
	require.NoError(t, err)
	require.Equal(t, updData, resData)

	resData, err = is.Get(id)
	require.ErrorIs(t, err, ErrNotFound)
}
