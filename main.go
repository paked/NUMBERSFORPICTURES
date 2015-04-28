package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/paked/gerrycode/communicator"
	"github.com/paked/models"
	"gopkg.in/mgo.v2/bson"
)

// Image represents a struct that can be described with numbers
type Image struct {
	ID  bson.ObjectId `bson:"_id" json:"id"`
	URL string        `bson:"url" json:"url"`
}

func (i Image) BID() bson.ObjectId {
	return i.ID
}

func (i Image) C() string {
	return "images"
}

// Number represents a number that can describe an Image
type Number struct {
	ID     bson.ObjectId `bson:"_id" json:"id"`
	For    bson.ObjectId `bson:"for" json:"for"`
	Number float64       `bson:"number" json:"number"`
}

func (n Number) BID() bson.ObjectId {
	return n.ID
}

func (n Number) C() string {
	return "numbers"
}

func main() {
	models.Init("localhost", "imagenumberdescriptions")

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	r.HandleFunc("/", homeHandler).Methods("GET")

	api.HandleFunc("/image/new", createImageHandler).Methods("POST")
	api.HandleFunc("/image/{image_id}/number/new", addNumberHandler).Methods("POST")

	http.Handle("/", r)

	fmt.Println("Listening on port 8080")
	fmt.Println(http.ListenAndServe("localhost:8080", nil))
}

// Serves the home page, a grand spanking "hello!"
//   GET /
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

// Handler to create a new "image" object and store it in the database
//   POST /api/image/new?url=<URL>
func createImageHandler(w http.ResponseWriter, r *http.Request) {
	c := communicator.New(w)
	location := r.FormValue("url")

	u, err := url.Parse(location)
	if err != nil {
		c.Fail("Not a good URL!")
		return
	}

	location = u.String()

	var i Image
	if err := models.Restore(&i, bson.M{"url": location}); err == nil {
		c.Fail("An image with that URL already exists")
		return
	}

	i = Image{
		ID:  bson.NewObjectId(),
		URL: location,
	}

	if err := models.Persist(i); err != nil {
		c.Error("An error with the database has occured")
		return
	}

	c.OKWithData("Here is your image", i)
}

// A handler to add a number to describe an image
//   POST /image/{image_id}/number/new?number=<A NUMBER>
func addNumberHandler(w http.ResponseWriter, r *http.Request) {
	c := communicator.New(w)
	imageId := mux.Vars(r)["image_id"]
	numberString := r.FormValue("number")

	number, err := strconv.ParseFloat(numberString, 64)
	if err != nil {
		c.Fail("That is not a valid float!")
		return
	}

	if !bson.IsObjectIdHex(imageId) {
		c.Fail("Not valid image id!")
		return
	}

	n := Number{
		ID:     bson.NewObjectId(),
		For:    bson.ObjectIdHex(imageId),
		Number: number,
	}

	if err := models.Persist(n); err != nil {
		c.Error("Something bad happened in the database...")
		return
	}

	c.OKWithData("Here is your number", n)
}
