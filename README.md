# similarimagefinder

* Finds similar JPEG images (using phash algorithm) in an specified folder (includes subfolders).
* Generates an html file displaying similar images in a row, each image has a "delete" link.
* User opens html file and decides which images to delete clicking the "delete" link.

## Steps on how to run it

1. Execute the command
  * ./similarimagefinder -folder=\<specify_root_folder_containing_images\> -outFile=\<html_result_file\>
2. Wait for message "Starting service on port", if no duplicate image has been found then application will exit.
3. Open html file (the one defined using the parameter -outFile) in a web browser (like Chrome)
4. Review and delete the images you do not want to keep clicking the "Delete" link. (Each row of images are a set of similar images)
5. Once you are done, close the html page and stop the application and delete the generated html result file.
