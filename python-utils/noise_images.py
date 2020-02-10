import random, os, sys
import cv2
import numpy as np

# Erodes the image
def erode(img, erosion_size):
    kernel = np.ones((erosion_size,erosion_size), np.uint8) 
    erosion_dst = cv2.erode(img, kernel, iterations=1)
    return erosion_dst


# Adds salt and pepper noise to the image
def add_noise(image,prob):
    output = np.zeros(image.shape,np.uint8)
    thres = 1 - prob 
    for i in range(image.shape[0]):
        for j in range(image.shape[1]):
            rdn = random.random()
            if rdn < prob:
                output[i][j] = 0
            elif rdn > thres:
                output[i][j] = 255
            else:
                output[i][j] = image[i][j]
    return output
    


if __name__ == "__main__":
    
    dir = sys.argv[1]
    outdir = sys.argv[2]
    erosion_size = int(sys.argv[3])
    prob = float(sys.argv[4])
    os.mkdir(outdir)
    for filename in os.listdir(dir):
        print(filename)
        img = cv2.imread(os.path.join(dir, filename))
        
        processed_img = erode(img, erosion_size)
        processed_img = add_noise(img, prob)
        cv2.imwrite(os.path.join(outdir, filename), processed_img)

    
    
    
    
    
    