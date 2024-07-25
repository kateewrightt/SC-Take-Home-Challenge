package folders

import (
	"github.com/gofrs/uuid"
)

// GetAllFolders fetches all folders for an org
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {

	// Unused variables might bring back when implementing pagination
	// var (
	//	err error
	//	f1  Folder
	//	fs  []*Folder
	// )
	// f := []Folder{}

	// Fetch all folders by org ID
	// Add error handling
	r, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}

	// Converts map vaules to a slice of folders
	// This loop iterates over the map returned by FetchAllFoldersByOrgID and converts
	// each value to a folder struct and then appends each folder struct to the slice f
	// for k, v := range r {
	//	f = append(f, *v)
	// }

	// Converts the slice of folders to a slice of pointers to folder
	// This second loop then takes the slice of folder structs f and
	// converts it back to a slice of pointers to folder struct fp
	//var fp []*Folder
	//for k1, v1 := range f {
	//	fp = append(fp, &v1)
	//}

	//var ffr *FetchFolderResponse
	//ffr = &FetchFolderResponse{Folders: fp}
	//return ffr, nil

	// Simplified creation of FetchFolderResponse
	ffr := &FetchFolderResponse{Folders: r}
	return ffr, nil
}

func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	// We want the retrieval of data here to be a modular component in our app
	// Here the default method is getting data through local json
	// We want to be able to generate this data in tests too
	// To faciliate this we will change the hardcoded function call into a interface
	// And inject the data retrieval as an application dependency
	folders := GetSampleData()

	resFolder := []*Folder{}
	for _, folder := range folders {
		if folder.OrgId == orgID {
			resFolder = append(resFolder, folder)
		}
	}
	return resFolder, nil
}
