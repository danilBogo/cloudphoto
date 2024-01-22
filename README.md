Requires an installed go

The order of commands:
1) go build
2) go install

Then you can use the commands:

1) cloudphoto init (program initialization)

2) cloudphoto upload --album ALBUM [--path PHOTOS_DIR] (sending photos to the cloud storage)

3) cloudphoto download --album ALBUM [--path PHOTOS_DIR] (uploading photos from cloud storage)

4) cloudphoto list [--album ALBUM] (view the list of albums and photos in the album)

5) cloudphoto delete --album ALBUM [--photo PHOTO] (deleting albums and photos)

6) cloudphoto mksite (generation and publication of photo archive web pages)
