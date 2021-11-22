# go-test

REQUIREMENTS:
* auth token matches --> OK
* POST data to the /upload --> OK
* should write the received file data to a temporary file -->  OK
* content type of the uploaded file is an image --> OK
* Images larger than 8 megabytes should also be rejected --> OK
* If the submission is bad, please return a 403 HTTP error code --> OK
* Write the image metadata (content type, size, etc) to a database of your choice, including all relevant HTTP information. --> OK

* nosurf : library for auth token generation & validation 

* Uploaded images are saved in the temp-images folder
* Image HTTP information (filename, filesize, content-type) are saved in the database (tbl_image)
