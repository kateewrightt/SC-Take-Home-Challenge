package folders_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

type MockFetcher struct{}

func (f MockFetcher) GetFolders() []*folders.Folder {
	testData := folders.GenerateData()
	return testData
}
func Test_GetAllFolders(t *testing.T) {
	t.Run("returns all folders with valid orgID", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.FromStringOrNil(folders.DefaultOrgID),
		}

		deps := folders.FetchFolderDependencies{
			DataFetcher: MockFetcher{},
		}
		foldersResp, err := folders.GetAllFolders(req, deps)

		// Reading the generate data code, we set 1/3 folders to be a random org id
		// The rest are default org id
		// In a real test setup, these should be arguments passed into the testing mock
		assert.Len(t, foldersResp.Folders, (folders.DataSetSize / 3 * 2))
		assert.Nil(t, err)
	})

	t.Run("returns no folders with empty orgID", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.Nil,
		}

		deps := folders.FetchFolderDependencies{
			DataFetcher: MockFetcher{},
		}
		foldersResp, err := folders.GetAllFolders(req, deps)

		assert.Len(t, foldersResp.Folders, 0)
		assert.Nil(t, err)
	})

	t.Run("returns no folders with invalid orgID", func(t *testing.T) {
		invalidUUID, _ := uuid.NewV4()
		req := &folders.FetchFolderRequest{
			OrgID: invalidUUID,
		}

		deps := folders.FetchFolderDependencies{
			DataFetcher: MockFetcher{},
		}
		foldersResp, err := folders.GetAllFolders(req, deps)

		assert.Len(t, foldersResp.Folders, 0)
		assert.Nil(t, err)
	})

	t.Run("error?", func(t *testing.T) {

	})
}
