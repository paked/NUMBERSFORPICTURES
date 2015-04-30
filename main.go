package main

import (
	"fmt"
	"math/rand"
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
	Number int           `bson:"number" json:"number"`
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

	api.HandleFunc("/images/new", createImageHandler).Methods("POST")
	api.HandleFunc("/images/{image_id}/numbers/new", addNumberHandler).Methods("POST")
	api.HandleFunc("/images/{image_id}/numbers", getNumbersHandler).Methods("GET")
	api.HandleFunc("/images/random", getRandomImageHandler).Methods("GET")
	api.HandleFunc("/images/{image_id}", getImageHandler).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))

	http.Handle("/", r)

	fmt.Println("Listening on port 8080")
	fmt.Println(http.ListenAndServe("localhost:8080", nil))
}

// createImageHandler is a handler to create a new "image" object and store it in the database
//   POST /api/images/new?url=<URL>
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

// addNumberHandler is a handler to add a number to describe an image
//   POST /api/images/{image_id}/numbers/new?number=<A NUMBER>
func addNumberHandler(w http.ResponseWriter, r *http.Request) {
	c := communicator.New(w)
	imageId := mux.Vars(r)["image_id"]
	numberString := r.FormValue("number")

	number, err := strconv.Atoi(numberString)
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

// getNumbersHandler is a handler to get all the numbers which describe an image
//   GET /api/images/<IMAGE_ID>/numbers
func getNumbersHandler(w http.ResponseWriter, r *http.Request) {
	c := communicator.New(w)
	imageID := mux.Vars(r)["image_id"]

	if !bson.IsObjectIdHex(imageID) {
		c.Fail("The image id you provided is not valid!")
		return
	}

	var numbers []Number
	number := Number{}

	iter, err := models.Fetch(number.C(), bson.M{"for": bson.ObjectIdHex(imageID)}, "_id")
	if err != nil {
		c.Error("Something went wrong during the fetching!")
		return
	}

	for iter.Next(&number) {
		numbers = append(numbers, number)
		fmt.Println(number)
	}

	c.OKWithData("Here are your numbers", numbers)
}

// getImageHandler is a handler that gives an image
//   GET /api/images/<IMAGE_ID>
func getImageHandler(w http.ResponseWriter, r *http.Request) {
	c := communicator.New(w)
	imageID := mux.Vars(r)["image_id"]

	if !bson.IsObjectIdHex(imageID) {
		c.Fail("The image id you provided was not valid.")
		return
	}

	var i Image
	if err := models.RestoreByID(&i, bson.ObjectIdHex(imageID)); err != nil {
		c.Fail("An image with that ID does not exist")
		return
	}

	c.OKWithData("Here is the image", i)
}

// getRandomImageHandler is a handler that gives a random image
// GET /api/image/random
func getRandomImageHandler(w http.ResponseWriter, r *http.Request) {
	c := communicator.New(w)

	var images []Image
	image := Image{}

	iter, err := models.Fetch(image.C(), bson.M{}, "_id")
	if err != nil {
		c.Error("Something went wrong during the fetching!")
		return
	}

	for iter.Next(&image) {
		images = append(images, image)
	}

	index := rand.Intn(len(images))

	c.OKWithData("Here is your random image", images[index])
}
