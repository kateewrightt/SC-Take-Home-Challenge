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

// This function checks various scenarios for fetching folders based on different org IDs.
func Test_GetAllFolders(t *testing.T) {

	// This test ensures that the function correctly handles a standard, expected input
	// A valid organisation ID is used to verify that the function retrieves the correct subset of folders
	t.Run("returns all folders with valid orgID", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.FromStringOrNil(folders.DefaultOrgID), // Using the default org ID.
		}

		// Set up the dependencies with the mock data fetcher.
		deps := folders.FetchFolderDependencies{
			DataFetcher: MockFetcher{}, // Injecting the mock data fetcher.
		}
		foldersResp, err := folders.GetAllFolders(req, deps)

		// Verify that each returned folder has the correct OrgID
		for _, folder := range foldersResp.Folders {
			assert.Equal(t, req.OrgID, folder.OrgId)
		}
		assert.Nil(t, err)
	})

	// This test checks the function's behaviour when given an empty organisation ID
	// An empty organisation ID should logically return no folders, testing the function's input validation
	t.Run("returns no folders with empty orgID", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.Nil, // Using an empty org ID
		}

		// Set up the dependencies with the mock data fetcher
		deps := folders.FetchFolderDependencies{
			DataFetcher: MockFetcher{}, // Injecting the mock data fetcher
		}
		foldersResp, err := folders.GetAllFolders(req, deps)

		// We expect to get no folders as the organisation ID is empty
		// This ensures that the function can handle cases where the input is not valid or missing
		assert.Len(t, foldersResp.Folders, 0)
		assert.Nil(t, err)
	})

	// This test ensures the function can handle invalid input correctly
	// An invalid organisation ID should return no folders, testing the robustness of the function against invalid data
	t.Run("returns no folders with invalid orgID", func(t *testing.T) {
		// Generate a random UUID to use as an invalid organisation ID.
		invalidUUID, _ := uuid.NewV4()
		req := &folders.FetchFolderRequest{
			OrgID: invalidUUID, // Using an invalid org ID
		}

		deps := folders.FetchFolderDependencies{
			DataFetcher: MockFetcher{}, // Injecting the mock data fetcher
		}

		// Call the GetAllFolders function with the request and dependencies
		foldersResp, err := folders.GetAllFolders(req, deps)

		// We expect to get no folders since the organsation ID does not match any in the data set
		// This verifies that the function can correctly identify and handle invalid organisation IDs
		assert.Len(t, foldersResp.Folders, 0)
		assert.Nil(t, err)
	})

	// Test to ensure only non-deleted folders are returned
	t.Run("returns only non-deleted folders", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.FromStringOrNil(folders.DefaultOrgID), // Using the default org ID
		}

		// Set up the dependencies with the mock data fetcher
		deps := folders.FetchFolderDependencies{
			DataFetcher: MockFetcher{}, // Injecting the mock data fetcher
		}
		foldersResp, err := folders.GetAllFolders(req, deps)

		// Check that no folder in the response is marked as deleted
		for _, folder := range foldersResp.Folders {
			assert.False(t, folder.Deleted)
		}
		assert.Nil(t, err)
	})
}

// This function checks various scenarios to ensure that pagination works correctly
func Test_Pagination(t *testing.T) {

	// This test ensures that the function can retrieve the first page of results correctly
	t.Run("returns first page with valid orgID", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.FromStringOrNil(folders.DefaultOrgID), // Using the default org ID
		}

		deps := folders.FetchFolderDependencies{
			DataFetcher: MockFetcher{}, // Injecting the mock data fetcher
		}
		pageSize := 5 // Setting the page size to 5
		token := ""   // No token for the first page

		// Call the FetchFoldersWithPagination function with the request, dependencies, page size, and token
		foldersResp, err := folders.FetchFoldersWithPagination(req, deps, pageSize, token)

		// We expect to get the first 5 folders
		// This verifies that the function correctly implements pagination starting from the first item
		assert.Nil(t, err)
		assert.Len(t, foldersResp.Folders, pageSize)
		assert.NotEmpty(t, foldersResp.Token) // Ensure the token is not empty for the next page
	})

	// This test ensures that the function can retrieve the second page of results correctly using a token
	t.Run("returns second page with valid token", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.FromStringOrNil(folders.DefaultOrgID), // Using the default org ID
		}

		deps := folders.FetchFolderDependencies{
			DataFetcher: MockFetcher{}, // Injecting the mock data fetcher
		}
		pageSize := 5 // Setting the page size to 5

		// Fetch the first page to get the token for the second page
		firstPageResp, _ := folders.FetchFoldersWithPagination(req, deps, pageSize, "")

		// Fetch the second page using the token from the first page response
		secondPageResp, err := folders.FetchFoldersWithPagination(req, deps, pageSize, firstPageResp.Token)

		// We expect to get the next 5 folders
		// This verifies that the function correctly continues pagination from the given token
		assert.Nil(t, err)
		assert.Len(t, secondPageResp.Folders, pageSize)
		assert.NotEqual(t, firstPageResp.Folders, secondPageResp.Folders) // Ensure the folders are different from the first page
	})

	// This test ensures that the function handles an invalid token correctly by returning an error
	t.Run("returns no folders with invalid token", func(t *testing.T) {
		req := &folders.FetchFolderRequest{
			OrgID: uuid.FromStringOrNil(folders.DefaultOrgID), // Using the default org ID
		}

		deps := folders.FetchFolderDependencies{
			DataFetcher: MockFetcher{}, // Injecting the mock data fetcher
		}
		pageSize := 5
		invalidToken := "invalidToken" // Setting an invalid token

		// Call the FetchFoldersWithPagination function with the request, dependencies, page size, and invalid token
		foldersResp, err := folders.FetchFoldersWithPagination(req, deps, pageSize, invalidToken)

		// We expect an error to be returned and no folders to be fetched
		// This verifies that the function correctly handles invalid tokens
		assert.NotNil(t, err)
		assert.Nil(t, foldersResp)
	})
}
