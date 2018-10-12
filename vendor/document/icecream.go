package document

import "gopkg.in/mgo.v2/bson"

//Icecream properties
type Icecream struct {
	ID                    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name                  string        `bson:"name" json:"name"`
	ImageClosed           string        `bson:"image_closed" json:"image_closed"`
	ImageOpen             string        `bson:"image_open" json:"image_open"`
	Description           string        `bson:"description" json:"description"`
	Story                 string        `bson:"story" json:"story"`
	SourcingValues        []string      `bson:"sourcing_values" json:"sourcing_values"`
	Ingredients           []string      `bson:"ingredients" json:"ingredients"`
	AllergyInfo           string        `bson:"allergy_info" json:"allergy_info"`
	DietaryCertifications string        `bson:"dietary_certifications" json:"dietary_certifications"`
	ProductID             string        `bson:"productID" json:"productID"`
}
