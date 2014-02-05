package imageUploader

import (
	"fmt"
	"time"

	"net/http"

	"appengine"
	"appengine/datastore"
)

type UploadFile struct {
	UploadDate time.Time
}

func init() {
	http.HandleFunc("/", rootPage)
	http.HandleFunc("/uploadTo", uploadToPage)
}

func uploadFileKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "UploadFile", "default_uploadfile", 0, nil)
}

func rootPage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	fmt.Fprint(w, uploadPage)

	q := datastore.NewQuery("uploads").Ancestor(uploadFileKey(c)).Order("-UploadDate").Limit(10)
	uploadfiles := make([]UploadFile, 0, 10)
	if _, err := q.GetAll(c, &uploadfiles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, val := range uploadfiles {
		fmt.Fprintf(w, "%v\n", val)
	}

}

const uploadPage = `
<!DOCTYPE html>
<html>
<head>
	<title>Upload</title>
</head>
<body>
    <form action="/uploadTo" method="post">
      <div><input type="submit" value="Upload"></div>
    </form>
</body>
</html>
`

func upload(w http.ResponseWriter, r *http.Request) {
}

func uploadToPage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	fmt.Fprintf(w, "Inserted file at time: (%v) ", time.Now())
	f := UploadFile{
		UploadDate: time.Now(),
	}
	key := datastore.NewIncompleteKey(c, "uploads", uploadFileKey(c))
	_, err := datastore.Put(c, key, &f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
