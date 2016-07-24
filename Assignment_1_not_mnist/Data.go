package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	imageSize  = 28    // Pixel width and height.
	pixelDepth = 255.0 // Number of levels per pixel.
)

func loadPNG(imageFileName string) {

	imageFile, err := os.Open(imageFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer imageFile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(imageFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new grayscale image
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	data := make([][]float32, w)

	for x := 0; x < w; x++ {
		data[x] = make([]float32, h)
		for y := 0; y < h; y++ {
			color := src.At(x, y)
			r, g, b, a := color.RGBA()
			fmt.Printf("(r,g,b,a)=(%d, %d, %d, %d)\n", r, g, b, a)
			data[x][y] = (float32(r) - pixelDepth/2) / pixelDepth
		}
	}
}

func pickle(dataFolders []string, minNumberImages int, force bool) []string {

	var datasetNames []string

	for _, folder := range dataFolders {
		pickleFilename := folder + ".pickle"
		fmt.Println("pickleFilename:", pickleFilename)
		append(dataset_names, pickleFilename)

		// Test if pickleFilename exists
		if _, err := os.Stat(pickleFilename); os.IsNotExist(err) && force == false {
			// pickleFilename does not exist
			fmt.Println(pickleFilename, "already present - Skipping pickling.")
		} else {
			fmt.Println("Pickling", pickleFilename)

			dataset := loadLetter(folder, minNumberImages)
			f, err := os.Create(pickleFilename)
			if err != nil {
				log.Fatal(err)
			}

			// dump serialize(dataset)
		}
	}

	return datasetNames
}

// loadLetter loads the data for a single letter label.
func loadLetter(folder String, minNumberImages int) {

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}
	for _, imageFile := range files {
	}

	/*
	   dataset = np.ndarray(shape=(len(image_files), image_size, image_size),
	                          dtype=np.float32)
	   image_index = 0
	   print(folder)
	   for image in os.listdir(folder):
	     image_file = os.path.join(folder, image)
	     try:
	       image_data = (ndimage.imread(image_file).astype(float) -
	                     pixel_depth / 2) / pixel_depth
	       if image_data.shape != (image_size, image_size):
	         raise Exception('Unexpected image shape: %s' % str(image_data.shape))
	       dataset[image_index, :, :] = image_data
	       image_index += 1
	     except IOError as e:
	       print('Could not read:', image_file, ':', e, '- it\'s ok, skipping.')

	   num_images = image_index
	   dataset = dataset[0:num_images, :, :]
	   if num_images < min_num_images:
	     raise Exception('Many fewer images than expected: %d < %d' %
	                     (num_images, min_num_images))

	   print('Full dataset tensor:', dataset.shape)
	   print('Mean:', np.mean(dataset))
	   print('Standard deviation:', np.std(dataset))
	   return dataset
	*/
}

func load(dataFolders []string, minNumberImages int, maxNumberImages int) ([][imageSize][imageSize]float32, []int32) {

	dataset := make([][imageSize][imageSize]float32, maxNumberImages)

	labels := make([]int32, maxNumberImages)

	labelIndex := 0

	imageIndex := 0
	for _, folder := range dataFolders {
		fmt.Println(folder)

		files, err := ioutil.ReadDir(folder)
		if err != nil {
			log.Fatal(err)
		}
		for _, imageFile := range files {

			if imageIndex >= maxNumberImages {
				log.Fatal("More images than expected: ", imageIndex, " >= ", maxNumberImages)
			}
			imagePath := filepath.Join(folder, imageFile.Name())
			fmt.Println(imagePath)

			var imageData [28][28]float32

			if len(imageData) != imageSize && len(imageData[0]) != imageSize {
				log.Fatal("Unexpected image shape: ", len(imageData), "x", len(imageData[0]))
			}

			dataset[imageIndex] = imageData
			labels[imageIndex] = int32(labelIndex)
			imageIndex++

			/*
			   try:
			     image_data = (ndimage.imread(image_file).astype(float) -
			                   pixelDepth / 2) / pixelDepth
			     dataset[image_index, :, :] = image_data
			     labels[image_index] = label_index
			     image_index += 1
			   except IOError as e:
			     print 'Could not read:', image_file, ':', e, '- it\'s ok, skipping.'
			*/
		}

		labelIndex++
	}

	numberImages := imageIndex
	dataset = dataset[0:numberImages][:][:]
	labels = labels[0:numberImages]
	if numberImages < minNumberImages {
		log.Fatal("Many fewer images than expected: ", numberImages, " < ", minNumberImages)
	}

	fmt.Println("Full dataset tensor: dataset[", len(dataset), "][", len(dataset[0]), "][", len(dataset[0][0]), "]")
	//	fmt.Println("Mean:", np.mean(dataset))
	//	fmt.Println("Standard deviation:", np.std(dataset))
	fmt.Println("Labels[", len(labels), "]")

	return dataset, labels
}

/*
def load(data_folders, min_num_images, max_num_images):
  dataset = np.ndarray(
    shape=(max_num_images, imageSize, imageSize), dtype=np.float32)
  labels = np.ndarray(shape=(max_num_images), dtype=np.int32)
  label_index = 0
  image_index = 0
  for folder in data_folders:
    print folder
    for image in os.listdir(folder):
      if image_index >= max_num_images:
        raise Exception('More images than expected: %d >= %d' % (
          num_images, max_num_images))
      image_file = os.path.join(folder, image)
      try:
        image_data = (ndimage.imread(image_file).astype(float) -
                      pixelDepth / 2) / pixelDepth
        if image_data.shape != (imageSize, imageSize):
          raise Exception('Unexpected image shape: %s' % str(image_data.shape))
        dataset[image_index, :, :] = image_data
        labels[image_index] = label_index
        image_index += 1
      except IOError as e:
        print 'Could not read:', image_file, ':', e, '- it\'s ok, skipping.'
    label_index += 1


  num_images = image_index
  dataset = dataset[0:num_images, :, :]
  labels = labels[0:num_images]
  if num_images < min_num_images:
    raise Exception('Many fewer images than expected: %d < %d' % (
        num_images, min_num_images))
  print 'Full dataset tensor:', dataset.shape
  print 'Mean:', np.mean(dataset)
  print 'Standard deviation:', np.std(dataset)
  print 'Labels:', labels.shape
  return dataset, labels

*/

func main() {
	/*
		const url = "http://commondatastorage.googleapis.com/books1000/"
		DownloadArchive(url, "notMNIST_small.tar.gz", 8458043, false)
		err := ExtractArchive("notMNIST_small.tar.gz", false)
		if err != nil {
			log.Fatal(err)
		}
	*/
	trainFolders := []string{"notMNIST_small/A", "notMNIST_small/B"}

	load(trainFolders, 450000, 550000)
}
