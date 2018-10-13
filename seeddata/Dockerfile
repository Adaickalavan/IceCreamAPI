#Start from a Debian image with the latest version of Mongo
FROM mongo

#Author
LABEL Author Adaickalavan Meiyappan

#Copy the local package files (from development computer) to the container's workspace (docker image)
COPY icecream.json /icecream.json

#Import seed data from json file into Mongo database
CMD mongoimport --host mongo --db product --collection icecream --type json --file /icecream.json --jsonArray
# --host : specifies a resolvable hostname for mongo database to which to connect to
# --db : specifies the name of database on which to run the mongoimport
# --collection : specifies the collection to import to
# --type : specifies the file type to import
# --file : specifies the location and name of a file containing the data to import 
# --jsonArray : accepts the import of data expressed with multiple MongoDB documents within a single JSON array
