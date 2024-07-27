package folders

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

// My approach focused on dividing a large dataset into smaller and more manageable chunks
// where each is retrievable by a token that represents the starting index of the next subset.
// The main goal here is to facilitate efficient data retrieval and avoid overwhelming the client
// with too much data at once. The function uses a token system to keep track of the current
// position within the dataset. The token is a base64 encoded string representing the starting
// index for the next page. This allows the client to request subsequent pages by providing the
// token, ensuring continuity in data retrieval.

// Initially the token is empty, indicating the start of the dataset. For each subsequent request
// the token is decoded to determine the starting index, and a new token is generated for the next page,
// encoding the end index.

// The function includes error handling to manage invalid tokens. If a token cannot be decoded or
// converted to an integer, an error is returned which ensures the function does not process invalid input.
// The pagination logic calculates the end index based on the starting index and page size. If the end
// index exceeds the datasets length then it is adjusted to avoid out of bounds errors and the subset of
// folders corresponding to the current page is sliced from the main dataset and returned along with
// the next token if more data remains.

// Weaknesses of my approach:
//  - The reliance on simple base64 encoding of the starting index does not validate the tokens integrity
//	   or prevent tampering which could lead to incorrect data retrieval if a client modifies the token.
//  - The function fetches all folders for the given organisation ID upfront before slicing the array for
//     pagination which could lead to high memory usage for large datasets.
//  - The page size is fixed which may not be flexible for different use cases or client preferences.

func FetchFoldersWithPagination(req *FetchFolderRequest, deps FetchFolderDependencies, pageSize int, token string) (*FetchFolderResponse, error) {
	folders, err := FetchAllFoldersByOrgID(req.OrgID, deps.DataFetcher)
	if err != nil {
		return nil, err
	}

	// Initalise starting index for pagination
	startIndex := 0
	if token != "" {
		// Decode token to get the starting index for the current page
		// The token is a base64 encoded string representing the index position in the folder lis
		decodedToken, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return nil, fmt.Errorf("invalid token: %v", err)
		}
		// Convert decoded token to an integer to determine the starting index
		startIndex, err = strconv.Atoi(string(decodedToken))
		if err != nil {
			return nil, fmt.Errorf("invalid token: %v", err)
		}
	}

	// Calculate the ending index for the current page of results.
	endIndex := startIndex + pageSize

	// Adjust the ending index if it exceeds the total number of folders.
	if endIndex > len(folders) {
		endIndex = len(folders)
	}

	// Slice the folders array to get the current page's items.
	paginatedFolders := folders[startIndex:endIndex]

	// Generate the next token if there are more folders to fetch
	var nextToken string
	if endIndex < len(folders) {
		// Encode the ending index as the next token
		// The token will be used in the next request to fetch the subsequent page of folders
		nextToken = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(endIndex)))
	}

	// Create and return the paginated response
	response := &FetchFolderResponse{
		Folders: paginatedFolders,
		Token:   nextToken,
	}

	return response, nil
}

// Function to fetch and print folders using pagination.
func FetchAndPrintFoldersWithPagination(req *FetchFolderRequest, deps FetchFolderDependencies, pageSize int) {
	token := ""
	for {
		// Fetch folders with the given page size and token
		response, err := FetchFoldersWithPagination(req, deps, pageSize, token)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		PrettyPrint(response)

		// Break the loop if there are no more pages
		if response.Token == "" {
			break
		}

		// Update the token for the next page fetch
		token = response.Token
	}
}
